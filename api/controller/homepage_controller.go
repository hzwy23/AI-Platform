package controller

import (
	"ai-platform/dbobj"
	"ai-platform/panda/hret"
	"ai-platform/panda/logger"
	"ai-platform/panda/route"
	service2 "ai-platform/server/service"
	"net/http"
)

type HomePageController struct {
	srv service2.DeviceScanService
}

func (r *HomePageController) Get(resp http.ResponseWriter, req *http.Request) {
	item := int64(0)
	err := dbobj.QueryForObject("select sum(power_total) from device_manage_info where delete_status = 0", dbobj.PackArgs(), &item)
	if err != nil {
		logger.Error(err)
	}

	tCnt := dbobj.Count("select count(*) from event_alarm_info where delete_status = 0 and handle_status = 0 and event_type_cd = ?", 1)
	offlineCnt := dbobj.Count("select count(*) from event_alarm_info where delete_status = 0 and handle_status = 0 and event_type_cd = ?", 2)
	lampCnt := dbobj.Count("select count(*) from event_alarm_info where delete_status = 0 and handle_status = 0 and event_type_cd = ?", 3)

	rst, sCnt, err := r.srv.FindAll()
	result := make(map[string]int64)
	result["TotalPower"] = item
	result["TemperatureWarnCnt"] = tCnt
	result["OfflineWarnCnt"] = offlineCnt
	result["LampWarnCnt"] = lampCnt
	result["ScanDeviceCnt"] = int64(len(rst))
	result["AddedDevice"] = int64(sCnt)
	hret.Success(resp, result)

}

func init() {
	ctl := HomePageController{
		srv: service2.NewDeviceScanService(),
	}
	route.Handler("GET", "/api/homepage/statistics", ctl.Get)
}
