package service

import (
	"fmt"
	"ai-platform/server/platform"
)

func demo(context *platform.Context) {
	fmt.Println("hello world", context.GetMessage().MsgBody, context.GetMsgId())
	context.Send(context.GetMessage())
}

func init() {
	platform.Register(uint16(12), demo)
}
