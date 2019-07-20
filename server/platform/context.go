package platform

import (
	"ai-platform/panda/logger"
	"ai-platform/protocol"
	"errors"
	"fmt"
	"sync"
)

type Context struct {
	msgId        uint16
	protoService protocol.CommunicationService
	message      *protocol.Message
	lock         *sync.RWMutex
}

func NewContext(protoService protocol.CommunicationService, data []byte) (*Context, error) {
	if len(data) == 0 {
		return nil, errors.New("消息格式不正确")
	}
	msg := protocol.ConvertMessage(data)
	if msg == nil {
		fmt.Println("异常数据格式：", data)
		return nil, errors.New("返回消息格式不不正确")
	}
	logger.Info("convert message is: ", *msg)
	r := &Context{
		message:      msg,
		protoService: protoService,
		msgId:        msg.MsgHeader.MsgID,
		lock:         new(sync.RWMutex),
	}
	return r, nil
}

func (r *Context) GetMsgId() uint16 {
	return r.msgId
}

func (r *Context) GetMessage() *protocol.Message {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.message
}

func (r *Context) Send(msg *protocol.Message) (int, error) {
	msgByte := protocol.ConvertToByte(msg)
	return r.protoService.Send(msgByte)
}
