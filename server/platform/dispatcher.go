package platform

import (
	"ai-platform/api/entity"
	"ai-platform/dbobj"
	"ai-platform/panda"
	"ai-platform/panda/hret"
	"ai-platform/panda/logger"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type DispatcherService interface {
	// 注册处理器
	Register(msgId uint16, handler Handler)
}

type Handler func(context *Context) (int, string)

type BusinessDispatcher struct {
	// 业务类型对应的处理函数
	msgIdHandler map[uint16]Handler
	lock         *sync.RWMutex
}

var logDeviceBuf = make(chan entity.PlatDeviceLogger, 40960)

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

func (r *BusinessDispatcher) dispatcher(context *Context) (int, string) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if handler, ok := r.msgIdHandler[context.GetMsgId()]; ok {
		return handler(context)
	} else {
		logger.Warn("无效的业务类型, 业务类型是：", context.GetMsgId())
		return 500, "无效的业务类型"
	}
}

func Register(msgId uint16, handler Handler) {
	defaultBusinessDispatcher.Register(msgId, handler)
}

func dispatcher(context *Context) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	code, retMsg := defaultBusinessDispatcher.dispatcher(context)
	go genLog(context.msgId, context.message.MsgBody, code, retMsg)
}

func genLog(msgId uint16, msg []byte, code int, retMsg string) {
	defer hret.RecvPanic()

	if msg == nil || len(msg) == 0 {
		logger.Error("收到无效的消息")
	} else {
		var rst interface{}
		json.Unmarshal(msg, &rst)
		body := rst.(map[string]interface{})
		logger.Debug("报文内容是：", rst)
		bodyStr, _ := json.Marshal(body)

		item := entity.PlatDeviceLogger{
			Direction:    "Input",
			BizType:      strconv.Itoa(int(msgId)),
			Message:      string(bodyStr),
			RetCode:      strconv.Itoa(code),
			RetMsg:       retMsg,
			SerialNumber: body["client_CPUID"].(string),
			HandleTime:   panda.CurTime(),
		}
		logger.Debug("消息是：", item)
		logDeviceBuf <- item
	}
}

func logSync() {
	var buf []entity.PlatDeviceLogger
	for {
		select {
		case <-time.After(time.Second * 2):
			// sync handle logs to database per 5 second.
			if len(buf) == 0 {
				continue
			}
			go savelogs(buf)
			buf = make([]entity.PlatDeviceLogger, 0)
		case val, ok := <-logDeviceBuf:
			if ok {
				buf = append(buf, val)
				if len(buf) > 20 {
					go savelogs(buf)
					buf = make([]entity.PlatDeviceLogger, 0)
				}
			}
		}
	}
}

func savelogs(data []entity.PlatDeviceLogger) {
	for _, item := range data {
		result, err := dbobj.Exec("insert into plat_device_logger(serial_number, handle_time, direction, biz_type, message, ret_code, ret_msg) values(?, ?, ?, ?, ?, ?, ?)",
			item.SerialNumber, item.HandleTime, item.Direction, item.BizType, item.Message, item.RetCode, item.RetMsg)
		if err != nil {
			logger.Error(result, err, item)
		}
	}
}

func init() {
	go logSync()
}
