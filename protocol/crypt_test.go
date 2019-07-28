package protocol

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	val := Encrypt(KEY, []byte{0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39})
	fmt.Println(val)
	val2 := Decrypt(KEY, val)
	fmt.Println(val2)

	val3 := Decrypt(KEY, val2)
	fmt.Println(val3)

	val4 := Decrypt(KEY, val3)
	fmt.Println(val4)
}
