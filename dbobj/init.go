package dbobj

import (
	"ai-platform/dbobj/dbhandle"

	_ "ai-platform/dbobj/mysql"
)

func init() {
	conf, err := dbhandle.GetConfig()
	if err != nil {
		panic("init database failed." + err.Error())
	}
	Default, err = conf.Get("DB.type")
	if err != nil {
		panic("get default database type failed." + err.Error())
	}
	InitDB(Default)
}
