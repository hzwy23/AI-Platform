package handle

import (
	"ai-platform/panda/crypto/aes"
	"ai-platform/panda/logger"
	"database/sql"
	"strconv"
)

func Connect(typ string)  (*sql.DB, error){

	red, err := GetConfig()
	if err != nil {
		panic("cant not read ./conf/dbobj.conf.please check this file.")
	}

	tns, _ := red.Get("DB.tns")
	usr, _ := red.Get("DB.user")
	pad, _ := red.Get("DB.passwd")
	mc, _ := red.Get("DB.maxConn")


	if len(pad) == 24 {
		pad, err = aes.Decrypt(pad)
		if err != nil {
			logger.Error("Decrypt mysql passwd failed.")
			return nil, err
		}
	}

	db, err := sql.Open(typ, usr+":"+pad+"@"+tns)

	maxConn := 100
	if len(mc) != 0 {
		mx, err := strconv.Atoi(mc)
		if err == nil {
			maxConn = mx
		}
	}
	// 设置连接池最大值
	db.SetMaxOpenConns(maxConn)
	db.SetConnMaxLifetime(0)

	logger.Info("create mysql handle success. max connect value is:", maxConn)
	return db, savePassword(pad)
}

func savePassword(pad string) error {

	if len(pad) != 24 {
		psd, err := aes.Encrypt(pad)
		if err != nil {
			logger.Error("decrypt passwd failed." + psd)
			return nil
		}
		psd = "\"" + psd + "\""

		red, err := GetConfig()
		if err != nil {
			panic("cant not read ./conf/dbobj.conf.please check this file.")
		}
		red.Set("DB.passwd", psd)
	}
	return nil
}