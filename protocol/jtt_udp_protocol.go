package protocol

import (
	"ai-platform/panda/logger"
	"net"
	"sync"
	"time"
	"errors"
)

// JTTUDPProtocol UDP JTT格式数据协议
type JTTUDPProtocol struct {
	buffer    []byte
	message   []byte
	conn      net.Conn
	lock      *sync.RWMutex
	msgLock   *sync.RWMutex
	isClosed  bool
	closeLock *sync.RWMutex
}


// NewUDPJTTProtocol 创建UDP协议
func NewUDPJTTProtocol(conn *net.UDPConn) *JTTUDPProtocol {
	r := &JTTUDPProtocol{
		buffer:    make([]byte, 0),
		message:   make([]byte, 0),
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
	tmp := r.buffer
	r.buffer = make([]byte, 0)
	r.lock.Unlock()

	if msg, ok := r.parse(tmp); ok {
		logger.Debug("receive message is:", msg)
		return msg, nil
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
		size,_, err := conn.ReadFromUDP(tmp)
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
		r.lock.Lock()
		logger.Debug("receive byte is: ", tmp[:size])
		r.buffer = append(r.buffer, tmp[:size]...)
		r.lock.Unlock()
	}
}

// 解析是否获取一个完成的报文
func (r *JTTUDPProtocol) parse(tmp []byte) ([]byte, bool) {
	endFlag := false

	for idx, item := range tmp {
		if item == HEADER_FLAG {
			r.msgLock.Lock()
			r.message = make([]byte, 0)
			r.message = append(r.message, item)
			r.msgLock.Unlock()
		} else {
			r.msgLock.Lock()
			r.message = append(r.message, item)
			r.msgLock.Unlock()
		}
		if item == FOOTER_FLAG {
			endFlag = true
			if idx+1 < len(tmp) {
				r.lock.Lock()
				r.buffer = append(tmp[idx+1:], r.buffer...)
				r.lock.Unlock()
			}
			break
		}
	}

	// 如果获取到结束符，则处理消息
	if endFlag {
		r.msgLock.RLock()
		endMsg:=r.message
		r.msgLock.RUnlock()
		return transfer(endMsg)
	}
	return nil, false
}
