package service

import (
	"ai-platform/dbobj"
	"ai-platform/panda/logger"
	"ai-platform/server/platform"
	"encoding/json"
	"fmt"
)

type deviceAttribute struct {
	// 序列号
	SerialNumber string `json:"client_CPUID"`
	// 是否DHCP
	DhcpFlag string `json:"client_Mode"`
	// 功率
	DevicePower string `json:"client_Power"`
	// 温度
	DeviceTemperature string `json:"client_Temp"`
	// 光敏阀值
	DeviceLightThreshold string `json:"client_CDSThreshold"`
	// 设备亮度
	DeviceBrightness string `json:"client_LightLevel"`
	// 总功耗
	PowerTotal string `json:"client_Consumption"`
	// 爆闪次数
	StrobeCount string `json:"client_FlashCount"`
	// 光照强度
	IntensityLight string `json:"client_CDS"`
}

func asyncAttribute(context *platform.Context) (int, string){
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	var data *deviceAttribute
	err := json.Unmarshal(context.GetMessage().MsgBody, data)
	if err != nil {
		logger.Error(err)
		return 0, err.Error()
	}
	fmt.Println(data)
	dbobj.Exec("update device_manage_info set dhcp_flag = ?, device_power = ?, device_temperature = ?, device_light_threshold = ?, device_brightness = ?, power_total = ?, strobe_count = ? where serial_number = ? and delete_status = 0",
		data.DhcpFlag, data.DevicePower, data.DeviceTemperature, data.DeviceLightThreshold, data.DeviceBrightness, data.PowerTotal, data.StrobeCount, data.SerialNumber)
	context.Send(0x0001, context.GetMessage().MsgBody)
	return 200,"OK"
}

func init() {
	platform.Register(0x0001, asyncAttribute)
}

