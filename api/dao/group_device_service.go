package dao

import (
	"ai-platform/api/entity"
	"ai-platform/panda"
)

type GroupDeviceService interface {
	AddGroup(groupName string, userId string) error
	FindByPage(pageNumber int, pageSize int) ([]entity.DeviceGroupInfo, error)
	DeleteByGroupId(groupId int) (int64, error)
	UpdateGroupName(item entity.DeviceGroupInfo) error
}

type groupDeviceServiceImpl struct {
	dao DeviceGroupInfoDao
}

func (r *groupDeviceServiceImpl) UpdateGroupName(item entity.DeviceGroupInfo) error {
	return r.dao.UpdateById(item)
}

func (r *groupDeviceServiceImpl) FindByPage(pageNumber int, pageSize int) ([]entity.DeviceGroupInfo, error) {
	return r.dao.FindAll()
}

func (r *groupDeviceServiceImpl) DeleteByGroupId(groupId int) (int64, error) {
	return r.dao.LogicDeleteById(groupId)
}

func (r *groupDeviceServiceImpl) AddGroup(groupName string, userId string) error {
	item := entity.DeviceGroupInfo{
		GroupName:    groupName,
		CreateBy:     userId,
		CreateDate:   panda.CurTime(),
		UpdateBy:     userId,
		UpdateDate:   panda.CurTime(),
		DeleteStatus: 0,
	}
	return r.dao.Insert(item)
}

func NewGroupDeviceService() GroupDeviceService {
	r := &groupDeviceServiceImpl{
		dao: NewDeviceGroupInfoDao(),
	}
	return r
}
