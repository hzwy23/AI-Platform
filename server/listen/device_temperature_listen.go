package listen

import (
	"ai-platform/api/dao"
	"ai-platform/dbobj"
	"ai-platform/panda/logger"
	"ai-platform/server/proto_data"
	"strconv"
)

var warnTemperature = 80

func CheckTemperature(device proto_data.DeviceAttribute) {
	tp, err := strconv.ParseFloat(device.DeviceTemperature,64)
	if err == nil && int(tp) > warnTemperature {
		// 设备温度异常
		logger.Warn("温度异常，设置设别状态为温度异常，并生成告警信息")
		// 产生告警日志
		dao.AddAlarmEvent(device.SerialNumber, 1)
	} else {
		// 温度正常，删除温度报警提示
		// 取消温度异常警告
		dbobj.Exec("update device_manage_info set device_status = concat(substr(device_status,1,2), 1) where serial_number = ? and delete_status = 0", device.SerialNumber)
		alarm.ChangeHandleStatus(device.SerialNumber, 1)
	}
}

func UpdateWarnTemperature(temperature int)  {
	warnTemperature = temperature
}

func init() {
	var rst string
	err := dbobj.QueryForObject("select item_value from sys_global_config where item_id = 3", dbobj.PackArgs(), &rst)
	if err == nil {
		t, err := strconv.Atoi(rst)
		if err == nil {
			warnTemperature = t
		}
	}
}
