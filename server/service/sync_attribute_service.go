package service

import (
	"ai-platform/dbobj"
	"ai-platform/panda/hret"
	"ai-platform/panda/logger"
	"ai-platform/server/listen"
	"ai-platform/server/platform"
	"ai-platform/server/proto_data"
	"encoding/json"
)

func asyncAttribute(context *platform.Context) (int, string) {
	defer hret.RecvPanic()
	var data proto_data.DeviceAttribute
	err := json.Unmarshal(context.GetMessage().MsgBody, &data)
	if err != nil {
		logger.Error(err)
		return 0, err.Error()
	}

	// 检测温度是否异常
	listen.CheckTemperature(data)

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

	thresHold, ok := proto_data.ThresholdValueReserve[data.DeviceLightThreshold]
	if !ok {
		thresHold = "1"
	}
	brightness, ok := proto_data.LightValueReserve[data.DeviceBrightness]
	if !ok {
		brightness = "1"
	}
	_, err = dbobj.Exec("update device_manage_info set device_attribute = ?, device_power = ?, device_temperature = ?, device_light_threshold = ?, device_brightness = ?, power_total = ?, strobe_count = ? where serial_number = ? and delete_status = 0",
		attr, data.DevicePower, data.DeviceTemperature,
		thresHold, brightness,
		data.PowerTotal, data.StrobeCount, data.SerialNumber)
	logger.Info("同步设备属性，设备同步信息是：", data)
	if err != nil {
		return 500, err.Error()
	}
	return 200, "OK"
}

func init() {
	platform.Register(0x0001, asyncAttribute)
}
