package service

import "ai-platform/api/entity"

type DeviceManageService interface {
	findBySerialNumber(serialNumber string) (entity.DeviceManageInfo, error)
	AddDevice() error
	UpdateDevice(item entity.DeviceManageInfo) error
	RemoveDevice(deviceId string) error
}

type deviceManageServiceImpl struct {
}

func (r *deviceManageServiceImpl) findBySerialNumber(serialNumber string) (entity.DeviceManageInfo, error) {
	panic("implement me")
}

func (r *deviceManageServiceImpl) AddDevice() error {
	panic("implement me")
}

func (r *deviceManageServiceImpl) UpdateDevice(item entity.DeviceManageInfo) error {
	panic("implement me")
}

func (r *deviceManageServiceImpl) RemoveDevice(deviceId string) error {
	panic("implement me")
}

func NewDeviceManageService() DeviceManageService {
	r := &deviceManageServiceImpl{}
	return r
}
