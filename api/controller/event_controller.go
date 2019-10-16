package controller

import (
	"ai-platform/dbobj"
	"ai-platform/api/dao"
	"ai-platform/api/entity"
	"ai-platform/panda/hret"
	"ai-platform/panda/route"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
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

	eventTypeCd := req.FormValue("eventTypeCd")

	if eventTypeCd == "3" {
		// 取消灯珠异常
		td, _ := strconv.Atoi(eventTypeCd)
		item, err := r.dao.FindById(td)

		if err != nil {
			hret.Error(resp, 500301, "记录不存在")
			return
		}
		dbobj.Exec("update device_manage_info set device_status = concat(substr(device_status,1,1), 1, substr(device_status,3,1)) where serial_number = ? and delete_status = 0", item.SerialNumber)
	}


	_, err := r.dao.CloseById(1, Id)
	if err != nil {
		hret.Error(resp, 500030, err.Error())
		return
	}
	hret.Success(resp, "Success")
}

func (r *EventController)Delete(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm();
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		hret.Error(resp, 500030, "参数解析失败，请联系管理员")
		return
	}
	var rst []entity.EventAlarmInfo
	err = json.Unmarshal(body,&rst)
	if err != nil {
		hret.Error(resp, 50060, "参数格式不正确，请联系管理员")
		return
	}
	for _, item := range rst {
		r.dao.LogicDeleteById(item.Id)
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
	route.Handler("POST", "/api/device/lamp/remove/event", ctl.Delete)
}
