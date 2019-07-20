package protocol

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	val := Encrypt(123, []byte{0x3c, 0x1f, 0x2a, 0x8b})
	fmt.Println(val)
	val2 := Decrypt(123, val)
	fmt.Println(val2)
	val3 := Encrypt(123, val2)
	fmt.Println(val3)
	val4 := Decrypt(123, val3)
	fmt.Println(val4)
}
