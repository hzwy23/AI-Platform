package controller

import (
	"ai-platform/api/dao"
	"ai-platform/panda/hret"
	"ai-platform/panda/route"
	"net/http"
	"strconv"
)

type DeviceLoggerController struct {
	dao dao.PlatDeviceLoggerDao
}

func (r *DeviceLoggerController)Get(resp http.ResponseWriter, req *http.Request)  {
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
	rst, total, err:=r.dao.FindAll(page, size)
	if err != nil {
		hret.Error(resp, 500030, err.Error())
		return
	}

	data := make(map[string]interface{})
	data["total"] = total
	data["content"] = rst

	hret.Success(resp, data)
}

func init()  {
	ctl := &DeviceLoggerController{
		dao: dao.NewPlatDeviceLoggerDao(),
	}
	route.Handler("GET", "/api/device/logger", ctl.Get)
}
