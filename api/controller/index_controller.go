package controller

import (
	"ai-platform/panda/route"
	"html/template"
	"net/http"
)

type indexController struct {
}

// 处理前端路由，前端刷新时，必须获取index.html内容
func (r *indexController) Index(resp http.ResponseWriter, req *http.Request, params route.Params) {
	t, _ := template.ParseFiles("./webui/index.html")
	t.Execute(resp, nil)
}

func init() {
	ctl := indexController{}
	route.GET("/index.html", ctl.Index)
	route.GET("/", ctl.Index)
}
