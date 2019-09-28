package dao

import (
	"ai-platform/api/entity"
	"ai-platform/dbobj"
)

type EventAlarmInfoDao interface {
	FindAll() ([]entity.EventAlarmInfo, error)
	FindByTypeCd(typeCd string) ([]entity.EventAlarmInfo, error)
	FindById(id int) (entity.EventAlarmInfo, error)
	LogicDeleteById(id int) (int64, error)
	CloseById(status int, id string) (int64, error)
	Insert(item entity.EventAlarmInfo) (int64, error)
	Update(item entity.EventAlarmInfo) (int64, error)
	ChangeHandleStatus(serialNumber string, eventTypeId int)
}

type eventAlarmInfoDaoImpl struct {
}

func (r *eventAlarmInfoDaoImpl) ChangeHandleStatus(serialNumber string, eventTypeId int) {
	dbobj.Exec("update event_alarm_info set handle_status = 1 where serial_number = ? and event_type_cd = ? and delete_status = 0", serialNumber, eventTypeId)
}

func (r *eventAlarmInfoDaoImpl) Insert(item entity.EventAlarmInfo) (int64, error) {
	ret, err := dbobj.Exec("insert into event_alarm_info(id, event_type_cd, occurrence_time, serial_number, device_name, device_ip, device_attribute, device_brightness, device_temperature, handle_status, delete_status) values(?,?,?,?,?,?,?,?,?,?,0)",
		item.Id, item.EventTypeCd, item.OccurrenceTime,
		item.SerialNumber, item.DeviceName, item.DeviceIp,
		item.DeviceAttribute, item.DeviceBrightness, item.DeviceTemperature,
		item.HandleStatus)
	if ret == nil {
		return 0, err
	}
	size, _ := ret.RowsAffected()
	return size, err
}

func (r *eventAlarmInfoDaoImpl) Update(item entity.EventAlarmInfo) (int64, error) {
	ret, err := dbobj.Exec("update event_alarm_info set device_attribute = ?, device_brightness = ?, device_temperature = ? where id = ? ",
		item.DeviceAttribute, item.DeviceBrightness, item.DeviceTemperature, item.Id)
	if ret == nil {
		return 0, err
	}
	size, _ := ret.RowsAffected()
	return size, err
}

func (r *eventAlarmInfoDaoImpl) CloseById(status int, id string) (int64, error) {
	result, err := dbobj.Exec("update event_alarm_info set handle_status = ? where id = ?", status, id)
	if result == nil {
		return 0, err
	}
	size, _ := result.RowsAffected()
	return size, err
}

func (r *eventAlarmInfoDaoImpl) FindAll() ([]entity.EventAlarmInfo, error) {
	rst := make([]entity.EventAlarmInfo, 0)
	err := dbobj.QueryForSlice("select id, event_type_cd, occurrence_time, serial_number, device_name, device_ip, device_attribute, device_brightness, device_temperature, handle_status from event_alarm_info where delete_status = 0", &rst)
	return rst, err
}

func (r *eventAlarmInfoDaoImpl) FindByTypeCd(typeCd string) ([]entity.EventAlarmInfo, error) {
	rst := make([]entity.EventAlarmInfo, 0)
	err := dbobj.QueryForSlice("select id, event_type_cd, occurrence_time, serial_number, device_name, device_ip, device_attribute, device_brightness, device_temperature, handle_status from event_alarm_info where delete_status = 0 and event_type_cd = ? order by handle_status asc, id desc limit 0, 100", &rst, typeCd)
	return rst, err
}

func (r *eventAlarmInfoDaoImpl) FindById(id int) (entity.EventAlarmInfo, error) {
	rst := entity.EventAlarmInfo{}
	err := dbobj.QueryForStruct("select id, event_type_cd, occurrence_time, serial_number, device_name, device_ip, device_attribute, device_brightness, device_temperature, handle_status from event_alarm_info where delete_status = 0 and id = ?", &rst, id)
	return rst, err
}

func (r *eventAlarmInfoDaoImpl) LogicDeleteById(id int) (int64, error) {
	result, err := dbobj.Exec("update event_alarm_info set delete_status = 1 where id = ?", id)
	if result == nil {
		return 0, err
	}
	size, _ := result.RowsAffected()
	return size, err
}

func NewEventAlarmInfoDao() EventAlarmInfoDao {
	return &eventAlarmInfoDaoImpl{}
}
