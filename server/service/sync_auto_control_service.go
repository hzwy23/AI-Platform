package service

import (
	"ai-platform/dbobj"
	"ai-platform/panda/logger"
	"ai-platform/server/platform"
	"ai-platform/server/utils"
	"encoding/json"
	"fmt"
)

type DeviceControlData struct {
	SerialNumber  string `json:"client_CPUID"`
	LightMode     string `json:"client_AutoFunction"`
	AutoStartTime string `json:"AutoTimeStart"`
	AutoEndTime   string `json:"AutoTimeStop"`
}

func asyncAutoControl(context *platform.Context) (int, string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	var data DeviceControlData
	err := json.Unmarshal(context.GetMessage().MsgBody, &data)
	if err != nil {
		logger.Error(err)
		return 0, err.Error()
	}

	var cmd utils.DeviceControlData
	err = dbobj.QueryForStruct("select serial_number, light_mode, auto_start_time, auto_end_time from device_manage_info where delete_status = 0 and serial_number = ?", &cmd, data.SerialNumber)
	if err != nil {
		logger.Error(err)
		return 0, err.Error()
	}

	LightMode := cmd.LightMode
	if LightMode == "1" {
		cmd.LightMode = "CDS"
	} else if LightMode == "2" {
		cmd.LightMode = "Timer"
	} else if LightMode == "3" {
		cmd.LightMode = "All"
	}
	fmt.Println(cmd)
	err = utils.UpdateLightMode(cmd.SerialNumber, cmd)
	if err != nil {
		return 500,err.Error()
	}
	return 200, "OK"
}

func init() {
	platform.Register(0x0005, asyncAutoControl)
}
