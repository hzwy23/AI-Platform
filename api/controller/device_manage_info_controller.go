package controller

import (
	"ai-platform/api/entity"
	"ai-platform/api/service"
	"ai-platform/api/vo"
	"ai-platform/dbobj"
	"ai-platform/panda"
	"ai-platform/panda/hret"
	"ai-platform/panda/jwt"
	"ai-platform/panda/logger"
	"ai-platform/panda/route"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type DeviceManageInfoController struct {
	service service.DeviceManageService
}

const sqlText = `SELECT
	t.device_id,
	t.serial_number,
	t.device_name,
	t.device_ip,
	t.device_port 
FROM
	device_manage_info t 
WHERE
	t.delete_status = 0 
	AND not EXISTS ( SELECT 1 FROM group_device_bind b where t.device_id = b.device_id and delete_status = 0 )`

func (r *DeviceManageInfoController) Get(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	groupId := req.FormValue("GroupId")
	rst, err := r.service.FindAll(groupId)
	if err != nil {
		hret.Error(resp, 5003445, err.Error())
		return
	}
	hret.Success(resp, rst)
}

// 添加设备
// 从扫描到的设备中，选择制定设备，添加到设备管理列表
func (r *DeviceManageInfoController) Post(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	//dchpFlag := req.FormValue("DhcpFlag")
	//dhcp, err := strconv.Atoi(dchpFlag)
	//if err != nil {
	//	hret.Error(resp,500100, err.Error())
	//	return
	//}
	//attr := req.FormValue("DeviceAttribute")
	//attribute, err := strconv.Atoi(attr)
	//if err != nil {
	//	hret.Error(resp,500101, err.Error())
	//	return
	//}

	claim, err := jwt.ParseHttp(req)
	if err != nil {
		logger.Error(err)
		hret.Error(resp, 403, "权限不足")
		return
	}

	item := &entity.DeviceManageInfo{
		// 设备序列号
		SerialNumber: req.FormValue("SerialNumber"),
		// 设备名称
		DeviceName: req.FormValue("DeviceName"),
		// 使用自动获取IP，0：否，1：是
		//DhcpFlag: uint8(dhcp),
		// 设备IP
		DeviceIp: req.FormValue("DeviceIp"),
		// 设备服务端口
		DevicePort:   req.FormValue("DevicePort"),
		DeviceStatus: 1,
		// 设备属性
		//DeviceAttribute: uint8(attribute),

		// MAC地址
		MacAddress: req.FormValue("MacAddress"),
		// 固件版本
		FirmwareVersion: req.FormValue("FirmwareVersion"),

		// 掩码
		Mask: req.FormValue("Mask"),
		// 网关
		Gateway:      req.FormValue("Gateway"),
		CreateBy:     claim.UserId,
		CreateDate:   panda.CurTime(),
		UpdateBy:     claim.UserId,
		UpdateDate:   panda.CurTime(),
		DeleteStatus: 0,
	}
	err = r.service.AddDevice(item, req.FormValue("GroupId"))
	if err != nil {
		hret.Error(resp, 500300, err.Error())
		return
	}

	hret.Success(resp, "Success")
}

func (r *DeviceManageInfoController) Put(resp http.ResponseWriter, req *http.Request) {

}

func (r *DeviceManageInfoController) Delete(resp http.ResponseWriter, req *http.Request, param route.Params) {
	deviceId := param.ByName("deviceId")
	sid, err := strconv.Atoi(deviceId)
	if err != nil {
		logger.Error("无效的设备号")
		hret.Error(resp, 500300, deviceId)
		return
	}
	err = r.service.RemoveDevice(sid)
	if err != nil {
		hret.Error(resp, 500500, err.Error())
		return
	}
	hret.Success(resp, "Success")
}

func (r *DeviceManageInfoController) RemoveFromGroup(resp http.ResponseWriter, req *http.Request, param route.Params) {
	deviceId := param.ByName("id")
	sid, err := strconv.Atoi(deviceId)
	if err != nil {
		logger.Error("无效的设备号")
		hret.Error(resp, 500300, deviceId)
		return
	}
	_, err = dbobj.Exec("update group_device_bind set delete_status = 1 where device_id = ?", sid)
	if err != nil {
		hret.Error(resp, 500500, err.Error())
		return
	}
	hret.Success(resp, "Success")
}

// 查询所有没有分组的设备
func (r *DeviceManageInfoController) GetUnGroupDevice(resp http.ResponseWriter, req *http.Request) {
	rst := make([]vo.UnGroupDeviceVo,0)
	err := dbobj.QueryForSlice(sqlText, &rst)
	if err != nil {
		hret.Error(resp, 500300,"查询未分组设备失败")
		return
	}
	hret.Success(resp, rst)
}

func (r *DeviceManageInfoController) UpdateDeviceGroup(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	body, err:=ioutil.ReadAll(req.Body)
	if err != nil {
		hret.Error(resp,500030, "添加设备失败，参数为空")
		return
	}
	args := make([]*vo.UnGroupDeviceVo,0)
	err = json.Unmarshal(body,&args)
	if err != nil {
		fmt.Println(err)
		hret.Error(resp,500031, "添加设备失败，参数解析失败")
		return
	}
	claim, err := jwt.ParseHttp(req)
	if err != nil {
		logger.Error(err)
		hret.Error(resp, 403, "权限不足")
		return
	}

	tx,_ := dbobj.Begin()
	for _, item := range args {
		result, err := tx.Exec("insert into group_device_bind(group_id, device_id, create_by, create_date, update_by, update_date, delete_status) values(?,?,?,?,?,?,0)",
			item.GroupId, item.DeviceId, claim.UserId, panda.CurTime(), claim.UserId, panda.CurTime())
		if err != nil {
			tx.Rollback()
			hret.Error(resp, 50030, "新增设备失败,请联系管理员")
			return
		}
		if result == nil {
			tx.Rollback()
			hret.Error(resp, 500308, "新增设备失败,请联系管理员")
			return
		}
	}
	tx.Commit()
	hret.Success(resp,"Success")
}

func (r *DeviceManageInfoController) ChangeGroup(resp http.ResponseWriter, req *http.Request)  {
	req.ParseForm()
	deviceId := req.FormValue("DeviceId")
	groupId := req.FormValue("GroupId")
	claim, err := jwt.ParseHttp(req)
	if err != nil {
		logger.Error(err)
		hret.Error(resp, 403, "权限不足")
		return
	}

	dbobj.Exec("update group_device_bind set delete_status = 1 where device_id = ?", deviceId)
	dbobj.Exec("insert into group_device_bind(group_id, device_id, create_by, create_date, update_by, update_date, delete_status) values(?,?,?,?,?,?,0)",
		groupId, deviceId, claim.UserId, panda.CurTime(), claim.UserId, panda.CurTime())
	hret.Success(resp,"Success")
}

func init() {
	ctl := &DeviceManageInfoController{
		service: service.NewDeviceManageService(),
	}
	route.Handler("GET", "/api/device/manage", ctl.Get)
	route.Handler("GET", "/api/device/manage/ungroup", ctl.GetUnGroupDevice)
	route.Handler("POST", "/api/device/manage", ctl.Post)
	route.Handler("POST", "/api/device/manage/group", ctl.UpdateDeviceGroup)
	route.Handler("PUT", "/api/device/manage", ctl.Put)
	route.Handler("PUT","/api/device/group/change", ctl.ChangeGroup)
	route.DELETE("/api/device/manage/:deviceId", ctl.Delete)
	route.DELETE("/api/device/bind/:id", ctl.RemoveFromGroup)
}
