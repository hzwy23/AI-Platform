package controller

import (
	"ai-platform/api/service"
	"ai-platform/panda/hret"
	"ai-platform/panda/jwt"
	"ai-platform/panda/logger"
	"ai-platform/panda/route"
	"fmt"
	"net/http"
	"strconv"
)

type GroupDeviceController struct {
	groupDeviceService service.GroupDeviceService
}

func (r *GroupDeviceController) Get(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	data, err := r.groupDeviceService.FindByPage(1, 10)
	if err != nil {
		logger.Error("请求数据失败，失败原因是：", err)
		hret.Error(resp, 5000020, err.Error())
		return
	}
	hret.Success(resp, data)
}

func (r *GroupDeviceController) Post(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	form := req.PostForm

	groupName := form.Get("GroupName")
	claim, err := jwt.ParseHttp(req)
	if err != nil {
		logger.Error(err)
		hret.Error(resp, 403, "权限不足")
		return
	}

	err = r.groupDeviceService.AddGroup(groupName, claim.UserId)
	if err != nil {
		hret.Error(resp, 50000001, err.Error(), groupName)
		return
	}
	hret.Success(resp, "OK")
}

func (r *GroupDeviceController) Put(resp http.ResponseWriter, req *http.Request) {

}

func (r *GroupDeviceController) Delete(resp http.ResponseWriter, req *http.Request, param route.Params) {
	fmt.Println(param)
	groupId := param.ByName("groupId")
	id, err := strconv.Atoi(groupId)
	if err != nil {
		logger.Info("分组编码有无，无法执行删除操作")
		hret.Error(resp, 500011, "删除错误")
		return
	}
	size, err := r.groupDeviceService.DeleteByGroupId(id)
	if err != nil {
		logger.Info(err)
		hret.Error(resp, 500012, err.Error())
		return
	}
	if size == 0 {
		logger.Error("分组不存在，不能执行删除操作")
		hret.Error(resp, 500013, "分组不存在，不能执行删除操作")
		return
	}
	hret.Success(resp, "success")
}

func init() {
	logger.Info("注册分组控制模块")
	handleFunc := &GroupDeviceController{
		groupDeviceService: service.NewGroupDeviceService(),
	}
	route.Handler("GET", "/api/device/group", handleFunc.Get)
	route.Handler("POST", "/api/device/group", handleFunc.Post)
	route.Handler("PUT", "/api/device/group", handleFunc.Put)
	route.DELETE("/api/device/group/:groupId", handleFunc.Delete)

}
