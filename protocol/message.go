package protocol

import (
	"ai-platform/panda/logger"
	"bytes"
	"encoding/binary"
	"errors"
)

// Message 报文数据
type Message struct {
	// 协议头标识
	HeaderFlag byte

	// 协议头
	MsgHeader Header

	// 数据内容
	MsgBody []byte

	// CRC校验码
	CrcCode uint16

	// 协议尾标识
	FooterFlag byte
}

// ConvertMessage 将字节转换成报文对象
func ConvertMessage(data []byte) *Message {
	size := len(data)
	if size < 26 {
		return nil
	}
	var message Message
	binary.Read(bytes.NewBuffer(data[:1]), binary.BigEndian, &message.HeaderFlag)
	binary.Read(bytes.NewBuffer(data[1:23]), binary.BigEndian, &message.MsgHeader)
	message.MsgBody = make([]byte, size-26)
	binary.Read(bytes.NewBuffer(data[23:size-3]), binary.BigEndian, &message.MsgBody)
	binary.Read(bytes.NewBuffer(data[size-3:size-1]), binary.BigEndian, &message.CrcCode)
	binary.Read(bytes.NewBuffer(data[size-1:]), binary.BigEndian, &message.FooterFlag)
	logger.Debug("接收到的对象是：", message)
	return &message
}

// ConvertToByte 将报文对象转换成字节
func ConvertToByte(message *Message) []byte {
	sbuf := bytes.NewBuffer(make([]byte, 0))
	ebuf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(ebuf, binary.BigEndian, message.MsgBody)
	binary.Write(ebuf, binary.BigEndian, message.CrcCode)
	binary.Write(ebuf, binary.BigEndian, message.FooterFlag)
	eb := ebuf.Bytes()
	message.MsgHeader.MsgLength = (uint32)(len(eb) + 23)
	binary.Write(sbuf, binary.BigEndian, message.HeaderFlag)
	binary.Write(sbuf, binary.BigEndian, message.MsgHeader)
	sb := sbuf.Bytes()
	result := append(sb, eb...)
	logger.Debug("发送：", result)
	return result
}


// UnPack 解包
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

// Pack 封包
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
			escapseCount++
			idx++
		} else if item == ESCAPE_HEADER_FLAG && msg[idx+1] == ALIAS_HEADER_FLAG_5A_APPEND {
			buf = append(buf, ESCAPE_HEADER_FLAG)
			escapseCount++
			idx++
		} else if item == ESCAPE_FOOTER_FLAG && msg[idx+1] == ESCAPE_FOOTER_FLAG_APPEND {
			buf = append(buf, FOOTER_FLAG)
			escapseCount++
			idx++
		} else if item == ESCAPE_FOOTER_FLAG && msg[idx+1] == ESCAPE_FOOTER_FLAG_5E_APPEND {
			buf = append(buf, ESCAPE_FOOTER_FLAG)
			escapseCount++
			idx++
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
	var size uint32
	sizeVal := realMsg[1:5]
	_ = binary.Read(bytes.NewBuffer(sizeVal), binary.BigEndian, &size)

	if uint32(len(realMsg)) == size {
		return size, true
	} 

	logger.Warn("报文字节字节数不一致。", realMsg)
	return 0, false
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

func transfer(message []byte) ([]byte, bool){
	// 转义
	realMsg := unpackescapse(message)
	// 检查长度
	size, flag := check(realMsg)
	// 解密
	realMsg = decrypt(realMsg, size)
	if flag {
		// CRC校验
		if crc(realMsg, size) {
			return realMsg, true
		}
		return nil, false
	}
	return nil, false
}
