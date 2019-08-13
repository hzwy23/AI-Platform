package dbobj

import (
	"ai-platform/dbobj/handle"
	_ "ai-platform/dbobj/handle/mysql"
)

func init() {
	conf, err := handle.GetConfig()
	if err != nil {
		panic("init database failed." + err.Error())
	}
	Default, err = conf.Get("DB.type")
	if err != nil {
		panic("get default database type failed." + err.Error())
	}
	InitDB(Default)
}
