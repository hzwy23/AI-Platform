package main

import (
	"ai-platform/api"
	"ai-platform/api/auth"
	_ "ai-platform/panda/cron"
	"ai-platform/server"
	"fmt"
	"time"
)

func main() {
	auth.AppRegister("api", api.Register)
	server.AIPlatformBootstrap()
	go func() {
		for {
			fmt.Println("*********************************************************************************************************************************")
			fmt.Println("*  欢迎使用智能灯控平台，当前版本为开发测试版，不可用户生产环境，若擅自将当前版本用于生产，我方不承当任何责任，后果自负。\t*")
			fmt.Println("*  欢迎使用智能灯控平台，当前版本为开发测试版，不可用户生产环境，若擅自将当前版本用于生产，我方不承当任何责任，后果自负。\t*")
			fmt.Println("*  欢迎使用智能灯控平台，当前版本为开发测试版，不可用户生产环境，若擅自将当前版本用于生产，我方不承当任何责任，后果自负。\t*")
			fmt.Println("*  \t\t\t\t                                             \t\t\t\t\t\t\t*")
			fmt.Println("*  \t\t\t\t          邮件联系：   hzwy23@163.com           \t\t\t\t\t\t*")
			fmt.Println("*  \t\t\t\t          微信联系：   hzwy23                   \t\t\t\t\t\t*")
			fmt.Println("*  \t\t\t\t          电话联系：   18581530028                 \t\t\t\t\t\t*")
			fmt.Println("*  \t\t\t\t                                             \t\t\t\t\t\t\t*")
			fmt.Println("*  欢迎使用智能灯控平台，当前版本为开发测试版，不可用户生产环境，若擅自将当前版本用于生产，我方不承当任何责任，后果自负。\t*")
			fmt.Println("*  欢迎使用智能灯控平台，当前版本为开发测试版，不可用户生产环境，若擅自将当前版本用于生产，我方不承当任何责任，后果自负。\t*")
			fmt.Println("*  欢迎使用智能灯控平台，当前版本为开发测试版，不可用户生产环境，若擅自将当前版本用于生产，我方不承当任何责任，后果自负。\t*")
			fmt.Println("*********************************************************************************************************************************")
			time.Sleep(time.Second * 30)
		}
	}()
	auth.Bootstrap()
}
