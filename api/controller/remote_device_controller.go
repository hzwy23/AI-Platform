package controller

import (
	"ai-platform/api/dao"
	"ai-platform/dbobj"
	"ai-platform/panda"
	"ai-platform/panda/hret"
	"ai-platform/panda/jwt"
	"ai-platform/panda/logger"
	"ai-platform/panda/route"
	"ai-platform/server/proto_data"
	"ai-platform/server/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type MotoData struct {
	SerialNumber string `json:"client_CPUID"`
	Command      string `json:"client_MotoControl"`
}

type DeviceNetworkData struct {
	SerialNumber    string `json:"client_CPUID"`
	DhcpFlag        string `json:"client_NetworkMode"`
	DeviceIp        string `json:"client_IP"`
	Mask            string `json:"client_MASK"`
	DevicePort      string `json:"client_PORT"`
	Gateway         string `json:"client_GATEWAY"`
	DNS             string `json:"client_DNS"`
	MacAddress      string `json:"client_MAC"`
	ServerIp        string `json:"server1_IP"`
	ServerPort      string `json:"server1_PORT"`
	BackServerIp    string `json:"server2_IP"`
	BackServerPort  string `json:"server2_PORT"`
	ServerDomain    string `json:"server_ADDR"`
	FirmwareVersion string `json:"client_FrameworkVersion"`
}

type DeviceAttrData struct {
	SerialNumber         string `json:"client_CPUID"`
	DeviceAttribute      string `json:"client_Mode"`
	DeviceBrightness     string `json:"client_LightLevel"`
	DeviceLightThreshold string `json:"client_CDSThreshold"`
	FlashDuration        string `json:"client_FlashDuration"`
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

// UpdateNetwork 修改网络参数
func (r *RemoteDeviceController) UpdateNetwork(resp http.ResponseWriter, req *http.Request, params route.Params) {
	req.ParseForm()

	if req.FormValue("Pin") != "123456" {
		hret.Error(resp, 500300, "设备密码不正确")
		return
	}

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
		SerialNumber:    serialNumber,
		DeviceIp:        req.FormValue("DeviceIp"),
		Mask:            req.FormValue("Mask"),
		DevicePort:      req.FormValue("DevicePort"),
		Gateway:         req.FormValue("GatewayAddr"),
		MacAddress:      "",
		DNS:             "",
		ServerIp:        "",
		ServerPort:      "",
		BackServerIp:    "",
		BackServerPort:  "",
		ServerDomain:    "",
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

// GlobalSetting 分组内设备统一配置
func (r *RemoteDeviceController) GlobalSetting(resp http.ResponseWriter, req *http.Request) {

	claim, err := jwt.ParseHttp(req)
	if err != nil {
		logger.Error(err)
		hret.Error(resp, 403, "权限不足")
		return
	}

	if !panda.IsAdmin(claim.UserId) {
		hret.Error(resp, 500500, "只有超级管理员才能进行统一设置")
		return
	}

	req.ParseForm()
	groupId := req.FormValue("GroupId")
	if len(groupId) == 0 {
		hret.Error(resp, 500300, "设备分组为空，不能进行全局设置")
		return
	}

	var rst []DeviceAttrData
	err = dbobj.QueryForSlice("select i.serial_number from group_device_bind t inner join device_manage_info i on t.device_id = i.device_id and i.delete_status = 0 where t.delete_status = 0 and t.group_id = ?", &rst, groupId)
	if err != nil {
		fmt.Println(err)
		hret.Error(resp, 500400, "设备组内没有查询到有效的设备信息")
		return
	}
	for _, item := range rst {
		r.cmd(item.SerialNumber, req)
	}
	hret.Success(resp, "Success")
}

func (r *RemoteDeviceController) cmd(serialNumber string, req *http.Request) {
	// 查看设备是否在线
	item, err := r.deviceDao.FindBySerialNumber(serialNumber)
	if err != nil || len(item.SerialNumber) == 0 {
		logger.Warn("设备未添加到管理列表，不能进行修改", serialNumber)
		return
	}

	err = updateDevice(serialNumber, req.Form)
	if err != nil {
		logger.Warn("连接设备：", serialNumber, ",错误信息是：", err)
		return
	}

	cmd := proto_data.DeviceControlData{
		SerialNumber:  serialNumber,
		AutoStartTime: req.FormValue("AutoStartTime"),
		AutoEndTime:   req.FormValue("AutoEndTime"),
	}
	if cmd.AutoEndTime == "Invalid date" {
		cmd.AutoEndTime = ""
	}
	if cmd.AutoStartTime == "Invalid date" {
		cmd.AutoStartTime = ""
	}

	LightMode := req.FormValue("LightMode")
	if LightMode == "1" {
		cmd.LightMode = "CDS"
	} else if LightMode == "2" {
		cmd.LightMode = "Timer"
	} else if LightMode == "3" {
		cmd.LightMode = "All"
	} else {
		cmd.LightMode = "All"
		LightMode = "3"
	}

	err = updateLightMode(serialNumber, cmd, LightMode)
	if err != nil {
		logger.Warn("连接设备：", serialNumber, ",错误信息是：", err)
		return
	}

}

// UpdateDeviceAttr 更新设备属性
func (r *RemoteDeviceController) UpdateDeviceAttr(resp http.ResponseWriter, req *http.Request, params route.Params) {

	req.ParseForm()
	serialNumber := params.ByName("SerialNumber")
	deviceName := req.FormValue("DeviceName")
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
	if req.FormValue("Pin") != item.Pin {
		logger.Warn("设备密码不正确，无法修改设备属性")
		hret.Error(resp, 500300, "设备密码不正确，无法修改设备属性")
		return
	}
	// 更新设备名称
	_, err = dbobj.Exec("update device_manage_info set device_name = ? where delete_status = 0 and serial_number = ?", deviceName, serialNumber)
	if err != nil {
		logger.Warn(err)
		hret.Error(resp, 500200, "修改设备名称失败，失败原因是："+err.Error())
		return
	}

	err = updateDevice(serialNumber, req.Form)
	if err != nil {
		logger.Error(err)
		hret.Error(resp, 500300, err.Error())
		return
	}

	cmd := proto_data.DeviceControlData{
		SerialNumber:  serialNumber,
		AutoStartTime: req.FormValue("AutoStartTime"),
		AutoEndTime:   req.FormValue("AutoEndTime"),
	}
	if cmd.AutoEndTime == "Invalid date" {
		cmd.AutoEndTime = ""
	}
	if cmd.AutoStartTime == "Invalid date" {
		cmd.AutoStartTime = ""
	}

	LightMode := req.FormValue("LightMode")
	cmd.LightMode = LightMode

	err = updateLightMode(serialNumber, cmd, LightMode)
	if err != nil {
		logger.Error(err)
		hret.Error(resp, 500300, err.Error())
		return
	}

	hret.Success(resp, "Success")
}

// 更新灯光控制策略
func updateLightMode(serialNumber string, cmd proto_data.DeviceControlData, lightMode string) error {
	body, _ := json.Marshal(cmd)
	err := utils.Command(0x8005, serialNumber, body)
	if err != nil {
		return err
	}
	if len(strings.TrimSpace(lightMode)) == 0 {
		lightMode = "3"
	}
	_, err = dbobj.Exec("update device_manage_info set light_mode = ?, auto_start_time = ?, auto_end_time = ? where serial_number = ? and delete_status = 0",
		lightMode, cmd.AutoStartTime, cmd.AutoEndTime, serialNumber)
	return err
}

// 更新属性
func updateDevice(serialNumber string, form url.Values) error {

	lightVal, yes := proto_data.LightValue[form.Get("DeviceBrightness")]
	if !yes {
		lightVal = "1"
	}

	thresholdVal, ok := proto_data.ThresholdValue[form.Get("DeviceLightThreshold")]
	if !ok {
		thresholdVal = "1"
	}

	deviceAttr := DeviceAttrData{
		SerialNumber:         serialNumber,
		DeviceAttribute:      form.Get("DeviceAttribute"),
		DeviceBrightness:     lightVal,
		DeviceLightThreshold: thresholdVal,
		FlashDuration:        form.Get("FlashDuration"),
	}

	body, _ := json.Marshal(deviceAttr)
	err := utils.Command(0x8001, serialNumber, body)
	if err != nil {
		return err
	}

	_, err = dbobj.Exec("update device_manage_info set device_attribute = ?, device_brightness = ?, device_light_threshold = ?, flash_duration = ? where delete_status = 0 and serial_number = ?",
		deviceAttr.DeviceAttribute,
		proto_data.LightValueReserve[deviceAttr.DeviceBrightness],
		proto_data.ThresholdValueReserve[deviceAttr.DeviceLightThreshold],
		deviceAttr.FlashDuration,
		deviceAttr.SerialNumber)

	return err
}

func init() {
	ctl := &RemoteDeviceController{
		deviceDao: dao.NewDeviceManageInfoDao(),
	}
	route.POST("/api/device/remote/control/:SerialNumber", ctl.Control)
	route.GET("/api/device/plus/:SerialNumber", ctl.Plus)
	route.GET("/api/device/minus/:SerialNumber", ctl.Minus)
	route.PUT("/api/device/manage/network/:SerialNumber", ctl.UpdateNetwork)
	route.PUT("/api/device/manage/attribute/:SerialNumber", ctl.UpdateDeviceAttr)
	route.Handler("PUT", "/api/device/global/setting", ctl.GlobalSetting)

}
