package service

import (
	"ai-platform/panda/hret"
	"ai-platform/server/listen"
	"ai-platform/server/platform"
	"ai-platform/server/proto_data"
	"encoding/json"
	"fmt"
	"time"
)

// 接收广播消息
func broadcast(context *platform.Context) (int, string) {
	defer hret.RecvPanic()

	bd := &proto_data.DeviceInfo{}
	err := json.Unmarshal(context.GetMessage().MsgBody, bd)
	if err == nil {
		start := time.Now().Unix()
		online := &listen.OnlineDevice{
			// 设备号
			SerialNumber: bd.SerialNumber,
			// 最近刷新的时间戳
			RefreshTime: start,
			// 在线标示
			OnlineStatus: 1,
			// 固件版本
			FirmwareVersion: bd.FirmwareVersion,
			// 设备IP
			DeviceIp: bd.DeviceIp,
			// 掩码
			Mask: bd.Mask,
			// 网关地址
			GatewayAddr: bd.GatewayAddr,
			// 设备端口
			DevicePort: bd.DevicePort,
			// Mac地址
			MacAddr: bd.MacAddr,
		}
		listen.UpdateOnlineDevice(bd.SerialNumber, online)
	} else {
		fmt.Println(err.Error())
		return 50030, err.Error()
	}
	return 200, "Ok"
}

func init() {

	platform.Register(0x0000, broadcast)

}
