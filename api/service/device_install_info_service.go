package service

import (
	"ai-platform/api/dao"
	"ai-platform/api/entity"
)

type DeviceInstallInfoService interface {
	// 查询所有分组信息
	FindAll() ([]entity.DeviceInstallInfo, error)

	// 根据ID更新分组
	UpdateById(item entity.DeviceInstallInfo) error

	// 逻辑删除分组
	LogicDeleteById(id int) (int64, error)

	// 新增分组
	Insert(item entity.DeviceInstallInfo) error
}

type DeviceInstallInfoServiceImpl struct {
	dao dao.DeviceInstallInfoDao
}

func (r *DeviceInstallInfoServiceImpl) FindAll() ([]entity.DeviceInstallInfo, error) {
	return r.dao.FindAll()
}

func (r *DeviceInstallInfoServiceImpl) UpdateById(item entity.DeviceInstallInfo) error {
	return r.dao.UpdateById(item)
}

func (r *DeviceInstallInfoServiceImpl) LogicDeleteById(id int) (int64, error) {
	return r.dao.LogicDeleteById(id)
}

func (r *DeviceInstallInfoServiceImpl) Insert(item entity.DeviceInstallInfo) error {
	return r.dao.Insert(item)
}

func NewDeviceInstallService() DeviceInstallInfoService {
	r := &DeviceInstallInfoServiceImpl{
		dao: dao.NewDeviceInstallDao(),
	}
	return r
}
