package service

import (
	"ai-platform/api/dao"
	"ai-platform/api/entity"
	"ai-platform/dbobj"
	"ai-platform/panda"
	"ai-platform/panda/logger"
	"ai-platform/server/platform"
	"encoding/json"
	"fmt"
)

var alarm = dao.NewEventAlarmInfoDao()

var device = dao.NewDeviceManageInfoDao()

type LampExceptionData struct {
	SerialNumber string `json:"client_CPUID"`
}

// 灯珠异常
func lampException(context *platform.Context) (int, string) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	var data LampExceptionData
	err := json.Unmarshal(context.GetMessage().MsgBody, &data)
	if err != nil {
		logger.Error(err)
		return 0, err.Error()
	}
	_, err = dbobj.Exec("update device_manage_info set device_attribute = 2 where delete_status = 0 and serial_number = ?", data.SerialNumber)
	if err != nil {
		return 500, err.Error()
	}
	return genLampException(data.SerialNumber)
}

func genLampException(serialNumber string) (int, string) {

	element, err := device.FindBySerialNumber(serialNumber)
	if err != nil {
		return 500, err.Error()
	}

	item := entity.EventAlarmInfo{
		EventTypeCd:       2,
		OccurrenceTime:    panda.CurTime(),
		SerialNumber:      serialNumber,
		DeviceName:        element.DeviceName,
		DeviceIp:          element.DeviceIp,
		DeviceAttribute:   element.DeviceAttribute,
		DeviceBrightness:  element.DeviceBrightness,
		DeviceTemperature: element.DeviceTemperature,
		HandleStatus:      0,
		DeleteStatus:      0,
	}
	_, err = alarm.Insert(item)
	if err != nil {
		return 501, err.Error()
	}
	return 200, "Ok"
}

func init()  {
	platform.Register(0x0006, lampException)
}
