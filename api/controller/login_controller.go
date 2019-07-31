package controller

import (
	"ai-platform/dbobj"
	"ai-platform/panda/hret"
	"ai-platform/panda/jwt"
	"ai-platform/panda/route"
	"net/http"
)

type LoginController struct {
}

func (login *LoginController) Post(w http.ResponseWriter, r *http.Request) {

	_ = r.ParseForm()

	form := r.Form

	username := r.FormValue("username")
	password := r.FormValue("password")
	pd := ""
	err := dbobj.QueryForObject("select password from sys_user_info where user_id = ?", dbobj.PackArgs(username), &pd)
	if err != nil {
		hret.Error(w, 403,"用户不存在")
		return
	}
	if password != pd {
		hret.Error(w,401,"用户密码错误，请重新输入")
		return
	}

	token, _ := jwt.GenToken(jwt.NewUserdata().SetUserId(form.Get("username")))

	cookie := http.Cookie{Name: "Authorization", Value: token, Path: "/", MaxAge: int(172800)}

	http.SetCookie(w, &cookie)

	_ = hret.Success(w, token)

}

func init() {
	login := &LoginController{}
	route.Handler("POST", "/api/login", login.Post)
}
