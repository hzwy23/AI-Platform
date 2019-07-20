package dao

import (
	"ai-platform/api/entity"
	"ai-platform/dbobj"
)

// 设备分组信息表
type DeviceGroupInfoDao interface {
	// 查询所有分组信息
	FindAll() ([]entity.DeviceGroupInfo, error)

	// 根据ID查找分组详细信息
	FindById(groupId int) (entity.DeviceGroupInfo, error)

	// 根据ID更新分组
	UpdateById(item entity.DeviceGroupInfo) error

	// 逻辑删除分组
	LogicDeleteById(groupId int) (int64, error)

	// 新增分组
	Insert(item entity.DeviceGroupInfo) error
}

func NewDeviceGroupInfoDao() DeviceGroupInfoDao {
	r := &deviceGroupInfoDaoImpl{}
	return r
}

type deviceGroupInfoDaoImpl struct {
}

// 查询所有分组信息
func (r *deviceGroupInfoDaoImpl) FindAll() ([]entity.DeviceGroupInfo, error) {
	rst := make([]entity.DeviceGroupInfo, 0)
	err := dbobj.QueryForSlice("select group_id, group_name, create_by, create_date, update_by, update_date from device_group_info where delete_status = 0", &rst)
	return rst, err
}

// 根据ID查找分组详细信息
func (r *deviceGroupInfoDaoImpl) FindById(groupId int) (entity.DeviceGroupInfo, error) {
	rst := entity.DeviceGroupInfo{}
	err := dbobj.QueryForStruct("select group_id, group_name, create_by, create_date, update_by, update_date from device_group_info where delete_status = 0 and group_id = ?", &rst, groupId)
	return rst, err
}

// 根据ID更新分组
func (r *deviceGroupInfoDaoImpl) UpdateById(item entity.DeviceGroupInfo) error {
	return nil
}

// 逻辑删除分组
func (r *deviceGroupInfoDaoImpl) LogicDeleteById(groupId int) (int64, error) {
	ret, err := dbobj.Exec("update device_group_info set delete_status = 1 where group_id = ?", groupId)
	if ret == nil {
		return 0, err
	}
	size, _ := ret.RowsAffected()
	return size, err
}

// 新增分组
func (r *deviceGroupInfoDaoImpl) Insert(item entity.DeviceGroupInfo) error {
	_, err := dbobj.Exec("insert into device_group_info(group_id, group_name, create_by, create_date, update_by, update_date, delete_status) values(?,?,?,?,?,?,0)",
		item.GroupId, item.GroupName, item.CreateBy, item.CreateDate, item.UpdateBy, item.UpdateDate)
	return err
}
