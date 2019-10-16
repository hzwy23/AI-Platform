package controller

import (
	"ai-platform/api/dao"
	"ai-platform/panda/hret"
	"ai-platform/panda/route"
	"ai-platform/server/listen"
	"net/http"
	"strconv"
)

type SysGlobalConfigController struct {
	dao dao.SysGlobalConfigDao
}

func (r *SysGlobalConfigController) Get(resp http.ResponseWriter, req *http.Request) {
	rst, err := r.dao.FindAll()
	if err != nil {
		hret.Error(resp, 5002332, err.Error())
		return
	}
	hret.Success(resp, rst)
}

func (r *SysGlobalConfigController) Put(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	itemId := req.FormValue("ItemId")
	itemValue := req.FormValue("ItemValue")

		// 修改温度告警阀值
		if itemId == "3" {
			t, err := strconv.ParseFloat(itemValue,32)
			if err != nil {
				hret.Error(resp, 500232, err.Error())
				return
			}
			listen.UpdateWarnTemperature(int(t))
		}


	// 修改心跳包时长
	if itemId == "4" {
		t, err := strconv.ParseFloat(itemValue,32)
		if err != nil {
			hret.Error(resp, 500232, err.Error())
			return
		}
		listen.ChangeHeartbeat(int64(t))
	}


	_, err := r.dao.Update(itemValue, itemId)
	if err != nil {
		hret.Error(resp, 500232, err.Error())
		return
	}

	hret.Success(resp, "Success")
}

func init() {
	ctl := &SysGlobalConfigController{
		dao: dao.NewSysGlobalConfigDao(),
	}
	route.Handler("GET", "/api/global/config", ctl.Get)
	route.Handler("PUT", "/api/global/config", ctl.Put)
}
