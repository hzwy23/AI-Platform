package platform

import (
	"ai-platform/panda/logger"
	"ai-platform/protocol"
	"net"
	"strconv"
	"time"
)

type Server interface {
	Start()
}

type defaultPlatformServer struct {
	ip       string
	port     int
	protocol string
}

func NewDefaultPlatformServer(ip string, port int, protocol string) Server {
	r := new(defaultPlatformServer)
	r.ip = ip
	r.port = port
	r.protocol = protocol
	return r
}

func UDP(conn *net.UDPConn) {
	udp := defaultPlatformServer{}
	udp.udpServer(conn)
}

func (r *defaultPlatformServer) Start() {
	listen, err := net.Listen(r.protocol, r.convert())
	if err != nil {
		logger.Error("开启监听失败，失败原因是：", err)
	}
	logger.Info("智能灯控平台启动成功")
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			logger.Error("接受请求失败，失败原因是：", err)
		}
		// 有请求到达后，开启协程处理
		logger.Info("客户端发起请求，客户端IP地址是：", conn.RemoteAddr())
		go r.server(conn)
	}
}

func (r *defaultPlatformServer) server(conn net.Conn) {
	protoService := protocol.NewJTTProtocol(conn)
	for {
		message, err := protoService.Parse()
		if err != nil {
			logger.Info("读取内容失败，失败原因是：", err)
			return
		}
		if message == nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		logger.Info("receive message is: ", message)
		context, err := NewContext(protoService, message)
		if err != nil {
			continue
		}
		dispatcher(context)
		time.Sleep(time.Millisecond * 100)
	}
}

func (r *defaultPlatformServer) udpServer(conn *net.UDPConn) {
	protoService := protocol.NewUDPJTTProtocol(conn)
	for {
		message, err := protoService.Parse()
		if err != nil {
			logger.Info("读取内容失败，失败原因是：", err)
			return
		}
		if message == nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		logger.Debug("receive message is: ", message)
		context, err := NewContext(protoService, message)
		if err != nil {
			continue
		}
		dispatcher(context)
		time.Sleep(time.Millisecond * 100)
	}
}

func (r *defaultPlatformServer) convert() string {
	portStr := strconv.Itoa(r.port)
	return r.ip + ":" + portStr
}
