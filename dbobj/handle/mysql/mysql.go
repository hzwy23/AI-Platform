package mysql

import (
	"ai-platform/dbobj/handle"
	"ai-platform/panda/logger"
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type mysql struct {
	db *sql.DB
}

func mysqlHandle() handle.DbObj {

	o := new(mysql)

	db, err := handle.Connect("mysql")
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	o.db = db
	return o
}

func (r *mysql) GetErrorCode(err error) string {
	ret := err.Error()
	if n := strings.Index(ret, ":"); n > 0 {
		return strings.TrimSpace(ret[:n])
	} else {
		logger.Error("this error information is not mysql return info")
		return ""
	}
}

func (r *mysql) GetErrorMsg(err error) string {
	ret := err.Error()
	if n := strings.Index(ret, ":"); n > 0 {
		return strings.TrimSpace(ret[n+1:])
	} else {
		logger.Error("this error information is not mysql return info")
		return ""
	}
}

func (r *mysql) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	rows, err := r.db.Query(sql, args...)
	if err != nil {
		if r.db.Ping() != nil {
			logger.Warn("Connection is broken")
			if val, ok := mysqlHandle().(*mysql); ok {
				r.db = val.db
			}
			return r.db.Query(sql, args...)
		}
	}
	return rows, err
}

func (r *mysql) Exec(sql string, args ...interface{}) (sql.Result, error) {
	result, err := r.db.Exec(sql, args...)
	if err != nil {
		if r.db.Ping() != nil {
			logger.Warn("Connection is broken")
			if val, ok := mysqlHandle().(*mysql); ok {
				r.db = val.db
			}
			return r.db.Exec(sql, args...)
		}
	}
	return result, err
}

func (r *mysql) Begin() (*sql.Tx, error) {
	tx, err := r.db.Begin()
	if err != nil {
		if r.db.Ping() != nil {
			logger.Warn("Connection is broken")
			if val, ok := mysqlHandle().(*mysql); ok {
				r.db = val.db
			}
			return r.db.Begin()
		}
	}
	return tx, err
}

func (r *mysql) Prepare(sql string) (*sql.Stmt, error) {
	stmt, err := r.db.Prepare(sql)
	if err != nil {
		if r.db.Ping() != nil {
			logger.Warn("Connection is broken")
			if val, ok := mysqlHandle().(*mysql); ok {
				r.db = val.db
			}
			return r.db.Prepare(sql)
		}
	}
	return stmt, err
}

func (r *mysql) QueryRow(sql string, args ...interface{}) *sql.Row {
	if r.db.Ping() != nil {
		logger.Warn("Connection is broken")
		if val, ok := mysqlHandle().(*mysql); ok {
			r.db = val.db
		}
	}
	return r.db.QueryRow(sql, args...)
}

func init() {
	handle.Register("mysql", mysqlHandle)
}
