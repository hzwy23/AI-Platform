package protocol

import (
	"ai-platform/panda/logger"
	"container/list"
	"errors"
	"net"
	"sync"
	"time"
)

type udpCmd struct {
	addr    *net.UDPAddr
	message []byte
}

// JTTUDPProtocol UDP JTT格式数据协议
type JTTUDPProtocol struct {
	buffer    map[string][]byte
	message   *list.List
	conn      net.Conn
	lock      *sync.RWMutex
	msgLock   *sync.RWMutex
	isClosed  bool
	closeLock *sync.RWMutex
}

// NewUDPJTTProtocol 创建UDP协议
func NewUDPJTTProtocol(conn *net.UDPConn) *JTTUDPProtocol {
	r := &JTTUDPProtocol{
		buffer:    make(map[string][]byte),
		message:   list.New(),
		conn:      conn,
		lock:      new(sync.RWMutex),
		msgLock:   new(sync.RWMutex),
		isClosed:  false,
		closeLock: new(sync.RWMutex),
	}
	go r.readFromUDP(conn)
	return r
}

// Send 发送数据
func (r *JTTUDPProtocol) Send(msgID uint16, msgData []byte) (int, error) {
	data, err := Pack(msgID, msgData)
	if err != nil {
		logger.Warn(err)
		return 0, err
	}
	return r.conn.Write(data)
}

// Parse 读取并解析数据
func (r *JTTUDPProtocol) Parse() ([]byte, error) {
	// 将buffer内容读取出去，并清空buffer
	r.lock.Lock()
	element := r.message.Front()
	if element != nil {
		r.message.Remove(element)
	}
	r.lock.Unlock()
	if element == nil {
		return nil, nil
	}
	msg, ok := element.Value.(udpCmd)
	if ok {
		return msg.message, nil
	}
	r.closeLock.RLock()
	defer r.closeLock.RUnlock()
	if r.isClosed {
		return nil, errors.New("连接已经断开")
	}
	return nil, nil
}

func (r *JTTUDPProtocol) readFromUDP(conn *net.UDPConn) {
	for {
		tmp := make([]byte, 2048)
		size, raddr, err := conn.ReadFromUDP(tmp)
		if err != nil {
			logger.Error("读取socket内容失败，失败原因是：", err)
			r.closeLock.Lock()
			r.isClosed = true
			r.closeLock.Unlock()
			break
		}
		if size == 0 {
			time.Sleep(100 * 1e6)
			continue
		}
		// 如果解析到message，则触发相应的处理逻辑
		logger.Debug("receive byte is: ", tmp[:size])
		if val, ok := r.buffer[raddr.String()]; ok {
			ret := append(val, tmp[:size]...)
			if msg, ok, idx := r.parse(ret); ok {
				r.lock.Lock()
				r.message.PushBack(udpCmd{
					addr:    raddr,
					message: msg,
				})
				if idx == len(ret) {
					delete(r.buffer, raddr.String())
				} else {
					r.buffer[raddr.String()] = ret[idx+1:]
				}
				r.lock.Unlock()
			} else {
				r.buffer[raddr.String()] = ret
			}
		} else {
			val = make([]byte, 0)
			ret := append(val, tmp[:size]...)
			if msg, ok, idx := r.parse(ret); ok {
				r.lock.Lock()
				r.message.PushBack(udpCmd{
					addr:    raddr,
					message: msg,
				})
				if idx == len(ret) {
					delete(r.buffer, raddr.String())
				} else {
					r.buffer[raddr.String()] = ret[idx+1:]
				}
				r.lock.Unlock()
			} else {
				r.buffer[raddr.String()] = ret
			}
		}
	}
}

// 解析是否获取一个完成的报文
func (r *JTTUDPProtocol) parse(tmp []byte) ([]byte, bool, int) {
	endFlag := false
	startFlag := false
	cnt := 0
	for idx, item := range tmp {
		if item == HEADER_FLAG {
			startFlag = true
		} else if startFlag && item == FOOTER_FLAG {
			endFlag = true
			cnt = idx
			break
		} else if item == FOOTER_FLAG {
			// 删除前边的数据
			cnt = idx
		}
	}
	// 如果获取到结束符，则处理消息
	if startFlag && endFlag {
		r.msgLock.RLock()
		endMsg := tmp[:cnt+1]
		r.msgLock.RUnlock()
		retMsg, yes := transfer(endMsg)
		return retMsg, yes, cnt
	}

	return nil, false, 0
}
