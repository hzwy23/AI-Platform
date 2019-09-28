package dao_test

import (
	"ai-platform/api/dao"
	"ai-platform/api/entity"
	"fmt"
	"testing"
)


func TestDeviceManageInfoDaoImpl_LogicDeleteById(t *testing.T) {
	r := dao.NewDeviceManageInfoDao()
	fmt.Println(r.LogicDeleteById(1))
	fmt.Println(r.LogicDeleteById(2))
	fmt.Println(r.LogicDeleteById(3))
	fmt.Println(r.LogicDeleteById(4))
}

func TestDeviceManageInfoDaoImpl_Insert(t *testing.T) {
	r := dao.NewDeviceManageInfoDao()

	item := entity.DeviceManageInfo{
		DeviceId:             5,
		SerialNumber:         "PIN11111111",
		DeviceName:           "hello test",
		DhcpFlag:             1,
		DeviceIp:             "192.168.1.1",
		DevicePort:           "8989",
		DeviceStatus:         1,
		DeviceAttribute:      1,
		DevicePower:          30,
		DeviceLightThreshold: 1,
		DeviceBrightness:     1,
		DeviceTemperature:    15,
		AutoStartTime:        "09:15",
		AutoEndTime:          "15:00",
		LightMode:            1,
		MacAddress:           "01-FD-3F-3F-4F",
		FirmwareVersion:      "V0.0.1",
		Longitude:            "12.32232",
		Latitude:             "122.23321",
		Address:              "hellow rodl",
		Mask:                 "hello world",
		Gateway:              "192.168.2.2",
		Pin:                  "123456",
		CreateBy:             "test",
		CreateDate:           "2019-01-23",
		UpdateBy:             "test",
		UpdateData:           "2019-02-21",
		DeleteStatus:         0,
		PowerTotal:           100,
		StrobeCount:          100,
	}

	fmt.Println(r.Insert(item))
}
