package dao_test

import (
	"ai-platform/api/dao"
	"ai-platform/api/entity"
	"fmt"
	"testing"
)

func TestNewPlatUserLoggerDao(t *testing.T) {
	r := dao.NewPlatUserLoggerDao()
	fmt.Println(r.FindAll(1, 10))
}

func TestNewPlatUserLoggerDao2(t *testing.T) {
	r := dao.NewPlatUserLoggerDao()
	item := entity.PlatUserLogger{
		Id:         3,
		UserId:     "1",
		HandleTime: "2019-01-01",
		ReqMethod:  "POST",
		ReqUrl:     "/music",
		ReqParam:   "param",
		RetMsg:     "hello",
		RetCode:    "200",
	}

	fmt.Println(r.Insert(item))
}
