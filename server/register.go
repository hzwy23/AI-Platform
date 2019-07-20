package server

import (
	"ai-platform/panda/config"
	"ai-platform/panda/logger"
	"ai-platform/server/platform"
	_ "ai-platform/server/service"
	"strconv"
)

func AIPlatformBootstrap() {
	// 启动服务
	// 从配置文件获取端口号
	c, err := config.Load("conf/app.conf", config.INI)
	if err != nil {
		logger.Warn("读取配置文件conf/app.conf文件失败，使用默认参数启动服务")
		server := platform.NewDefaultPlatformServer("", 8989, "tcp")
		server.Start()
		return
	}

	configPort, err := c.Get("ai.platform.port")
	if err != nil {
		logger.Warn("http.port 不存在，使用默认端口[8989]启动服务")
		configPort = "8989"
	}
	port, err := strconv.Atoi(configPort)
	if err != nil {
		logger.Warn("无效的配置端口", configPort, "，使用默认端口[ 8989 ]启动服务")
		port = 8989
	}

	logger.Info("AI智能灯控平台启动，监听端口是：", port)
	server := platform.NewDefaultPlatformServer("", port, "tcp")
	go server.Start()

}
