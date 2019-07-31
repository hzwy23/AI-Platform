package dao

import (
	"ai-platform/api/entity"
	"ai-platform/dbobj"
	"ai-platform/panda"
	"errors"
)

// 设备分组信息表
type DeviceInstallInfoDao interface {
	// 查询所有分组信息
	FindAll() ([]entity.DeviceInstallInfo, error)

	// 根据ID更新分组
	UpdateById(item entity.DeviceInstallInfo) error

	// 逻辑删除分组
	LogicDeleteById(groupId int) (int64, error)

	// 新增分组
	Insert(item entity.DeviceInstallInfo) error
}

type deviceInstallDaoImpl struct {
}

func (r *deviceInstallDaoImpl) FindAll() ([]entity.DeviceInstallInfo, error) {
	rst := make([]entity.DeviceInstallInfo, 0)
	err := dbobj.QueryForSlice("select id,serial_number, device_address, lat, lon, create_date, create_by, update_date, update_by from device_install_info where delete_status = 0", &rst)
	return rst, err
}

func (r *deviceInstallDaoImpl) UpdateById(item entity.DeviceInstallInfo) error {
	rst, err := dbobj.Exec("update device_install_info set serial_number = ?, device_address = ?, lat = ?, lon = ?,  update_date = ?, update_by = ?  where  id = ? and  delete_status = 0",
		item.SerialNumber, item.DeviceAddress, item.Lat, item.Lon, panda.CurTime(), item.UpdateBy, item.Id)
	if err != nil {
		return err
	}
	size, err := rst.RowsAffected()
	if err != nil {
		return err
	}
	if size > 0 {
		return nil
	}
	return errors.New("安装的设备不存在")
}

func (r *deviceInstallDaoImpl) LogicDeleteById(id int) (int64, error) {
	rst, err := dbobj.Exec("update device_install_info set delete_status = 1 where id = ? ", id)
	if err != nil {
		return 0, err
	}
	size, err := rst.RowsAffected()
	if err != nil {
		return 0, err
	}
	if size > 0 {
		return size, nil
	}
	return 0, errors.New("安装设备不存在")
}

func (r *deviceInstallDaoImpl) Insert(item entity.DeviceInstallInfo) error {
	rst, err := dbobj.Exec("insert into device_install_info(id, serial_number, device_address, lat, lon, create_date, create_by, update_date, update_by, delete_status) values(?,?,?,?,?,?,?,?,?,?)",
		item.Id, item.SerialNumber, item.DeviceAddress, item.Lat, item.Lon, item.CreateDate, item.CreateBy, item.UpdateDate, item.UpdateBy, 0)
	if err != nil {
		return err
	}
	size, err := rst.RowsAffected()
	if err != nil {
		return err
	}
	if size > 0 {
		return nil
	}
	return errors.New("新增安装设备信息失败")
}

func NewDeviceInstallDao() DeviceInstallInfoDao {
	r := &deviceInstallDaoImpl{}
	return r
}
