package protocol

import (
	"ai-platform/panda/logger"
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"time"
)

func NewJTTProtocol(conn net.Conn) *JTTProtocol {
	r := &JTTProtocol{
		buffer:    make([]byte, 0),
		message:   make([]byte, 0),
		conn:      conn,
		lock:      new(sync.RWMutex),
		isClosed:  false,
		closeLock: new(sync.RWMutex),
	}
	go r.read()
	return r
}

// JTT格式数据协议
type JTTProtocol struct {
	buffer    []byte
	message   []byte
	conn      net.Conn
	lock      *sync.RWMutex
	isClosed  bool
	closeLock *sync.RWMutex
}

// 发送数据
func (r *JTTProtocol) Send(data []byte) (int, error) {
	return r.conn.Write(data)
}

// 读取数据
func (r *JTTProtocol) Parse() ([]byte, error) {
	if msg, ok := r.parse(); ok {
		logger.Info("receive message is:", msg)
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
		r.buffer = append(r.buffer, tmp[:size]...)
		r.lock.Unlock()
	}

}

// 解析是否获取一个完成的报文
func (r *JTTProtocol) parse() ([]byte, bool) {
	endFlag := false
	r.lock.RLock()
	tmp := r.buffer
	r.buffer = make([]byte, 0)
	r.lock.RUnlock()
	for idx, item := range tmp {
		if item == HEADER_FLAG {
			r.message = make([]byte, 0)
		}
		r.message = append(r.message, item)
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
	if endFlag && r.check() {
		return r.escapse(), true
	}
	return nil, false
}

func (r *JTTProtocol) escapse() []byte {
	buf := make([]byte, 0)
	buf = append(buf, r.message[0])
	escapseCount := 0
	msg := r.message[1 : len(r.message)-1]
	for idx := 0; idx < len(msg); idx++ {
		item := msg[idx]
		if item == ESCAPE_HEADER_FLAG && msg[idx+1] == ESCAPE_HEADER_FLAG_APPEND {
			buf = append(buf, HEADER_FLAG)
			escapseCount += 1
			idx += 1
		} else if item == ESCAPE_HEADER_FLAG && msg[idx+1] == ALIAS_HEADER_FLAG_5A_APPEND {
			buf = append(buf, ESCAPE_HEADER_FLAG)
			escapseCount += 1
			idx += 1
		} else if item == ESCAPE_FOOTER_FLAG && msg[idx+1] == ESCAPE_FOOTER_FLAG_APPEND {
			buf = append(buf, FOOTER_FLAG)
			escapseCount += 1
			idx += 1
		} else if item == ESCAPE_FOOTER_FLAG && msg[idx+1] == ESCAPE_FOOTER_FLAG_5E_APPEND {
			buf = append(buf, ESCAPE_FOOTER_FLAG)
			escapseCount += 1
			idx += 1
		} else {
			buf = append(buf, item)
		}
	}
	buf = append(buf, FOOTER_FLAG)
	if escapseCount > 0 {
		size := bytes.NewBuffer(make([]byte, 0))
		_ = binary.Write(size, binary.BigEndian, (uint32)(len(buf)))
		for idx, item := range buf[1:5] {
			buf[idx+1] = item
		}
	}
	return buf
}

func (r *JTTProtocol) check() bool {
	var size uint32 = 0
	sizeVal := r.message[1:5]
	_ = binary.Read(bytes.NewBuffer(sizeVal), binary.BigEndian, &size)
	if uint32(len(r.message)) == size {
		crcCode := uint16(0)
		_ = binary.Read(bytes.NewBuffer(r.message[size-3:size-1]), binary.BigEndian, &crcCode)
		crc, _ := CRC16CCITT(r.message[23 : size-3])
		if crcCode == crc {
			body := Decrypt(KEY, r.message[23:size-3])
			copy(r.message[23:size-3], body)
			return true
		}
		logger.Warn("CRC16-CCITT校验失败:", r.message)
		return false
	} else {
		logger.Info("报文字节字节数不一致。", r.message)
		return false
	}
}
