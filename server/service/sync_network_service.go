package service

import (
	"ai-platform/api/dao"
	"ai-platform/api/entity"
	"ai-platform/dbobj"
	"ai-platform/panda"
	"ai-platform/panda/hret"
	"ai-platform/panda/logger"
	"ai-platform/server/platform"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type onlineDevice struct {
	// 设备号
	SerialNumber string
	// 最近刷新的时间戳
	RefreshTime int64
	// 设备信息
	*DeviceInfo
}

type DeviceInfo struct {
	// 设别序列号
	SerialNumber string `json:"client_CPUID"`
	// 软件版本号
	FirmwareVersion string `json:"client_FrameworkVersion"`
	// 设备IP地址
	DeviceIp string `json:"client_IP"`
	// 设备掩码
	Mask string `json:"client_MASK"`
	// 网关地址
	GatewayAddr string `json:"client_GATEWAY"`
	// 设备端口号
	DevicePort string `json:"client_PORT"`
	// 设备mac地址
	MacAddr string `json:"client_MAC"`
}

var deviceScan = make(map[string]*onlineDevice, 0)
var lock = &sync.RWMutex{}
var alarm dao.EventAlarmInfoDao
var device dao.DeviceManageInfoDao

func GetOnlineDevice() map[string]*onlineDevice {
	lock.RLock()
	defer lock.RUnlock()
	return deviceScan
}

// 定时清理离线设备，并发送告警信息
func removeOfflineDevice() {
	// 如果程序异常退出，重新拉起
	defer hret.RecvPanic(removeOfflineDevice)

	for {
		for key, val := range deviceScan {
			// 设备持续掉线30s将会判定为设备离线
			duration := time.Now().Unix() - val.RefreshTime
			if duration > 30 {
				logger.Info("从设备扫描列表中删除设备", key)
				lock.Lock()
				delete(deviceScan, key)
				lock.Unlock()
				element, err := device.FindBySerialNumber(key)
				if err != nil {
					fmt.Println(err)
					continue
				} else if len(element.SerialNumber) == 0 {
					continue
				}

				dbobj.Exec("update device_manage_info set device_status = 4 where serial_number = ? and delete_status = 0", key)
				item := entity.EventAlarmInfo{
					EventTypeCd:       2,
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

		time.Sleep(time.Millisecond * 500)
		logger.Debug("delete offline device")
	}
}

// 接收广播消息
func broadcast(context *platform.Context) (int, string) {
	defer hret.RecvPanic()

	bd := &DeviceInfo{}
	err := json.Unmarshal(context.GetMessage().MsgBody, bd)
	if err == nil {
		start := time.Now().Unix()
		online := &onlineDevice{
			bd.SerialNumber,
			start,
			bd,
		}
		lock.Lock()
		deviceScan[bd.SerialNumber] = online
		lock.Unlock()

		// 设备上线
		dbobj.Exec("update device_manage_info set device_status = 1 where serial_number = ? and delete_status = 0", bd.SerialNumber)

		// 取消告警
		alarm.ChangeHandleStatus(bd.SerialNumber, 2)

	} else {
		fmt.Println(err.Error())
		return 50030, err.Error()
	}
	return 200, "Ok"
}

func init() {

	alarm = dao.NewEventAlarmInfoDao()

	device = dao.NewDeviceManageInfoDao()

	platform.Register(0x0000, broadcast)

	go removeOfflineDevice()

}
