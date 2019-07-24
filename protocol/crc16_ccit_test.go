package protocol

import (
	"fmt"
	"testing"
)

func TestCRC16CCITT(t *testing.T) {
	crcval, crcstr := CRC16CCITT([]byte{48, 49, 50, 51, 52, 53, 54, 55, 56})
	fmt.Println(crcval, crcstr)
}
