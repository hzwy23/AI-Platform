package protocol

import (
	"ai-platform/panda/logger"
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

// JTTProtocol TCP JTT格式数据协议
type JTTProtocol struct {
	buffer    []byte
	message   []byte
	conn      net.Conn
	lock      *sync.RWMutex
	msgLock   *sync.RWMutex
	isClosed  bool
	closeLock *sync.RWMutex
}

// 数据解析
// 1. 转义处理
// 2. 长度校验
// 3. 解密body部分
// 4. CRC校验body部分

// NewJTTProtocol 封装数据
// 1. 生成数据body
// 2. CRC校验Body部分，生成CRC校验码，并设置到报文相关字段CrcCode
// 3. 加密body部分
// 4. 统计报文长度，并设置Header中报文长度相关字段MsgLength
// 4. 转义处理
// 5. 发送数据
func NewJTTProtocol(conn net.Conn) *JTTProtocol {
	r := &JTTProtocol{
		buffer:    make([]byte, 0),
		message:   make([]byte, 0),
		conn:      conn,
		lock:      new(sync.RWMutex),
		msgLock:   new(sync.RWMutex),
		isClosed:  false,
		closeLock: new(sync.RWMutex),
	}
	go r.read()
	return r
}

// Send 发送数据
func (r *JTTProtocol) Send(msgID uint16, msgData []byte) (int, error) {
	data, err := Pack(msgID, msgData)
	if err != nil {
		logger.Warn(err)
		return 0, err
	}
	return r.conn.Write(data)
}

// Parse 读取数据
func (r *JTTProtocol) Parse() ([]byte, error) {

	// 将buffer内容读取出去，并清空buffer
	r.lock.Lock()
	tmp := r.buffer
	r.buffer = make([]byte, 0)
	r.lock.Unlock()

	// 解析读取的buffer数据，并从中获取有效的报文信息
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

func (r *JTTProtocol) read() {
	for {
		tmp := make([]byte, 256)
		size, err := r.conn.Read(tmp)
		if err == io.EOF {
			logger.Info("连接已断开：", err)
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
func (r *JTTProtocol) parse(tmp []byte) ([]byte, bool) {
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
		endMsg := r.message
		r.msgLock.RUnlock()
		return transfer(endMsg)
	}
	return nil, false
}
