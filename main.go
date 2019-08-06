package main

import (
	"ai-platform/api"
	"ai-platform/api/auth"
	"ai-platform/server"
)

func main() {
	auth.AppRegister("api", api.Register)
	server.AIPlatformBootstrap()
	auth.Bootstrap()
}
