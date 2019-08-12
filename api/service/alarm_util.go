package service

import (
	"ai-platform/api/dao"
	"ai-platform/api/entity"
	"ai-platform/dbobj"
	"ai-platform/panda"
	"ai-platform/panda/logger"
	"fmt"
)

var device = dao.NewDeviceManageInfoDao()
var alarm = dao.NewEventAlarmInfoDao()

func AddAlarmEvent(key string, eventTypeCd int) {
	element, err := device.FindBySerialNumber(key)
	if err != nil {
		fmt.Println(err)
		return
	} else if len(element.SerialNumber) == 0 {
		return
	}

	status := 1
	if eventTypeCd == 1 {
		status = 3
	} else if eventTypeCd == 2 {
		status = 4
	} else if eventTypeCd == 3 {
		status = 2
	}

	dbobj.Exec("update device_manage_info set device_status = ? where serial_number = ? and delete_status = 0", status, key)
	logger.Info("产生异常事件，设备号是：", key, ",事件类型是：", eventTypeCd)
	// 检查设备是否存在未处理的异常
	var cnt = 0
	err = dbobj.QueryForObject("select count(*) from event_alarm_info where delete_status = 0 and handle_status = 0 and serial_number = ? and event_type_cd = ?", dbobj.PackArgs(key, eventTypeCd), &cnt)
	if err == nil && cnt == 0 {
		item := entity.EventAlarmInfo{
			EventTypeCd:       eventTypeCd,
			OccurrenceTime:    panda.CurTime(),
			SerialNumber:      key,
			DeviceName:        element.DeviceName,
			DeviceIp:          element.DeviceIp,
			DeviceAttribute:   element.DeviceAttribute,
			DeviceBrightness:  element.DeviceBrightness,
			DeviceTemperature: element.DeviceTemperature,
			HandleStatus:      0,
			DeleteStatus:      0,
		}
		alarm.Insert(item)
	}
}
