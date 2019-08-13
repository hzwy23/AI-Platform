package service

import (
	"ai-platform/api/dao"
	"ai-platform/api/service"
	"ai-platform/dbobj"
	"ai-platform/panda/logger"
	"ai-platform/server/platform"
	"ai-platform/server/proto_data"
	"encoding/json"
	"fmt"
)

var alarm = dao.NewEventAlarmInfoDao()

var device = dao.NewDeviceManageInfoDao()

// 灯珠异常
func lampException(context *platform.Context) (int, string) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	var data proto_data.LampExceptionData
	err := json.Unmarshal(context.GetMessage().MsgBody, &data)
	if err != nil {
		logger.Error(err)
		return 0, err.Error()
	}
	_, err = dbobj.Exec("update device_manage_info set device_attribute = 2 where delete_status = 0 and serial_number = ?", data.SerialNumber)
	logger.Info("生成灯珠异常告警信息,", data)
	if err != nil {
		return 500, err.Error()
	}
	service.AddAlarmEvent(data.SerialNumber, 3)
	return 200, "Success"
}

func init() {
	platform.Register(0x0006, lampException)
}
