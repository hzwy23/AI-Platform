package controller

import (
	"ai-platform/api/dao"
	"ai-platform/api/utils"
	"ai-platform/panda/hret"
	"ai-platform/panda/logger"
	"ai-platform/panda/route"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type MotoData struct {
	SerialNumber string `json:"client_CPUID"`
	Command      string `json:"client_MotoControl"`
}

type DeviceNetworkData struct {
	SerialNumber string `json:"client_CPUID"`
	DhcpFlag string `json:"client_NetworkMode"`
	DeviceIp string `json:"client_IP"`
	Mask string `json:"client_MASK"`
	DevicePort string `json:"client_PORT"`
	Gateway string `json:"client_GATEWAY"`
	DNS string `json:"client_DNS"`
	MacAddress string `json:"client_MAC"`
	ServerIp string `json:"server1_IP"`
	SERVERPort string `json:"server1_PORT"`
	BackServerIp string `json:"server2_IP"`
	BackServerPort string `json:"server2_PORT"`
	ServerDomain string `json:"server_ADDR"`
	FirmwareVersion string `json:"client_FrameworkVersion"`
}

type DeviceAttrData struct {
	SerialNumber string `json:"client_CPUID"`
	DeviceAttribute string `json:"client_Mode"`
	DeviceBrightness string `json:"client_LightLevel"`
	DeviceLightThreshold string `json:"client_CDSThreshold"`
}

type DeviceControlData struct {
	SerialNumber string `json:"client_CPUID"`
	LightMode string `json:"client_AutoFunction"`
	AutoStartTime string `json:"AutoTimeStart"`
	AutoEndTime string `json:"AutoTimeStop"`

}

type RemoteDeviceController struct {
	deviceDao dao.DeviceManageInfoDao
}

func (r *RemoteDeviceController) Control(resp http.ResponseWriter, req *http.Request, params route.Params) {
	req.ParseForm()
	serialNumber := params.ByName("SerialNumber")
	if len(serialNumber) == 0 {
		hret.Error(resp, 500300, "无效的设备序列号")
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		hret.Error(resp, 500030, "添加设备失败，参数为空")
		return
	}

	err = utils.Command(0x7fff, serialNumber, body)
	if err != nil {
		hret.Error(resp, 500100, err.Error())
		return
	}

	hret.Success(resp, "Success")

}

func (r *RemoteDeviceController) Plus(resp http.ResponseWriter, req *http.Request, params route.Params) {
	req.ParseForm()
	serialNumber := params.ByName("SerialNumber")
	if len(serialNumber) == 0 {
		hret.Error(resp, 500300, "无效的设备序列号")
		return
	}
	cmd := MotoData{
		SerialNumber: serialNumber,
		Command:      "1",
	}
	body, _ := json.Marshal(cmd)
	err := utils.Command(0x8004, serialNumber, body)
	if err != nil {
		hret.Error(resp, 5001002, err.Error())
		return
	}
	hret.Success(resp, "Success")
}

func (r *RemoteDeviceController) Minus(resp http.ResponseWriter, req *http.Request, params route.Params) {
	req.ParseForm()
	serialNumber := params.ByName("SerialNumber")
	if len(serialNumber) == 0 {
		hret.Error(resp, 500300, "无效的设备序列号")
		return
	}
	cmd := MotoData{
		SerialNumber: serialNumber,
		Command:      "2",
	}
	body, _ := json.Marshal(cmd)
	err := utils.Command(0x8004, serialNumber, body)
	if err != nil {
		hret.Error(resp, 5001002, err.Error())
		return
	}
	hret.Success(resp, "Success")
}


// 修改网络参数
func (r *RemoteDeviceController) UpdateNetwork(resp http.ResponseWriter, req *http.Request, params route.Params)  {

	req.ParseForm()

	serialNumber := params.ByName("SerialNumber")
	if len(serialNumber) == 0 {
		hret.Error(resp, 500300, "无效的设备序列号")
		return
	}

	// 查看设备是否在线
	item, err := r.deviceDao.FindBySerialNumber(serialNumber)
	if err != nil && len(item.SerialNumber) > 0 {
		logger.Warn("设备已添加管理列表，不允许修改网络参数", serialNumber)
		hret.Error(resp, 500200, "设备已添加管理列表，不允许修改网络参数")
		return
	}

	// 指令数据
	cmd := DeviceNetworkData{
		SerialNumber: serialNumber,
		DeviceIp:req.FormValue("DeviceIp"),
		Mask: req.FormValue("Mask"),
		DevicePort: req.FormValue("DevicePort"),
		Gateway: req.FormValue("GatewayAddr"),
		MacAddress: "",
		DNS: "",
		ServerIp: "",
		SERVERPort: "",
		BackServerIp: "",
		BackServerPort: "",
		ServerDomain: "",
		FirmwareVersion: "",
	}
	if req.FormValue("DhcpFlag") == "true" {
		cmd.DhcpFlag = "Auto"
	} else {
		cmd.DhcpFlag = "User"
	}
	body, _ := json.Marshal(cmd)
	err = utils.Command(0x8000, serialNumber, body)
	if err != nil {
		hret.Error(resp, 5001002, err.Error())
		return
	}
	hret.Success(resp, "Success")
}


// 更新设备属性
func (r *RemoteDeviceController) UpdateDeviceAttr(resp http.ResponseWriter, req *http.Request, params route.Params) {

	req.ParseForm()
	serialNumber := params.ByName("SerialNumber")
	if len(serialNumber) == 0 {
		hret.Error(resp, 500300, "无效的设备序列号")
		return
	}

	// 查看设备是否在线
	item, err := r.deviceDao.FindBySerialNumber(serialNumber)
	if err != nil || len(item.SerialNumber) == 0 {
		logger.Warn("设备未添加到管理列表，不能进行修改", serialNumber)
		hret.Error(resp, 500200, "设备未添加到管理列表，不能进行修改")
		return
	}

	err = updateDevice(serialNumber, req.Form)
	if err != nil {
		logger.Error(err)
		hret.Error(resp, 500300, err.Error())
		return
	}
	err = updateLightMode(serialNumber, req.Form)
	if err != nil {
		logger.Error(err)
		hret.Error(resp, 500300, err.Error())
		return
	}
	hret.Success(resp, "Success")
}

// 更新属性
func updateDevice(serialNumber string, form url.Values) error {
	deviceAttr := DeviceAttrData{
		SerialNumber: serialNumber,
		DeviceAttribute: form.Get("DeviceAttribute"),
		DeviceBrightness: form.Get("DeviceBrightness"),
		DeviceLightThreshold: form.Get("DeviceLightThreshold"),
	}
	body, _ := json.Marshal(deviceAttr)
	err := utils.Command(0x8001, serialNumber, body)
	if err != nil {
		return err
	}
	return nil
}

// 更新灯光控制策略
func updateLightMode(serialNumber string, form url.Values) error {

	cmd := DeviceControlData{
		SerialNumber: serialNumber,
		AutoStartTime: form.Get("AutoStartTime"),
		AutoEndTime: form.Get("AutoEndTime"),
	}
	LightMode := form.Get("LightMode")
	if LightMode == "1" {
		cmd.LightMode = "CDS"
	} else if LightMode == "2" {
		cmd.LightMode = "Timer"
	} else if LightMode == "3" {
		cmd.LightMode = "All"
	}
	body, _ := json.Marshal(cmd)
	err := utils.Command(0x8001, serialNumber, body)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	ctl := &RemoteDeviceController{
		deviceDao:dao.NewDeviceManageInfoDao(),
	}
	route.POST("/api/device/remote/control/:SerialNumber", ctl.Control)
	route.GET("/api/device/plus/:SerialNumber", ctl.Plus)
	route.GET("/api/device/minus/:SerialNumber", ctl.Minus)
	route.PUT("/api/device/manage/network/:SerialNumber", ctl.UpdateNetwork)
	route.PUT("/api/device/manage/attribute/:SerialNumber", ctl.UpdateDeviceAttr)
}
