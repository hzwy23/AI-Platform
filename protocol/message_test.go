package protocol_test

import (
	"ai-platform/protocol"
	"fmt"
	"testing"
	"unsafe"
)

func TestConvertToByte(t *testing.T) {

	header := protocol.Header{
		MsgSN:       12122,
		MsgID:       12,
		MsgCenterID: 34,
		VersionFlag: [3]byte{0x01, 0x02, 0x03},
		EncryptFlag: 0x01,
		EncryptKey:  45,
	}

	fmt.Println(int(unsafe.Sizeof(protocol.Header{})))

	message := &protocol.Message{
		HeaderFlag: protocol.HEADER_FLAG,
		MsgHeader:  header,
		MsgBody:    []byte("abc"),
		CrcCode:    12,
		FooterFlag: protocol.FOOTER_FLAG,
	}

	buf := protocol.ConvertToByte(message)
	fmt.Println(buf)

	st := protocol.ConvertMessage(buf)
	fmt.Println(st, st.MsgHeader, st.MsgBody, st.CrcCode, st.HeaderFlag, st.FooterFlag)
	fmt.Println(st.MsgHeader)
	fmt.Println(buf)
}
