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

// 数据解析
// 1. 转义处理
// 2. 长度校验
// 3. 解密body部分
// 4. CRC校验body部分

// 封装数据
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
		isClosed:  false,
		closeLock: new(sync.RWMutex),
	}
	go r.read()
	return r
}

func NewJTTProtocolUDP(conn *net.UDPConn) *JTTProtocol {
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


func NewUDPJTTProtocol(conn *net.UDPConn) *JTTProtocol {
	r := &JTTProtocol{
		buffer:    make([]byte, 0),
		message:   make([]byte, 0),
		conn:      conn,
		lock:      new(sync.RWMutex),
		isClosed:  false,
		closeLock: new(sync.RWMutex),
	}
	go r.readFromUdp(conn)
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
func (r *JTTProtocol) Send(msgId uint16, msgData []byte) (int, error) {
	data, err := Pack(msgId, msgData)
	if err != nil {
		logger.Warn(err)
		return 0, err
	}
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

func UnPack(msgData []byte) ([]byte, error) {
	if len(msgData) < 26 {
		logger.Error(msgData)
		return nil, errors.New("接收数据报文信息不正确")
	}
	// 转义
	realMsg := unpackescapse(msgData)
	// 检查长度
	size, flag := check(realMsg)
	// 解密
	realMsg = decrypt(realMsg, size)
	if flag {
		// CRC校验
		if crc(realMsg, size) {
			return realMsg, nil
		} else {
			return nil, errors.New("CRC校验失败")
		}
	} else {
		return nil, errors.New("报文长度校验失败")
	}
}

// 封包
func Pack(msgID uint16, body []byte) ([]byte, error) {
	header := Header{
		MsgLength:   0,
		MsgSN:       1,
		MsgID:       msgID,
		MsgCenterID: 1,
		VersionFlag: [3]byte{0x01, 0x01, 0x01},
		EncryptFlag: 0x01,
		EncryptKey:  KEY,
	}
	crc, _ := CRC16CCITT(body)

	encryptBody := Encrypt(KEY, body)

	message := &Message{
		HeaderFlag: HEADER_FLAG,
		MsgHeader:  header,
		MsgBody:    encryptBody,
		CrcCode:    crc,
		FooterFlag: FOOTER_FLAG,
	}

	data := ConvertToByte(message)
	msgLen := len(data)
	// 修改报文长度字段
	_ = binary.Write(bytes.NewBuffer(data[1:5]), binary.BigEndian, msgLen)

	msgData, err := packEscape(data)

	return msgData, err
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
		logger.Debug("receive byte is: ", tmp[:size])
		r.buffer = append(r.buffer, tmp[:size]...)
		r.lock.Unlock()
	}
}

func (r *JTTProtocol) readFromUdp(conn *net.UDPConn) {
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
	if endFlag {
		// 转义
		realMsg := unpackescapse(r.message)
		// 检查长度
		size, flag := check(realMsg)
		// 解密
		realMsg = decrypt(realMsg, size)
		if flag {
			// CRC校验
			if crc(realMsg, size) {
				return realMsg, true
			} else {
				return nil, false
			}
		} else {
			return nil, false
		}
	}
	return nil, false
}

func packEscape(msgData []byte) ([]byte, error) {
	if len(msgData) < 23 {
		logger.Warn("无效的报文")
		return nil, errors.New("无效的报文")
	}
	buf := make([]byte, 0)
	buf = append(buf, msgData[0])
	for _, val := range msgData[1 : len(msgData)-1] {
		if val == HEADER_FLAG {
			buf = append(buf, ESCAPE_HEADER_FLAG, ESCAPE_HEADER_FLAG_APPEND)
		} else if val == ESCAPE_HEADER_FLAG {
			buf = append(buf, ESCAPE_HEADER_FLAG, ALIAS_HEADER_FLAG_5A_APPEND)
		} else if val == FOOTER_FLAG {
			buf = append(buf, ESCAPE_FOOTER_FLAG, ESCAPE_FOOTER_FLAG_APPEND)
		} else if val == ESCAPE_FOOTER_FLAG {
			buf = append(buf, ESCAPE_FOOTER_FLAG, ESCAPE_FOOTER_FLAG_5E_APPEND)
		} else {
			buf = append(buf, val)
		}
	}
	buf = append(buf, FOOTER_FLAG)
	return buf, nil
}

// 转义处理
func unpackescapse(msgData []byte) []byte {
	if len(msgData) == 0 {
		return nil
	}
	buf := make([]byte, 0)
	buf = append(buf, msgData[0])
	escapseCount := 0
	msg := msgData[1 : len(msgData)-1]
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

// 长度校验
func check(realMsg []byte) (uint32, bool) {
	if len(realMsg) < 26 {
		return 0, false
	}
	// 获取数据包长度
	var size uint32 = 0
	sizeVal := realMsg[1:5]
	_ = binary.Read(bytes.NewBuffer(sizeVal), binary.BigEndian, &size)

	if uint32(len(realMsg)) == size {
		return size, true
	} else {
		logger.Info("报文字节字节数不一致。", realMsg)
		return 0, false
	}
}

// CRC 校验
func crc(realMsg []byte, size uint32) bool {
	// CRC校验
	crcCode := uint16(0)
	_ = binary.Read(bytes.NewBuffer(realMsg[size-3:size-1]), binary.BigEndian, &crcCode)
	crc, _ := CRC16CCITT(realMsg[23 : size-3])

	if crcCode == crc {
		return true
	}

	logger.Warn("CRC16-CCITT校验失败:", realMsg)
	return false
}

func decrypt(realMsg []byte, size uint32) []byte {
	body := Decrypt(KEY, realMsg[23:size-3])
	copy(realMsg[23:size-3], body)
	return realMsg
}
