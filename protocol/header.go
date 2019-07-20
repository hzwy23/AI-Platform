package protocol

import (
	"bytes"
	"encoding/binary"
)

// 数据协议头
type Header struct {
	// 报文总长度
	MsgLength uint32
	// 报文序列号
	MsgSN uint32
	// 业务类型
	MsgID uint16
	//
	MsgCenterID uint32
	// 版本好
	VersionFlag [3]byte
	// 加密标志
	EncryptFlag byte
	// 密钥
	EncryptKey uint32
}

func (r *Header) toByte() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.BigEndian, r)
	return buf.Bytes()
}
