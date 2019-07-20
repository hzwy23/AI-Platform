package dao

import (
	"ai-platform/api/entity"
	"ai-platform/dbobj"
)

type DeviceManageInfoDao interface {
	FindAll() ([]entity.DeviceManageInfo, error)
	FindByDeviceId(deviceId int) (entity.DeviceManageInfo, error)
	LogicDeleteById(deviceId int) (int64, error)
	Update(item entity.DeviceManageInfo) (int64, error)
	Insert(item entity.DeviceManageInfo) (int64, error)
}

type deviceManageInfoDaoImpl struct {
}

func NewDeviceManageInfoDao() DeviceManageInfoDao {
	r := &deviceManageInfoDaoImpl{}
	return r
}

func (r *deviceManageInfoDaoImpl) FindAll() ([]entity.DeviceManageInfo, error) {
	rst := make([]entity.DeviceManageInfo, 0)
	err := dbobj.QueryForSlice("select device_id, serial_number, device_name, dhcp_flag, device_ip, device_port, device_status, device_attribute, device_power, device_light_threshold, device_brightness, device_temperature, auto_start_time, auto_end_time, light_mode, mac_address, firmware_version, longitude, latitude, address, mask, gateway, pin, create_by, create_date, update_by, update_data, delete_status, power_total, strobe_count from device_manage_info where delete_status = 0", &rst)
	return rst, err
}

func (r *deviceManageInfoDaoImpl) FindByDeviceId(deviceId int) (entity.DeviceManageInfo, error) {
	rst := entity.DeviceManageInfo{}
	err := dbobj.QueryForStruct("select device_id, serial_number, device_name, dhcp_flag, device_ip, device_port, device_status, device_attribute, device_power, device_light_threshold, device_brightness, device_temperature, auto_start_time, auto_end_time, light_mode, mac_address, firmware_version, longitude, latitude, address, mask, gateway, pin, create_by, create_date, update_by, update_data, delete_status, power_total, strobe_count from device_manage_info where delete_status = 0 and device_id = ?", &rst, deviceId)
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
	result, err := dbobj.Exec("update device_manage_info set device_name = ?, dhcp_flag = ?, device_ip = ?, device_port = ?, device_status = ?, device_attribute = ?, device_power = ?, device_light_threshold = ?, device_brightness = ?, device_temperature = ?, auto_start_time = ?, auto_end_time = ?, light_mode = ?, mac_address = ?, firmware_version = ?, longitude = ?, latitude = ?, address = ?, mask = ?, gateway = ?, pin = ?, update_by = ?, update_data = ?, power_total = ?, strobe_count = ? where delete_status = 0 and device_id = ?",
		item.DeviceName, item.DhcpFlag, item.DeviceIp,
		item.DevicePort, item.DeviceStatus, item.DeviceAttribute,
		item.DevicePower, item.DeviceLightThreshold, item.DeviceBrightness,
		item.DeviceTemperature, item.AutoStartTime, item.AutoEndTime,
		item.LightMode, item.MacAddress, item.FirmwareVersion,
		item.Longitude, item.Latitude, item.Address,
		item.Mask, item.Gateway, item.Pin, item.UpdateBy,
		item.UpdateData, item.PowerTotal, item.StrobeCount, item.DeviceId)

	if result != nil {
		return 0, err
	}
	size, err := result.RowsAffected()
	return size, nil
}

func (r *deviceManageInfoDaoImpl) Insert(item entity.DeviceManageInfo) (int64, error) {
	result, err := dbobj.Exec("insert into device_manage_info(device_id, serial_number, device_name, dhcp_flag, device_ip, device_port, device_status, device_attribute, device_power, device_light_threshold, device_brightness, device_temperature, auto_start_time, auto_end_time, light_mode, mac_address, firmware_version, longitude, latitude, address, mask, gateway, pin, create_by, create_date, update_by, update_data, delete_status, power_total, strobe_count) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		item.DeviceId, item.SerialNumber, item.DeviceName,
		item.DhcpFlag, item.DeviceIp, item.DevicePort,
		item.DeviceStatus, item.DeviceAttribute,
		item.DevicePower, item.DeviceLightThreshold, item.DeviceBrightness,
		item.DeviceTemperature, item.AutoStartTime, item.AutoEndTime,
		item.LightMode, item.MacAddress, item.FirmwareVersion,
		item.Longitude, item.Latitude, item.Address,
		item.Mask, item.Gateway, item.Pin,
		item.CreateBy, item.CreateDate, item.UpdateBy, item.UpdateData, 0, item.PowerTotal, item.StrobeCount)

	if result == nil {
		return 0, err
	}
	size, err := result.RowsAffected()
	return size, err
}
