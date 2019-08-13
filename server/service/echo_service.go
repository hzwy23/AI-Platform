package service

import (
	"ai-platform/server/platform"
	"fmt"
)

func echo(context *platform.Context) (int, string) {
	fmt.Println("echo:", context.GetMessage().MsgBody, context.GetMsgId())
	context.Send(0x7fff, context.GetMessage().MsgBody)
	return 200, "Success"
}

func init() {
	platform.Register(0x7fff, echo)
}
