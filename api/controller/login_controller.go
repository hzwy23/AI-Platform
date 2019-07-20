package controller

import (
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

	token, _ := jwt.GenToken(jwt.NewUserdata().SetUserId(form.Get("username")))

	cookie := http.Cookie{Name: "Authorization", Value: token, Path: "/", MaxAge: int(17280)}

	http.SetCookie(w, &cookie)

	_ = hret.Success(w, token)

}

func init() {
	login := &LoginController{}
	route.Handler("POST", "/api/login", login.Post)
}
