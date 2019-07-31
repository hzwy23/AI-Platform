package controller

import (
	"ai-platform/api/service"
	"ai-platform/panda/hret"
	"ai-platform/panda/logger"
	"ai-platform/panda/route"
	"encoding/json"
	"net/http"
	"strconv"
)

type DeviceScanController struct {
	deviceScan service.DeviceScanService
}

func (r *DeviceScanController) Get(resp http.ResponseWriter, req *http.Request) {
	dt, idx, err := r.deviceScan.FindAll()
	if err != nil {
		logger.Warn(err)
		hret.Error(resp, 5002001, err.Error())
	}

	ok := hret.RetContent{
		Code:    200,
		Message: "execute successfully.",
		Rows:    dt,
		Total:   int64(idx),
	}

	ijs, err := json.Marshal(ok)
	if err != nil {
		logger.Error(err)
		resp.WriteHeader(http.StatusExpectationFailed)
		resp.Write([]byte(`{"code":"` + strconv.Itoa(http.StatusExpectationFailed) + `","msg":"format json type info failed.","details":"format json type info failed."}`))
		return
	}
	resp.Write(ijs)
}

func init() {
	scan := &DeviceScanController{
		deviceScan: service.NewDeviceScanService(),
	}
	route.Handler("GET", "/api/scan/device", scan.Get)
}
