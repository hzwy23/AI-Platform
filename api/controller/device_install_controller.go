package controller

import (
	"ai-platform/api/entity"
	"ai-platform/api/service"
	"ai-platform/panda"
	"ai-platform/panda/hret"
	"ai-platform/panda/jwt"
	"ai-platform/panda/logger"
	"ai-platform/panda/route"
	"net/http"
	"strconv"
)

type DeviceInstallController struct {
	service service.DeviceInstallInfoService
}

func (r *DeviceInstallController) Get(resp http.ResponseWriter, req *http.Request) {
	rst, err := r.service.FindAll()
	if err != nil {
		hret.Error(resp, 500010, err.Error())
		return
	}
	hret.Success(resp, rst)
}

func (r *DeviceInstallController) Post(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	SerialNumber := req.FormValue("SerialNumber")
	DeviceAddress := req.FormValue("DeviceAddress")
	Lat := req.FormValue("Lat")
	Lon := req.FormValue("Lon")
	claim, err := jwt.ParseHttp(req)
	if err != nil {
		logger.Error(err)
		hret.Error(resp, 403, "权限不足")
		return
	}

	arg := entity.DeviceInstallInfo{
		SerialNumber:  SerialNumber,
		DeviceAddress: DeviceAddress,
		Lat:           Lat,
		Lon:           Lon,
		CreateDate:    panda.CurTime(),
		UpdateDate:    panda.CurTime(),
		CreateBy:      claim.UserId,
		UpdateBy:      claim.UserId,
		DeleteStatus:  0,
	}
	err = r.service.Insert(arg)
	if err != nil {
		hret.Error(resp, 5001001, err.Error())
		return
	}
	hret.Success(resp, "Success")
}

func (r *DeviceInstallController) Put(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	Id := req.FormValue("Id")
	Sid, err := strconv.Atoi(Id)
	if err != nil {
		hret.Error(resp, 500400, err.Error())
		return
	}
	claim, err := jwt.ParseHttp(req)
	if err != nil {
		logger.Error(err)
		hret.Error(resp, 403, "权限不足")
		return
	}

	SerialNumber := req.FormValue("SerialNumber")
	DeviceAddress := req.FormValue("DeviceAddress")
	Lat := req.FormValue("Lat")
	Lon := req.FormValue("Lon")
	arg := entity.DeviceInstallInfo{
		Id:            Sid,
		SerialNumber:  SerialNumber,
		DeviceAddress: DeviceAddress,
		Lat:           Lat,
		Lon:           Lon,
		UpdateDate:    panda.CurTime(),
		UpdateBy:      claim.UserId,
	}
	err = r.service.UpdateById(arg)
	if err != nil {
		hret.Error(resp, 5001001, err.Error())
		return
	}
	hret.Success(resp, "Success")
}

func (r *DeviceInstallController) Delete(resp http.ResponseWriter, req *http.Request, param route.Params) {
	id := param.ByName("id")
	sid, err := strconv.Atoi(id)
	if err != nil {
		hret.Error(resp, 5002001, err.Error())
		return
	}
	_, err = r.service.LogicDeleteById(sid)
	if err != nil {
		hret.Error(resp, 5003001, err.Error())
		return
	}
	hret.Success(resp, "Success")
}

func init() {
	ctl := &DeviceInstallController{
		service: service.NewDeviceInstallService(),
	}
	route.Handler("GET", "/api/device/install", ctl.Get)
	route.Handler("POST", "/api/device/install", ctl.Post)
	route.Handler("PUT", "/api/device/install", ctl.Put)
	route.DELETE("/api/device/install/:id", ctl.Delete)
}
