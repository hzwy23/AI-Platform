package service

import (
	"ai-platform/server/platform"
	"fmt"
)

func echo(context *platform.Context) (int, string){
	fmt.Println("echo:", context.GetMessage().MsgBody, context.GetMsgId())
	context.Send(888, context.GetMessage().MsgBody)
	return 200,"OK"
}

func init() {
	platform.Register(uint16(888), echo)
}
