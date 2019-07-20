package protocol

import (
	"ai-platform/panda/logger"
	"bytes"
	"encoding/binary"
)

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
	return &message
}

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
	logger.Info("发送：", result)
	return result
}
