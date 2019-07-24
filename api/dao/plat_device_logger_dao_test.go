package dao_test

import (
	"ai-platform/api/dao"
	"ai-platform/api/entity"
	"fmt"
	"testing"
)

func TestNewPlatDeviceLoggerDao(t *testing.T) {
	r := dao.NewPlatDeviceLoggerDao()
	fmt.Println(r.FindAll(1, 10))
	fmt.Println(r.FindAll(2, 10))
}

func TestPlatDeviceLoggerDaoImpl_Insert(t *testing.T) {
	item := entity.PlatDeviceLogger{
		Id:         11,
		DeviceId:   1,
		HandleTime: "2019-01-01",
		DeviceName: "test",
		Direction:  "Input",
		BizType:    "1",
		Message:    "hello",
		RetCode:    "200",
		RetMsg:     "message is message",
	}
	r := dao.NewPlatDeviceLoggerDao()
	fmt.Println(r.Insert(item))
}
