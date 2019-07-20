package dao

import (
	"ai-platform/api/entity"
	"ai-platform/dbobj"
)

type PlatDeviceLoggerDao interface {
	FindAll(pageNumber int, pageSize int) ([]entity.PlatDeviceLogger, error)
	Insert(item entity.PlatDeviceLogger) (int64, error)
}

type platDeviceLoggerDaoImpl struct {
}

func (r *platDeviceLoggerDaoImpl) Insert(item entity.PlatDeviceLogger) (int64, error) {
	rst, err := dbobj.Exec("insert into plat_device_logger(id, device_id, handle_time, device_name, direction, biz_type, message, ret_code, ret_msg) values(?,?,?,?,?,?,?,?,?)",
		item.Id, item.DeviceId, item.HandleTime, item.DeviceName, item.Direction, item.BizType, item.Message, item.RetCode, item.RetMsg)
	if rst == nil {
		return 0, err
	}
	size, _ := rst.RowsAffected()
	return size, err
}

func (r *platDeviceLoggerDaoImpl) FindAll(pageNumber int, pageSize int) ([]entity.PlatDeviceLogger, error) {
	rst := make([]entity.PlatDeviceLogger, 0)
	if pageNumber > 0 {
		pageNumber -= 1
	}
	start := pageNumber * pageSize
	end := (pageNumber + 1) * pageSize

	err := dbobj.QueryForSlice("select id, device_id, handle_time, device_name, direction, biz_type, message, ret_code, ret_msg from plat_device_logger limit ?,?", &rst, start, end)
	return rst, err
}

func NewPlatDeviceLoggerDao() PlatDeviceLoggerDao {
	r := &platDeviceLoggerDaoImpl{}
	return r
}
