package utils

import (
	"encoding/json"
)

type DeviceControlData struct {
	SerialNumber  string `json:"client_CPUID"`
	LightMode     string `json:"client_AutoFunction"`
	AutoStartTime string `json:"AutoTimeStart"`
	AutoEndTime   string `json:"AutoTimeStop"`
}

// 更新灯光控制策略
func UpdateLightMode(serialNumber string,  cmd DeviceControlData) error {
	body, _ := json.Marshal(cmd)
	err := Command(0x8005, serialNumber, body)
	if err != nil {
		return err
	}
	return nil
}

