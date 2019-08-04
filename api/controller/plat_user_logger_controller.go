package controller

import (
	"ai-platform/api/dao"
	"ai-platform/panda/hret"
	"ai-platform/panda/route"
	"net/http"
	"strconv"
)

type PlatUserLoggerController struct {
	dao dao.PlatUserLoggerDao
}

func (r *PlatUserLoggerController) Get(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	pageNumber := req.FormValue("pageNumber")
	pageSize := req.FormValue("pageSize")
	page, err := strconv.Atoi(pageNumber)
	if err != nil {
		page = 1
	}
	size, err := strconv.Atoi(pageSize)
	if err != nil {
		size = 10
	}
	rst, total, err := r.dao.FindAll(page, size)
	if err != nil {
		hret.Error(resp, 500030, err.Error())
		return
	}

	data := make(map[string]interface{})
	data["total"] = total
	data["content"] = rst
	hret.Success(resp, data)
}

func init() {
	ctl := &PlatUserLoggerController{
		dao: dao.NewPlatUserLoggerDao(),
	}
	route.Handler("GET", "/api/platform/logger", ctl.Get)
}
