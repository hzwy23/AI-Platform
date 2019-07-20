package dao

import (
	"ai-platform/api/entity"
	"ai-platform/dbobj"
)

type SysGlobalConfigDao interface {
	FindAll() ([]entity.SysGlobalConfig, error)
	Update(itemValue string, itemId string) (int64, error)
}

type sysGlobalConfigDaoImpl struct {
}

func (r *sysGlobalConfigDaoImpl) FindAll() ([]entity.SysGlobalConfig, error) {
	rst := make([]entity.SysGlobalConfig, 0)
	err := dbobj.QueryForSlice("select item_id, item_name, item_value from sys_global_config", &rst)
	return rst, err
}

func (r *sysGlobalConfigDaoImpl) Update(itemValue string, itemId string) (int64, error) {
	rst, err := dbobj.Exec("update sys_global_config set item_value = ? where item_id = ?", itemValue, itemId)
	if rst == nil {
		return 0, err
	}
	size, _ := rst.RowsAffected()
	return size, err
}

func NewSysGlobalConfigDao() SysGlobalConfigDao {
	r := &sysGlobalConfigDaoImpl{}
	return r
}
