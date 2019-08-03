package controller

import (
	"ai-platform/api/entity"
	"ai-platform/api/vo"
	"ai-platform/dbobj"
	"ai-platform/panda/hret"
	"ai-platform/panda/jwt"
	"ai-platform/panda/logger"
	"ai-platform/panda/route"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UserController struct {

}

func (r *UserController) Get(resp http.ResponseWriter, req *http.Request)  {
	req.ParseForm()
	claim, err := jwt.ParseHttp(req)
	if err != nil {
		logger.Error(err)
		hret.Error(resp, 403, "权限不足")
		return
	}
	var rst entity.SysUserInfo
	err = dbobj.QueryForStruct("select user_id, nickname, remark, mobile_phone, avatar, weixin, qq, email from sys_user_info where delete_status = 0 and user_id = ?", &rst, claim.UserId)
	if err != nil {
		hret.Error(resp,500100, err.Error())
		return
	}
	hret.Success(resp, rst)
}

func (r *UserController) UpdateProfile(resp http.ResponseWriter, req *http.Request)  {
	req.ParseForm()

	body, err:=ioutil.ReadAll(req.Body)
	if err != nil {
		hret.Error(resp,500030, "添加设备失败，参数为空")
		return
	}
	var arg vo.UserProffileVo
	err = json.Unmarshal(body,&arg)
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
	fmt.Println(arg, claim.UserId)
	result, err := dbobj.Exec("update sys_user_info set nickname = ?, remark = ?, mobile_phone = ?, weixin = ?, qq = ?, email = ? where user_id = ? and delete_status = 0",
		arg.NickName, arg.Remark, arg.MobilePhone, arg.Weixin, arg.QQ, arg.Email, claim.UserId)
	if err != nil {
		hret.Error(resp, 500100,err.Error())
		return
	}
	if size,_ := result.RowsAffected(); size == 1 {
		hret.Success(resp, "Success")
		return
	}
	hret.Success(resp,"信息未发生变更")
}

func (r *UserController) Put(resp http.ResponseWriter, req *http.Request)  {
	req.ParseForm()
	claim, err := jwt.ParseHttp(req)
	if err != nil {
		logger.Error(err)
		hret.Error(resp, 403, "权限不足")
		return
	}
	oldPassword := req.FormValue("OldPassword")
	newPassword := req.FormValue("NewPassword")
	confirmPassword := req.FormValue("ConfirmPassword")
	if newPassword != confirmPassword {
		hret.Error(resp, 500321,"两次输入新密码不匹配，请重新输入")
		return
	}
	result, err := dbobj.Exec("update sys_user_info set password = ? where password = ? and delete_status = 0 and user_id = ?", newPassword, oldPassword, claim.UserId)
	if err != nil{
		hret.Error(resp, 500400, "修改密码失败，请联系管理员")
		return
	}
	if size,_:=result.RowsAffected(); size == 1 {
		hret.Success(resp, "Success")
		return
	}
	hret.Success(resp, "密码未发生修改")
}

func init(){
	ctl := &UserController{};
	route.Handler("GET","/api/account/profiles", ctl.Get)
	route.Handler("PUT","/api/account/password", ctl.Put)
	route.Handler("PUT","/api/account/profiles", ctl.UpdateProfile)
}
