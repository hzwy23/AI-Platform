package listen

import (
	"ai-platform/api/dao"
	"ai-platform/dbobj"
	"ai-platform/panda/hret"
	"ai-platform/panda/logger"
	"time"
)

type OnlineDevice struct {
	// 设备号
	SerialNumber string
	// 最近刷新的时间戳
	RefreshTime int64
	// 在线标示
	OnlineStatus uint8
	// 固件版本
	FirmwareVersion string
	// 设备IP
	DeviceIp string
	// 掩码
	Mask string
	// 网关地址
	GatewayAddr string
	// 设备端口
	DevicePort string
	// Mac地址
	MacAddr string
}

var device dao.DeviceManageInfoDao
var alarm dao.EventAlarmInfoDao
var beatDuration int64

// 获取在线设备信息
func GetOnlineDevice() ([]OnlineDevice, error) {
	var rst []OnlineDevice
	err := dbobj.QueryForSlice("select serial_number, refresh_time, online_status, firmware_version, device_ip, mask, gateway_addr, device_port, mac_addr from device_scan_info order by device_ip", &rst)
	return rst, err
}

func GetOnlineItem(serialNumber string) (OnlineDevice, error) {
	var rst OnlineDevice
	err := dbobj.QueryForStruct("select serial_number, refresh_time, online_status, firmware_version, device_ip, mask, gateway_addr, device_port, mac_addr from device_scan_info where serial_number = ?", &rst, serialNumber)
	return rst, err
}

func UpdateOnlineDevice(key string, element *OnlineDevice) {
	var cnt = 0
	err := dbobj.QueryForObject("select count(*) from device_scan_info where serial_number = ?", dbobj.PackArgs(key), &cnt)
	if err != nil {
		logger.Warn("查询设备是否已经被扫描到，失败原因是：", err)
		return
	}

	if cnt > 0 {
		// 更新操作
		dbobj.Exec("update device_scan_info set refresh_time = ?,  online_status = 1, firmware_version = ?, device_ip = ?, mask = ?, gateway_addr = ?, device_port = ?, mac_addr = ? where serial_number = ?",
			element.RefreshTime, element.FirmwareVersion, element.DeviceIp, element.Mask, element.GatewayAddr, element.DevicePort, element.MacAddr, element.SerialNumber)
	} else {
		// 新增操作
		dbobj.Exec("insert into device_scan_info(serial_number, refresh_time, online_status, firmware_version, device_ip, mask, gateway_addr, device_port, mac_addr) values(?,?,?,?,?,?,?,?,?)",
			element.SerialNumber, element.RefreshTime, element.OnlineStatus, element.FirmwareVersion, element.DeviceIp, element.Mask, element.GatewayAddr, element.DevicePort, element.MacAddr)
	}

	// 设置设备上线
	dbobj.Exec("update device_manage_info set device_status = 1 where serial_number = ? and delete_status = 0", key)

	// 取消离线告警
	alarm.ChangeHandleStatus(key, 2)
}

// 定时清理离线设备，并发送告警信息
func removeOfflineDevice() {
	// 如果程序异常退出，重新拉起
	defer hret.RecvPanic()

	for {
		rst, err := GetOnlineDevice()
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		for _, val := range rst {
			// 设备持续掉线30s将会判定为设备离线
			duration := time.Now().Unix() - val.RefreshTime
			logger.Info(val.SerialNumber, ',',duration,',',beatDuration,",", time.Now().Unix(),',', val.RefreshTime)
			if duration > beatDuration {
				logger.Info("从设备扫描列表中删除设备", val.SerialNumber)
				dbobj.Exec("delete from device_scan_info where serial_number = ?", val.SerialNumber)
				dao.AddAlarmEvent(val.SerialNumber, 2)
			}
		}
		logger.Info("sync device status")
		time.Sleep(time.Second * 2)
	}
}

func checkAddedDevice()  {
	defer hret.RecvPanic()
	for {
		rst, err := device.FindAll("")
		if err != nil {
			logger.Error(err)
			time.Sleep(time.Second*5)
		}
		for _, item := range rst {
			var cnt = 0
			err = dbobj.QueryForObject("select count(*) from device_scan_info where serial_number = ?", dbobj.PackArgs(item.SerialNumber), &cnt)
			if err != nil || cnt == 0 {
				dao.AddAlarmEvent(item.SerialNumber, 2)			}
		}
		time.Sleep(time.Second*10)
	}

}

func init() {

	defer hret.RecvPanic()

	alarm = dao.NewEventAlarmInfoDao()

	device = dao.NewDeviceManageInfoDao()

	err := dbobj.QueryForObject("select item_value from sys_global_config where item_id = 4",dbobj.PackArgs(), &beatDuration)
	if err != nil {
		logger.Warn("获取设备心跳时长失败. 失败原因是：", err)
		beatDuration = 60
	}
	logger.Info("设备心跳保持时间是：", beatDuration)

	go removeOfflineDevice()

	go checkAddedDevice()

}
