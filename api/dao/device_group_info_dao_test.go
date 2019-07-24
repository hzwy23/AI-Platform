package dao_test

import (
	"ai-platform/api/dao"
	"ai-platform/api/entity"
	"fmt"
	"testing"
)

func TestNewDeviceGroupInfoDao(t *testing.T) {
	obj := dao.NewDeviceGroupInfoDao()
	fmt.Println(obj.FindAll())
}

func TestDeviceGroupInfoDaoImpl_Insert(t *testing.T) {
	obj := dao.NewDeviceGroupInfoDao()
	item := entity.DeviceGroupInfo{
		GroupId:    2,
		GroupName:  "ceshi",
		CreateBy:   "test",
		CreateDate: "2019-01-01",
		UpdateBy:   "test",
		UpdateDate: "2019-01-01",
	}

	err := obj.Insert(item)
	if err != nil {
		fmt.Println(err)
	}
}

func TestDeviceGroupInfoDaoImpl_FindById(t *testing.T) {
	obj := dao.NewDeviceGroupInfoDao()
	fmt.Println(obj.FindById(2))
}

func TestDeviceGroupInfoDaoImpl_LogicDeleteById(t *testing.T) {
	obj := dao.NewDeviceGroupInfoDao()
	fmt.Println(obj.LogicDeleteById(2))
}
