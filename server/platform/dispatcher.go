package platform

type DispatcherService interface {
	// 注册处理器
	Register(msgId uint16, handler Handler)
}
