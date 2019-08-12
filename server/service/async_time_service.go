package service

import (
	"ai-platform/panda/logger"
	"ai-platform/server/platform"
	"ai-platform/server/proto_data"
	"ai-platform/server/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// asyncTimeService 同步时间
// 设备向平台发送0x0002请求，平台给设备返回0x8002
func asyncTimeService(context *platform.Context) (int, string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	var data proto_data.AsyncTimeData
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
	err = utils.Command(0x8002, data.SerialNumber, body)
	if err != nil {
		return 400, err.Error()
	}
	return 0, "Ok"
}

func init() {
	platform.Register(0x0002, asyncTimeService)
}
