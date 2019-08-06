package test

import (
	"ai-platform/protocol"
	"ai-platform/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type DeviceControlData struct {
	SerialNumber  string `json:"client_CPUID"`
	LightMode     string `json:"client_AutoFunction"`
	AutoStartTime string `json:"AutoTimeStart"`
	AutoEndTime   string `json:"AutoTimeStop"`
}

func WriteDeviceControl() {
	conn, err := net.Dial("udp", "192.168.2.255:8900")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	cmd := utils.DeviceControlData{
		SerialNumber:  "1111111111111111",
		AutoStartTime: "09:00",
		AutoEndTime:   "18:00",
		LightMode: "All",
	}

	body, _ := json.Marshal(cmd)
	data ,err := protocol.Pack(0x0005, body)
	conn.Write(data)
	fmt.Println("write",err)
}