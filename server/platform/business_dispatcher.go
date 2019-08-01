package platform

import (
	"ai-platform/dbobj"
	"ai-platform/panda"
	"ai-platform/panda/logger"
	"encoding/json"
	"sync"
)

type Handler func(context *Context) (int, string)

type BusinessDispatcher struct {
	// 业务类型对应的处理函数
	msgIdHandler map[uint16]Handler
	lock         *sync.RWMutex
}

var defaultBusinessDispatcher = NewBusinessDispatcher()

func NewBusinessDispatcher() *BusinessDispatcher {
	r := &BusinessDispatcher{
		msgIdHandler: make(map[uint16]Handler),
		lock:         new(sync.RWMutex),
	}
	return r
}

func (r *BusinessDispatcher) Register(msgId uint16, handler Handler) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.msgIdHandler[msgId] = handler
}

func (r *BusinessDispatcher) dispatcher(context *Context) (int, string){
	r.lock.RLock()
	defer r.lock.RUnlock()
	if handler, ok := r.msgIdHandler[context.GetMsgId()]; ok {
		return handler(context)
	} else {
		logger.Error("无效的业务类型, 业务类型是：", context.GetMsgId())
		return 500, "无效的业务类型"
	}
}

func Register(msgId uint16, handler Handler) {
	defaultBusinessDispatcher.Register(msgId, handler)
}

func dispatcher(context *Context) {
	code, retMsg := defaultBusinessDispatcher.dispatcher(context)
	go func() {
		msgId := context.msgId
		msg := context.message.MsgBody
		var rst interface{}
		json.Unmarshal(msg, &rst)
		body := rst.(map[string]interface{})
		logger.Info("报文内容是：",rst)
		bodyStr,_ := json.Marshal(body)
		result, err := dbobj.Exec("insert into plat_device_logger(serial_number, handle_time, direction, biz_type, message, ret_code, ret_msg) values(?, ?, ?, ?, ?, ?, ?)",
			body["client_CPUID"], panda.CurTime(), "Input", msgId, bodyStr, code, retMsg)
		if err != nil {
			logger.Error(result, err, *context.message)
		}
	}()
}
