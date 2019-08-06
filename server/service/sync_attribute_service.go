package service

import (
	"ai-platform/dbobj"
	"ai-platform/panda/logger"
	"ai-platform/server/platform"
	"encoding/json"
	"fmt"
)

type DeviceAttribute struct {
	// 序列号
	SerialNumber string `json:"client_CPUID"`
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
	// 设备属性  <Option value="1">常亮</Option>
	//         <Option value="2">频闪</Option>
	//         <Option value="3">爆灯</Option>
	DeviceAttribute string `json:"client_Mode"`
}

func asyncAttribute(context *platform.Context) (int, string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	var data DeviceAttribute
	err := json.Unmarshal(context.GetMessage().MsgBody, &data)
	if err != nil {
		logger.Error(err)
		return 0, err.Error()
	}

	attr := 1
	if data.DeviceAttribute == "Auto" {
		attr = 4
	} else if data.DeviceAttribute == "ON" {
		attr = 1
	} else if data.DeviceAttribute == "OFF" {
		attr = 2
	} else if data.DeviceAttribute == "Flash" {
		attr = 3
	}

	_, err = dbobj.Exec("update device_manage_info set device_attribute = ?, device_power = ?, device_temperature = ?, device_light_threshold = ?, device_brightness = ?, power_total = ?, strobe_count = ? where serial_number = ? and delete_status = 0",
		attr, data.DevicePower, data.DeviceTemperature, data.DeviceLightThreshold, data.DeviceBrightness, data.PowerTotal, data.StrobeCount, data.SerialNumber)
	if err != nil {
		return 500, err.Error()
	}
	return 200, "OK"
}

func init() {
	platform.Register(0x0001, asyncAttribute)
}
