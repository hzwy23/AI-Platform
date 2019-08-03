package dao

import (
	"ai-platform/api/entity"
	"ai-platform/dbobj"
)

type PlatUserLoggerDao interface {
	FindAll(pageNumber int, pageSize int) ([]entity.PlatUserLogger, int64, error)
	Insert(item entity.PlatUserLogger) (int64, error)
}

type platUserLoggerDaoImpl struct {
}

func (r *platUserLoggerDaoImpl) FindAll(pageNumber int, pageSize int) ([]entity.PlatUserLogger, int64, error) {
	rst := make([]entity.PlatUserLogger, 0)
	start := int64((pageNumber - 1) * pageSize)
	if start < 0 {
		start = 0
	}

	total := dbobj.Count("select count(*) from plat_user_logger")
	if start > total {
		start = 0
	}

	err := dbobj.QueryForSlice("select id, user_id, handle_time, req_method, req_url, req_param, ret_msg, ret_code from plat_user_logger order by id desc limit ?,?", &rst, start, pageSize)
	return rst, total, err
}

func (r *platUserLoggerDaoImpl) Insert(item entity.PlatUserLogger) (int64, error) {
	rst, err := dbobj.Exec("insert into plat_user_logger(id, user_id, handle_time, req_method, req_url, req_param, ret_msg, ret_code) values(?,?,?,?,?,?,?,?)",
		item.Id, item.UserId, item.HandleTime, item.ReqMethod, item.ReqUrl, item.ReqParam, item.RetMsg, item.RetCode)
	if rst == nil {
		return 0, err
	}
	size, _ := rst.RowsAffected()
	return size, err
}

func NewPlatUserLoggerDao() PlatUserLoggerDao {
	r := &platUserLoggerDaoImpl{}
	return r
}
