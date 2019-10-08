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
	tp, err := strconv.Atoi(device.DeviceTemperature)
	if err != nil && tp > warnTemperature {
		// 设备温度异常
		dbobj.Exec("update device_manage_info set device_status = 3 where serial_number = ? and delete_status = 0", device.SerialNumber)
		logger.Warn("温度异常，设置设别状态为温度异常，并生成告警信息")
		// 产生告警日志
		dao.AddAlarmEvent(device.SerialNumber, 1)
	}
}

func init() {
	var rst string
	_, err := dbobj.Exec("select item_value from sys_global_config where item_id = 3", dbobj.PackArgs(), &rst)
	if err == nil {
		t, err := strconv.Atoi(rst)
		if err != nil {
			warnTemperature = t
		}
	}
}
