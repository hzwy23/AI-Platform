package protocol

import (
	"fmt"
	"testing"
)

func TestCRC16CCITT(t *testing.T) {
	crcval, crcstr := CRC16CCITT([]byte{0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39})
	fmt.Println(crcval, crcstr)
}
