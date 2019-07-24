package service

import (
	"ai-platform/server/platform"
	"fmt"
)

func broadcast(context *platform.Context) {
	fmt.Println("broadcast:", context.GetMessage().MsgBody, context.GetMsgId())
	context.Send(110, context.GetMessage().MsgBody)
}

func init() {
	platform.Register(uint16(110), broadcast)
}
