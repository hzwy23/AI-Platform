package service

import (
	"ai-platform/panda/logger"
	"ai-platform/server/platform"
	"ai-platform/server/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type AsyncTimeData struct {
	SerialNumber string `json:"client_CPUID"`
	TimeZone string `json:"TimeZone"`
	Year string `json:"Year"`
	Month string `json:"Month"`
	Day string `json:"Day"`
	Hour string `json:"Hour"`
	Minute string `json:"Minute"`
	Second string `json:"Second"`
	Week string `json:"Week"`
}

func asyncTimeService(context *platform.Context)  (int, string)  {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	var data AsyncTimeData
	err := json.Unmarshal(context.GetMessage().MsgBody, &data)
	if err != nil {
		logger.Error(err)
		return 0, err.Error()
	}
	now := time.Now()
	data.Year = strconv.Itoa(now.Year())
	data.Month = strconv.Itoa(int(now.Month()))
	data.Day = strconv.Itoa(now.Day())
	data.Hour = strconv.Itoa(now.Hour())
	data.Minute = strconv.Itoa(now.Minute())
	data.Second = strconv.Itoa(now.Second())
	data.Week = strconv.Itoa(int(now.Weekday()))
	data.TimeZone, _ = now.Zone()
	body, _ := json.Marshal(data)
	err = utils.Command(context.GetMsgId(),data.SerialNumber, body)
	if err != nil {
		return 400, err.Error()
	}
	return 0, "Ok"
}

func init()  {
	platform.Register(0x0002,asyncTimeService)
}
