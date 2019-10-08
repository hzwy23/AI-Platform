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
	"strconv"
	"strings"
	"time"
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

    // 参数解析
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error("参数解析错误")
		return
	}
	var tmp interface{}
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		return
	}
	args := make(map[string]interface{})
	for key, val := range tmp.(map[string]interface{}) {
		args[key] = val
	}

	groupId := strconv.Itoa(int(args["GroupId"].(float64)))
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
		r.cmd(item.SerialNumber, args)
		time.Sleep(time.Millisecond * 100)
	}
	hret.Success(resp, "Success")
}

func (r *RemoteDeviceController) cmd(serialNumber string, args map[string]interface{}) {
	// 查看设备是否在线
	item, err := r.deviceDao.FindBySerialNumber(serialNumber)
	if err != nil || len(item.SerialNumber) == 0 {
		fmt.Println("设备未添加到管理列表，无法进行设置操作", serialNumber)
		logger.Warn("设备未添加到管理列表，不能进行修改", serialNumber)
		return
	}


	err = updateDevice(serialNumber, args)
	if err != nil {
		fmt.Println("[统一设置]连接设备：", serialNumber, ",错误信息是：", err)
		logger.Warn("[统一设置]连接设备：", serialNumber, ",错误信息是：", err)
		return
	}


	timer := make([]proto_data.DeviceTimerData,0)
	for _, item := range args["Timer"].([]interface{}){
		element := proto_data.DeviceTimerData{}
		for key, value := range item.(map[string]interface{}) {

			if key == "AutoStartTime" {
				if value == "Invalid date" {
					value = ""
				}
				element.AutoStartTime = value.(string)
			} else if key == "AutoEndTime" {
				if value == "Invalid date" {
					value = ""
				}
				element.AutoEndTime = value.(string)
			}
		}
		timer = append(timer, element)
	}

	cmd := proto_data.DeviceControlData{
		SerialNumber:  serialNumber,
		Timer: timer,
	}

	LightMode := args["LightMode"].(float64)
	cmd.LightMode = strconv.Itoa(int(LightMode))

	err = updateLightMode(serialNumber, cmd, cmd.LightMode)
	if err != nil {
		logger.Warn("连接设备：", serialNumber, ",错误信息是：", err)
		return
	}
	fmt.Println("[统一设置]设备",serialNumber,"统一成功")
}

// UpdateDeviceAttr 更新设备属性
func (r *RemoteDeviceController) UpdateDeviceAttr(resp http.ResponseWriter, req *http.Request, params route.Params) {

	req.ParseForm()
	serialNumber := params.ByName("SerialNumber")
	if len(serialNumber) == 0 {
		hret.Error(resp, 500300, "无效的设备序列号")
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		hret.Error(resp, 500301, "参数解析失败")
		return
	}
	var rst interface{}
	err = json.Unmarshal(body, &rst)
	if err != nil {
		hret.Error(resp, 500302, "参数解析失败")
		return
	}
	args := make(map[string]interface{})
	for key, val := range rst.(map[string]interface{}) {
		args[key] = val
	}

	// 查看设备是否在线
	item, err := r.deviceDao.FindBySerialNumber(serialNumber)
	if err != nil || len(item.SerialNumber) == 0 {
		logger.Warn("设备未添加到管理列表，不能进行修改", serialNumber)
		hret.Error(resp, 500200, "设备未添加到管理列表，不能进行修改")
		return
	}
	if args["Pin"] != item.Pin {
		logger.Warn("设备密码不正确，无法修改设备属性")
		hret.Error(resp, 500300, "设备密码不正确，无法修改设备属性")
		return
	}
	// 更新设备名称
	_, err = dbobj.Exec("update device_manage_info set device_name = ? where delete_status = 0 and serial_number = ?", args["DeviceName"], serialNumber)
	if err != nil {
		logger.Warn(err)
		hret.Error(resp, 500200, "修改设备名称失败，失败原因是："+err.Error())
		return
	}

	err = updateDevice(serialNumber, args)
	if err != nil {
		logger.Error(err)
		hret.Error(resp, 500300, err.Error())
		return
	}

	timer := make([]proto_data.DeviceTimerData,0)
	for _, item := range args["Timer"].([]interface{}){
		element := proto_data.DeviceTimerData{}
		for key, value := range item.(map[string]interface{}) {

			if key == "AutoStartTime" {
				if value == "Invalid date" {
					value = ""
				}
				element.AutoStartTime = value.(string)
			} else if key == "AutoEndTime" {
				if value == "Invalid date" {
					value = ""
				}
				element.AutoEndTime = value.(string)
			}
		}
		timer = append(timer, element)
	}

	cmd := proto_data.DeviceControlData{
		SerialNumber:  serialNumber,
		Timer: timer,
	}

	LightMode := args["LightMode"].(float64)
	cmd.LightMode = strconv.Itoa(int(LightMode))

	err = updateLightMode(serialNumber, cmd, cmd.LightMode)
	if err != nil {
		logger.Error(err)
		hret.Error(resp, 500300, err.Error())
		return
	}

	hret.Success(resp, "Success")
}

// 更新灯光控制策略
func updateLightMode(serialNumber string, cmd proto_data.DeviceControlData, lightMode string) error {
	// 发送指令到终端设备，修改设备属性
	body, _ := json.Marshal(cmd)
	err := utils.Command(0x8005, serialNumber, body)
	if err != nil {
		return err
	}

	// 保存修改后的数据到数据库
	if len(strings.TrimSpace(cmd.LightMode)) == 0 {
		cmd.LightMode = "3"
	}
	startTime := ""
	endTime := ""
	for _, val := range cmd.Timer {
		startTime += val.AutoStartTime + ","
		endTime += val.AutoEndTime+","
	}
	startTime = strings.TrimRight(startTime, ",")
	endTime = strings.TrimRight(endTime, ",")

	_, err = dbobj.Exec("update device_manage_info set light_mode = ?, auto_start_time = ?, auto_end_time = ? where serial_number = ? and delete_status = 0",
		lightMode, startTime, endTime, serialNumber)
	return err
}

// 更新属性
func updateDevice(serialNumber string, args map[string]interface{}) error {
	light := strconv.Itoa(int(args["DeviceBrightness"].(float64)))
	lightVal, yes := proto_data.LightValue[light]
	if !yes {
		lightVal = "1"
	}

	threshold := strconv.Itoa(int(args["DeviceLightThreshold"].(float64)))
	thresholdVal, ok := proto_data.ThresholdValue[threshold]
	if !ok {
		thresholdVal = "1"
	}

	attr := strconv.Itoa(int(args["DeviceAttribute"].(float64)))
	flashDuration := strconv.Itoa(int(args["FlashDuration"].(float64)))

	deviceAttr := DeviceAttrData{
		SerialNumber:         serialNumber,
		DeviceAttribute:      attr,
		DeviceBrightness:     lightVal,
		DeviceLightThreshold: thresholdVal,
		FlashDuration:        flashDuration,
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
