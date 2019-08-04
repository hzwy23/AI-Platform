package utils

import (
	"ai-platform/dbobj"
	"ai-platform/panda"
	"ai-platform/panda/logger"
	"ai-platform/protocol"
	"ai-platform/server/service"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

func Command(msgId uint16, serialNumber string, command []byte) error {
	if item, ok := service.GetOnlineDevice()[serialNumber]; ok {
		address := item.DeviceIp + ":" + item.DevicePort
		tcpAddr, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			logger.Error("", err)
			go writeLog(msgId, command, 10010, err.Error())
			return err
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			logger.Error("连接失败，", err)
			go writeLog(msgId, command, 10011, err.Error())
			return err
		}
		defer conn.Close()
		data, err := protocol.Pack(msgId, command)
		logger.Info("发送控制指令到设备，指令时：", data)
		size, err := conn.Write(data)
		if err != nil {
			logger.Error(err.Error())
			go writeLog(msgId, command, 10012, err.Error())
			return err
		}
		if size == len(data) {
			logger.Info("向设备发送指令成功，指令内容是：", data)
			go writeLog(msgId, command, 200, "Success")
			return nil
		}
		go writeLog(msgId, command, 10013, "指令发送过程中出现丢包")
		return errors.New("发送指令长度不一致，发送指令格式错误")
	}
	go writeLog(msgId, command, 10014, "设备不在线，无法发送控制指令")
	return errors.New("设备不在线，无法发送控制指令")
}

func writeLog(msgId uint16, msg []byte, retCode int, retMsg string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	if msg == nil || len(msg) == 0 {
		logger.Error("收到无效的消息")
	} else {
		var rst interface{}
		json.Unmarshal(msg, &rst)
		body := rst.(map[string]interface{})
		logger.Debug("报文内容是：", rst)
		bodyStr, _ := json.Marshal(body)
		result, err := dbobj.Exec("insert into plat_device_logger(serial_number, handle_time, direction, biz_type, message, ret_code, ret_msg) values(?, ?, ?, ?, ?, ?, ?)",
			body["client_CPUID"], panda.CurTime(), "Output", msgId, bodyStr, retCode, retMsg)
		if err != nil {
			logger.Error(result, err, msg)
		}
	}
}
