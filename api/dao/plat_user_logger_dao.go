package dao

import (
	"ai-platform/api/entity"
	"ai-platform/dbobj"
)

type PlatUserLoggerDao interface {
	FindAll(pageNumber int, pageSize int) ([]entity.PlatUserLogger, error)
	Insert(item entity.PlatUserLogger) (int64, error)
}

type platUserLoggerDaoImpl struct {
}

func (r *platUserLoggerDaoImpl) FindAll(pageNumber int, pageSize int) ([]entity.PlatUserLogger, error) {
	rst := make([]entity.PlatUserLogger, 0)
	if pageNumber > 0 {
		pageNumber -= 1
	} else {
		pageNumber = 0
		pageSize = 10
	}
	start := pageNumber * pageSize
	end := (pageNumber + 1) * pageSize

	err := dbobj.QueryForSlice("select id, user_id, handle_time, req_method, req_url, req_param, ret_msg, ret_code from plat_user_logger limit ?,?", &rst, start, end)
	return rst, err
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
