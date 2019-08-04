package controller

import (
	"ai-platform/panda/route"
	"fmt"
	"net/http"
)

type indexController struct {

}

func (r *indexController)Index(resp http.ResponseWriter, req *http.Request, params route.Params)  {
	fmt.Println(params, req.RequestURI)
}

func init()  {
	ctl := &indexController{}
	route.GET("/", ctl.Index)
}
