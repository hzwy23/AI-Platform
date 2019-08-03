package dao

import (
	"ai-platform/api/entity"
	"ai-platform/dbobj"
	"fmt"
)

type PlatDeviceLoggerDao interface {
	FindAll(pageNumber int, pageSize int) ([]entity.PlatDeviceLogger, int64, error)
	Insert(item entity.PlatDeviceLogger) (int64, error)
}

type platDeviceLoggerDaoImpl struct {
}

func (r *platDeviceLoggerDaoImpl) Insert(item entity.PlatDeviceLogger) (int64, error) {
	rst, err := dbobj.Exec("insert into plat_device_logger(id, device_id, handle_time, direction, biz_type, message, ret_code, ret_msg) values(?,?,?,?,?,?,?,?)",
		item.Id, item.SerialNumber, item.HandleTime, item.Direction, item.BizType, item.Message, item.RetCode, item.RetMsg)
	if rst == nil {
		return 0, err
	}
	size, _ := rst.RowsAffected()
	return size, err
}

func (r *platDeviceLoggerDaoImpl) FindAll(pageNumber int, pageSize int) ([]entity.PlatDeviceLogger, int64,  error) {
	rst := make([]entity.PlatDeviceLogger, 0)

	start := int64((pageNumber - 1) * pageSize)
	if start < 0 {
		start = 0
	}

	total := dbobj.Count("select count(*) from plat_device_logger")
	if start > total {
		start = 0
	}

	fmt.Println(start, pageSize)
	err := dbobj.QueryForSlice("select id, direction, biz_type, message, ret_code, ret_msg, serial_number, handle_time from plat_device_logger order by id desc limit ?,?", &rst, start, pageSize)

	return rst, total, err
}

func NewPlatDeviceLoggerDao() PlatDeviceLoggerDao {
	r := &platDeviceLoggerDaoImpl{}
	return r
}
