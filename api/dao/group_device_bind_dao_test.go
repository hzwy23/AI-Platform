package dao_test

import (
	"fmt"
	"ai-platform/api/dao"
	"ai-platform/api/entity"
	"testing"
)

func TestNewGroupDeviceBindDaoImpl(t *testing.T) {
	r := dao.NewGroupDeviceBindDao()
	fmt.Println(r.FindAll())
}

func TestGroupDeviceBindDaoImpl_Insert(t *testing.T) {
	item := entity.GroupDeviceBind{
		Id:           4,
		GroupId:      1,
		DeviceId:     1,
		CreateBy:     "test",
		CreateDate:   "2019-06-01",
		UpdateBy:     "test",
		UpdateDate:   "2019-06-01",
		DeleteStatus: 0,
	}

	r := dao.NewGroupDeviceBindDao()
	fmt.Println(r.Insert(item))
}

func TestGroupDeviceBindDaoImpl_FindByDeviceId(t *testing.T) {
	r := dao.NewGroupDeviceBindDao()
	fmt.Println(r.FindById(3))
}

func TestGroupDeviceBindDaoImpl_FindByGroupId(t *testing.T) {
	r := dao.NewGroupDeviceBindDao()
	fmt.Println(r.FindByDeviceId(1))
}

func TestGroupDeviceBindDaoImpl_FindByGroupId2(t *testing.T) {
	r := dao.NewGroupDeviceBindDao()
	fmt.Println(r.FindByGroupId(2))
}

func TestGroupDeviceBindDaoImpl_Update(t *testing.T) {
	r := dao.NewGroupDeviceBindDao()
	fmt.Println(r.Update(2, 3))
}
