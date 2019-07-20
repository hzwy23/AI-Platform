package protocol

import (
	"fmt"
	"testing"
)

func TestCRC16CCITT(t *testing.T) {
	crcval, crcstr := CRC16CCITT([]byte{0x31, 0x32, 0x33, 0x34})
	fmt.Println(crcval, crcstr)
}
