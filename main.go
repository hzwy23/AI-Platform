package main

import (
	"ai-platform/api"
	"ai-platform/api/auth"
	_ "ai-platform/panda/cron"
	"ai-platform/server"
)

func main() {
	auth.AppRegister("api", api.Register)
	server.AIPlatformBootstrap()
	auth.Bootstrap()
}
