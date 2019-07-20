package dao

import (
	"ai-platform/api/entity"
	"ai-platform/dbobj"
)

type GroupDeviceBindDao interface {
	FindAll() ([]entity.GroupDeviceBind, error)
	FindByGroupId(groupId int) ([]entity.GroupDeviceBind, error)
	FindByDeviceId(deviceId int) ([]entity.GroupDeviceBind, error)
	FindById(id int) (entity.GroupDeviceBind, error)
	LogicDeleteById(id int) (int64, error)
	Insert(item entity.GroupDeviceBind) (int64, error)
	Update(groupId int, id int) (int64, error)
	GlobalSetting() error
}

type groupDeviceBindDaoImpl struct {
}

func (r *groupDeviceBindDaoImpl) FindAll() ([]entity.GroupDeviceBind, error) {
	rst := make([]entity.GroupDeviceBind, 0)
	err := dbobj.QueryForSlice("select id, group_id, device_id, create_by, create_date, update_by, update_date, delete_status from group_device_bind where delete_status = 0", &rst)
	return rst, err
}

func (r *groupDeviceBindDaoImpl) FindByGroupId(groupId int) ([]entity.GroupDeviceBind, error) {
	rst := make([]entity.GroupDeviceBind, 0)
	err := dbobj.QueryForSlice("select id, group_id, device_id, create_by, create_date, update_by, update_date, delete_status from group_device_bind where delete_status = 0 and group_id = ?", &rst, groupId)
	return rst, err
}

func (r *groupDeviceBindDaoImpl) FindByDeviceId(deviceId int) ([]entity.GroupDeviceBind, error) {
	rst := make([]entity.GroupDeviceBind, 0)
	err := dbobj.QueryForSlice("select id, group_id, device_id, create_by, create_date, update_by, update_date, delete_status from group_device_bind where delete_status = 0 and device_id = ?", &rst, deviceId)
	return rst, err
}

func (r *groupDeviceBindDaoImpl) FindById(id int) (entity.GroupDeviceBind, error) {
	rst := entity.GroupDeviceBind{}
	err := dbobj.QueryForStruct("select id, group_id, device_id, create_by, create_date, update_by, update_date, delete_status from group_device_bind where delete_status = 0 and id = ?", &rst, id)
	return rst, err
}

func (r *groupDeviceBindDaoImpl) LogicDeleteById(id int) (int64, error) {
	rst, err := dbobj.Exec("update group_device_bind set delete_status = 1 where id = ? and delete_status = 0", id)
	if rst == nil {
		return 0, err
	}
	size, _ := rst.RowsAffected()
	return size, err
}

func (r *groupDeviceBindDaoImpl) Insert(item entity.GroupDeviceBind) (int64, error) {
	rst, err := dbobj.Exec("insert into group_device_bind(id, group_id, device_id, create_by, create_date, update_by, update_date, delete_status) values(?,?,?,?,?,?,?,0)",
		item.Id, item.GroupId, item.DeviceId, item.CreateBy, item.CreateDate, item.UpdateBy, item.UpdateDate)
	if rst == nil {
		return 0, err
	}
	size, _ := rst.RowsAffected()
	return size, err
}

func (r *groupDeviceBindDaoImpl) Update(groupId int, id int) (int64, error) {
	rst, err := dbobj.Exec("update group_device_bind set group_id = ? where id = ?", groupId, id)
	if rst == nil {
		return 0, err
	}
	size, _ := rst.RowsAffected()
	return size, err
}

func (r *groupDeviceBindDaoImpl) GlobalSetting() error {
	panic("implement me")
}

func NewGroupDeviceBindDao() GroupDeviceBindDao {
	r := &groupDeviceBindDaoImpl{}
	return r
}
