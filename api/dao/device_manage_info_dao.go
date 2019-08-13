package dao

import (
	"ai-platform/api/entity"
	"ai-platform/dbobj"
	"ai-platform/panda/logger"
	"errors"
	"strings"
)

type DeviceManageInfoDao interface {
	FindAll(groupId string) ([]entity.DeviceManageInfo, error)
	FindByDeviceId(deviceId int) (entity.DeviceManageInfo, error)
	FindBySerialNumber(serialNumber string) (entity.DeviceManageInfo, error)
	LogicDeleteById(deviceId int) (int64, error)
	Update(item entity.DeviceManageInfo) (int64, error)
	Insert(item entity.DeviceManageInfo, GroupId string) (int64, error)
}

type deviceManageInfoDaoImpl struct {
}

func (r *deviceManageInfoDaoImpl) FindBySerialNumber(serialNumber string) (entity.DeviceManageInfo, error) {
	rst := entity.DeviceManageInfo{}
	err := dbobj.QueryForStruct("select device_id, serial_number, device_name, dhcp_flag, device_ip, device_port, device_status, device_attribute, device_power, device_light_threshold, device_brightness, device_temperature, auto_start_time, auto_end_time, light_mode, mac_address, firmware_version, '', '', mask, gateway, pin, create_by, create_date, update_by, update_data, delete_status, power_total, strobe_count, '0', '', flash_duration from device_manage_info where delete_status = 0 and serial_number = ?", &rst, serialNumber)
	return rst, err
}

func NewDeviceManageInfoDao() DeviceManageInfoDao {
	r := &deviceManageInfoDaoImpl{}
	return r
}

func (r *deviceManageInfoDaoImpl) FindAll(groupId string) ([]entity.DeviceManageInfo, error) {
	rst := make([]entity.DeviceManageInfo, 0)
	err := errors.New("未知错误")
	if len(strings.TrimSpace(groupId)) == 0 {
		err = dbobj.QueryForSlice("select t.device_id, t.serial_number, t.device_name, t.dhcp_flag, t.device_ip, t.device_port, t.device_status, t.device_attribute, t.device_power, t.device_light_threshold, t.device_brightness, t.device_temperature, t.auto_start_time, t.auto_end_time, t.light_mode, t.mac_address, t.firmware_version, di.lon, di.lat, t.mask, t.gateway, t.pin, t.create_by, t.create_date, t.update_by, t.update_data, t.delete_status, t.power_total, t.strobe_count, p.group_id, d.group_name, t.flash_duration from device_manage_info t left join group_device_bind p on t.device_id = p.device_id and p.delete_status = 0 left join device_group_info d on p.group_id = d.group_id and d.delete_status = 0 left join device_install_info di on t.serial_number = di.serial_number and di.delete_status = 0 where t.delete_status = 0", &rst)
	} else {
		err = dbobj.QueryForSlice("select t.device_id, t.serial_number, t.device_name, t.dhcp_flag, t.device_ip, t.device_port, t.device_status, t.device_attribute, t.device_power, t.device_light_threshold, t.device_brightness, t.device_temperature, t.auto_start_time, t.auto_end_time, t.light_mode, t.mac_address, t.firmware_version, di.lon, di.lat, t.mask, t.gateway, t.pin, t.create_by, t.create_date, t.update_by, t.update_data, t.delete_status, t.power_total, t.strobe_count, p.group_id, d.group_name, t.flash_duration from device_manage_info t left join group_device_bind p on t.device_id = p.device_id and p.delete_status = 0 left join device_group_info d on p.group_id = d.group_id and d.delete_status = 0 left join device_install_info di on t.serial_number = di.serial_number and di.delete_status = 0 where t.delete_status = 0 and p.group_id = ?", &rst, groupId)
	}
	return rst, err
}

func (r *deviceManageInfoDaoImpl) FindByDeviceId(deviceId int) (entity.DeviceManageInfo, error) {
	rst := entity.DeviceManageInfo{}
	err := dbobj.QueryForStruct("select t.device_id, t.serial_number, t.device_name, t.dhcp_flag, t.device_ip, t.device_port, t.device_status, t.device_attribute, t.device_power, t.device_light_threshold, t.device_brightness, t.device_temperature, t.auto_start_time, t.auto_end_time, t.light_mode, t.mac_address, t.firmware_version, di.lon, di.lat, t.mask, t.gateway, t.pin, t.create_by, t.create_date, t.update_by, t.update_data, t.delete_status, t.power_total, t.strobe_count, '0', '', t.flash_duration from device_manage_info t left join device_install_info di on t.serial_number = di.serial_number and di.delete_status = 0 where t.delete_status = 0 and t.device_id = ?", &rst, deviceId)
	return rst, err
}

func (r *deviceManageInfoDaoImpl) LogicDeleteById(deviceId int) (int64, error) {
	ret, err := dbobj.Exec("update device_manage_info set delete_status = 1 where device_id = ?", deviceId)
	if ret == nil {
		return 0, err
	}
	size, _ := ret.RowsAffected()
	return size, err
}

func (r *deviceManageInfoDaoImpl) Update(item entity.DeviceManageInfo) (int64, error) {
	result, err := dbobj.Exec("update device_manage_info set device_name = ?, dhcp_flag = ?, device_ip = ?, device_port = ?, device_status = ?, device_attribute = ?, device_power = ?, device_light_threshold = ?, device_brightness = ?, device_temperature = ?, auto_start_time = ?, auto_end_time = ?, light_mode = ?, mac_address = ?, firmware_version = ?, mask = ?, gateway = ?, pin = ?, update_by = ?, update_data = ?, power_total = ?, strobe_count = ? where delete_status = 0 and device_id = ?",
		item.DeviceName, item.DhcpFlag, item.DeviceIp,
		item.DevicePort, item.DeviceStatus, item.DeviceAttribute,
		item.DevicePower, item.DeviceLightThreshold, item.DeviceBrightness,
		item.DeviceTemperature, item.AutoStartTime, item.AutoEndTime,
		item.LightMode, item.MacAddress, item.FirmwareVersion,
		item.Mask, item.Gateway, item.Pin, item.UpdateBy,
		item.UpdateDate, item.PowerTotal, item.StrobeCount, item.DeviceId)
	logger.Info("更新设备信息，跟新内容是：", item)
	if result != nil {
		return 0, err
	}
	size, err := result.RowsAffected()
	return size, nil
}

func (r *deviceManageInfoDaoImpl) Insert(item entity.DeviceManageInfo, GroupId string) (int64, error) {
	tx, err := dbobj.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec("insert into device_manage_info(device_id, serial_number, device_name, dhcp_flag, device_ip, device_port, device_status, device_attribute, device_power, device_light_threshold, device_brightness, device_temperature, auto_start_time, auto_end_time, light_mode, mac_address, firmware_version, mask, gateway, pin, create_by, create_date, update_by, update_data, delete_status, power_total, strobe_count, flash_duration) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		item.DeviceId, item.SerialNumber, item.DeviceName,
		item.DhcpFlag, item.DeviceIp, item.DevicePort,
		item.DeviceStatus, item.DeviceAttribute,
		item.DevicePower, item.DeviceLightThreshold, item.DeviceBrightness,
		item.DeviceTemperature, item.AutoStartTime, item.AutoEndTime,
		item.LightMode, item.MacAddress, item.FirmwareVersion,
		item.Mask, item.Gateway, item.Pin,
		item.CreateBy, item.CreateDate, item.UpdateBy, item.UpdateDate, 0, item.PowerTotal, item.StrobeCount, item.FlashDuration)

	if result == nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if len(GroupId) > 0 {
		result, err = tx.Exec("insert into group_device_bind(group_id, device_id, create_by, create_date, update_by, update_date, delete_status) values(?,?,?,?,?,?,0)",
			GroupId, id, item.CreateBy, item.CreateDate, item.UpdateBy, item.UpdateDate)
		if result == nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	size, err := result.RowsAffected()
	return size, err
}
