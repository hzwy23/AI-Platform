package controller

import (
	"ai-platform/api/service"
	"ai-platform/dbobj"
	"ai-platform/panda/hret"
	"ai-platform/panda/logger"
	"ai-platform/panda/route"
	"net/http"
)

type HomePageController struct {
	srv service.DeviceScanService
}

func (r *HomePageController) Get(resp http.ResponseWriter, req *http.Request) {
	item := int64(0)
	err := dbobj.QueryForObject("select sum(power_total) from device_manage_info where delete_status = 0", dbobj.PackArgs(), &item)
	if err != nil {
		logger.Error(err)
	}
	tCnt, offlineCnt, lampCnt := int64(0), int64(0), int64(0)
	err = dbobj.QueryForObject("select count(*) from event_alarm_info where delete_status = 0 and handle_status = 0 and event_type_cd = ?", dbobj.PackArgs(1), &tCnt)
	err = dbobj.QueryForObject("select count(*) from event_alarm_info where delete_status = 0 and handle_status = 0 and event_type_cd = ?", dbobj.PackArgs(2), &offlineCnt)
	err = dbobj.QueryForObject("select count(*) from event_alarm_info where delete_status = 0 and handle_status = 0 and event_type_cd = ?", dbobj.PackArgs(3), &lampCnt)
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
		srv: service.NewDeviceScanService(),
	}
	route.Handler("GET", "/api/homepage/statistics", ctl.Get)
}
