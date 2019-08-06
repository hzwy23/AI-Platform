package test

import (
	"ai-platform/protocol"
	"encoding/json"
	"fmt"
	"net"
)

type AsyncTimeData struct {
	SerialNumber string `json:"client_CPUID"`
	TimeZone string `json:"TimeZone"`
	Year string `json:"Year"`
	Month string `json:"Month"`
	Day string `json:"Day"`
	Hour string `json:"Hour"`
	Minute string `json:"Minute"`
	Second string `json:"Scond"`
	Week string `json:"Week"`
}

func WriteBroadcast() {
	conn, err := net.Dial("udp", "192.168.2.255:8900")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()


	msg := AsyncTimeData{
		SerialNumber: "1111111111",
		TimeZone:"GB",
		Year: "2019",
		Month: "09",
		Day: "31",
		Hour: "21",
		Minute: "12",
		Second: "42",
		Week: "2",
	}
	body, _ := json.Marshal(msg)
	data ,err := protocol.Pack(0x0002, body)
	conn.Write(data)
	fmt.Println("write",err)
}