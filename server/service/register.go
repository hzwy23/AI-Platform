package service

import (
	"fmt"
	"ai-platform/server/platform"
)

func demo(context *platform.Context) {
	fmt.Println("hello world", context.GetMessage().MsgBody, context.GetMsgId())
	context.Send(10, []byte{'h','e','l','l','o'})
}

func init() {
	platform.Register(uint16(12), demo)
}
