package dao

type SysUserInfoDao interface {
}

func NewSysUserInfoDaoImpl() SysUserInfoDao {
	r := &SysUserInfoDaoImpl{}
	return r
}

type SysUserInfoDaoImpl struct{}
