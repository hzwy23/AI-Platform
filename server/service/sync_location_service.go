package service

import (
	"ai-platform/server/platform"
	"fmt"
)

func asyncLocation(context *platform.Context) (int, string){
	fmt.Println("echo:", context.GetMessage().MsgBody, context.GetMsgId())
	context.Send(0x0003, context.GetMessage().MsgBody)
	return 200,"OK"
}

func init() {
	platform.Register(0x0003, asyncLocation)
}

