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
	"sync"
	"time"
)


type onlineDevice struct {
	SerialNumber string
	// 最近刷新的时间戳
	RefreshTime time.Time
	*DeviceInfo
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

func removeOfflineDevice()  {
	for key, val := range deviceScan {
		d,_:=time.ParseDuration("30s")
		if time.Now().After(val.RefreshTime.Add(d)) {
			logger.Info("从设备扫描列表中删除设备", key)
			lock.Lock()
			delete(deviceScan, key)
			dbobj.Exec("update device_manage_info set device_status = 4 where serial_number = ? and delete_status = 0", key)
			item := entity.EventAlarmInfo{
				EventTypeCd:       2,
				OccurrenceTime:    panda.CurTime(),
				SerialNumber:      key,
				DeviceName:        "设备不存在",
				DeviceIp:          "-",
				DeviceAttribute:   0,
				DeviceBrightness:  0,
				DeviceTemperature: "0",
				HandleStatus:      0,
				DeleteStatus:      0,
			}
			// todo 生成离线异常信息
			element, err := device.FindBySerialNumber(key)
			if err != nil {
				fmt.Println(err)
			} else {
				item = entity.EventAlarmInfo{
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
			}
			alarm.Insert(item)
			lock.Unlock()
		}

	}
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

func broadcast(context *platform.Context) (int, string){
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	bd := &DeviceInfo{}
	err := json.Unmarshal(context.GetMessage().MsgBody, bd)
	if err == nil {
		online := &onlineDevice{
			bd.SerialNumber,
			time.Now(),
			bd,
		}
		lock.Lock()
		deviceScan[bd.SerialNumber] = online
		// 设备上线
		dbobj.Exec("update device_manage_info set device_status = 1 where serial_number = ? and delete_status = 0", bd.SerialNumber)
		// 取消告警
		alarm.ChangeHandleStatus(bd.SerialNumber, 2)
		lock.Unlock()
	}
	return 200, "Ok"
}

func init() {
	alarm = dao.NewEventAlarmInfoDao()
	device = dao.NewDeviceManageInfoDao()

	platform.Register(0x0000, broadcast)
	go func() {
		for {
			removeOfflineDevice()
			time.Sleep(time.Millisecond*500)
			logger.Info("delete offline device")
		}
	}()
}
