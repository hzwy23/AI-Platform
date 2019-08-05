package panda_test

import (
	"ai-platform/panda"
	"fmt"
	"testing"
)

func TestGetKey(t *testing.T) {
	s := panda.JoinKey("hello world", "my name is", "c huang zhan wei")
	fmt.Println(s)
	fmt.Println(panda.GetKey(s, 4))
	fmt.Println([]byte(s))
	fmt.Println(byte(0x1f))
}
