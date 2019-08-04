package controller

import (
	"ai-platform/api/dao"
	"ai-platform/panda/hret"
	"ai-platform/panda/route"
	"net/http"
)

type EventController struct {
	dao dao.EventAlarmInfoDao
}

func (r *EventController) Get(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	eventTypeId := req.FormValue("EventTypeId")
	if len(eventTypeId) == 0 {
		hret.Error(resp, 500300, "无效的事件类型")
		return
	}
	item, err := r.dao.FindByTypeCd(eventTypeId)
	if err != nil {
		hret.Error(resp, 50030, err.Error())
		return
	}
	hret.Success(resp, item)
}

func (r *EventController) Put(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	Id := req.FormValue("Id")
	if len(Id) == 0 {
		hret.Error(resp, 500300, "无效的事件")
		return
	}
	_, err := r.dao.CloseById(1, Id)
	if err != nil {
		hret.Error(resp, 500030, err.Error())
		return
	}
	hret.Success(resp, "Success")
}

func init() {
	ctl := EventController{
		dao: dao.NewEventAlarmInfoDao(),
	}
	route.Handler("GET", "/api/device/offline", ctl.Get)
	route.Handler("GET", "/api/device/temperature", ctl.Get)
	route.Handler("GET", "/api/device/lamp/exception", ctl.Get)
	route.Handler("PUT", "/api/device/event", ctl.Put)
}
