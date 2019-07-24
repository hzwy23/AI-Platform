package dao_test

import (
	"ai-platform/api/dao"
	"ai-platform/api/entity"
	"fmt"
	"testing"
)

func TestNewEventAlarmInfoDaoImpl(t *testing.T) {
	r := dao.NewEventAlarmInfoDao()
	fmt.Println(r.FindAll())
}

func TestEventAlarmInfoDaoImpl_Insert(t *testing.T) {
	item := entity.EventAlarmInfo{
		Id:                3,
		EventTypeCd:       1,
		OccurrenceTime:    "2019-01-01",
		DeviceId:          "123",
		DeviceName:        "hello world",
		DeviceIp:          "192.168.2.1",
		DeviceAttribute:   1,
		DeviceBrightness:  1,
		DeviceTemperature: 30,
		HandleStatus:      0,
		DeleteStatus:      0,
	}
	r := dao.NewEventAlarmInfoDao()
	fmt.Println(r.Insert(item))

}

func TestEventAlarmInfoDaoImpl_CloseById(t *testing.T) {
	r := dao.NewEventAlarmInfoDao()
	fmt.Println(r.CloseById(1, 1))
}

func TestEventAlarmInfoDaoImpl_FindById(t *testing.T) {
	r := dao.NewEventAlarmInfoDao()
	fmt.Println(r.FindById(1))
}

func TestEventAlarmInfoDaoImpl_FindByTypeCd(t *testing.T) {
	r := dao.NewEventAlarmInfoDao()
	fmt.Println(r.FindByTypeCd(1))
}

func TestEventAlarmInfoDaoImpl_LogicDeleteById(t *testing.T) {
	r := dao.NewEventAlarmInfoDao()
	fmt.Println(r.LogicDeleteById(2))
}
