package dao

import (
	"ai-platform/api/entity"
	"errors"
)

type DeviceManageService interface {
	FindAll(groupId string) ([]entity.DeviceManageInfo, error)
	FindBySerialNumber(serialNumber string) (entity.DeviceManageInfo, error)
	AddDevice(item *entity.DeviceManageInfo, GroupId string) error
	UpdateDevice(item entity.DeviceManageInfo) error
	RemoveDevice(deviceId int) error
	FindByDeviceId(deviceId int) (entity.DeviceManageInfo, error)
}

type deviceManageServiceImpl struct {
	dao DeviceManageInfoDao
}

func (r *deviceManageServiceImpl) FindAll(groupId string) ([]entity.DeviceManageInfo, error) {
	return r.dao.FindAll(groupId)
}

func (r *deviceManageServiceImpl)FindByDeviceId(deviceId int) (entity.DeviceManageInfo, error)  {
	return r.dao.FindByDeviceId(deviceId);
}

func (r *deviceManageServiceImpl) FindBySerialNumber(serialNumber string) (entity.DeviceManageInfo, error) {
	return r.dao.FindBySerialNumber(serialNumber)
}

func (r *deviceManageServiceImpl) AddDevice(item *entity.DeviceManageInfo, GroupId string) error {
	size, err := r.dao.Insert(*item, GroupId)
	if err != nil {
		return err
	}
	if size > 0 {
		return nil
	}
	return errors.New("添加设备失败，请联系管理员")
}

func (r *deviceManageServiceImpl) UpdateDevice(item entity.DeviceManageInfo) error {
	panic("implement me")
}

func (r *deviceManageServiceImpl) RemoveDevice(deviceId int) error {
	size, err := r.dao.LogicDeleteById(deviceId)
	if err != nil {
		return err
	}
	if size > 0 {
		return nil
	}
	return errors.New("设备已经被删除或不存在")
}

func NewDeviceManageService() DeviceManageService {
	r := &deviceManageServiceImpl{
		dao: NewDeviceManageInfoDao(),
	}
	return r
}
