package main

import (
	"ai-platform/panda/logger"
	"ai-platform/protocol"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type DeviceInfo struct {
	// 设别序列号
	SerialNumber string
	// 软件版本号
	FirmwareVersion string
	// 设备IP地址
	DeviceIp string
	// 设备掩码
	Mask string
	// 网关地址
	GatewayAddr string
	// 设备端口号
	DevicePort string
	// 设备mac地址
	MacAddr string
}

func reconnect() (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "121.42.143.130:8989")
	if err != nil {
		logger.Error("", err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logger.Error("连接失败，", err)
	}
	return conn, err
}

func write(conn *net.TCPConn, message []byte) {
	idx := 0
	for {
		_, err := conn.Write(message[idx : idx+1])
		fmt.Println("发送数据：", message[idx:idx+1])
		idx += 1
		if len(message) == idx {
			idx = 0
		}
		if err != nil {
			logger.Error(err)
			conn, err = reconnect()
			if err != nil {
				for {
					time.Sleep(time.Second * 3)
					conn, err = reconnect()
					if err == nil {
						break
					}
				}
			}
		}

		time.Sleep(time.Millisecond * 100)

	}
}

func read(conn *net.TCPConn) {
	for {
		rbuf := make([]byte, 128)
		rsize, err := conn.Read(rbuf)

		buf,err := protocol.UnPack(rbuf[:rsize])
		if err != nil {
			return
		}
		msg:= protocol.ConvertMessage(buf)

		fmt.Println("读取到字节流是：", rbuf[:rsize])
		fmt.Println("读取到字节流是：", buf)
		fmt.Println("读取到对象是：", msg)
		fmt.Println("数据body是：",string(msg.MsgBody))
		if err != nil {
			logger.Error(err)
			conn, err = reconnect()
			if err != nil {
				for {
					time.Sleep(time.Second * 3)
					conn, err = reconnect()
					if err == nil {
						break
					}
				}
			}
		}
	}
}

func getBroadcast()  {
	conn, err := net.Dial("udp", "192.168.2.255:8900")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println("连接成功")
	dt := DeviceInfo{
		// 设别序列号
		SerialNumber: "DTP00001",
		// 软件版本号
		FirmwareVersion: "V0.0.1",
		// 设备IP地址
		DeviceIp: "192.168.1.1",
		// 设备掩码
		Mask:"255.255.255.0",
		// 网关地址
		GatewayAddr:"192.168.1.1",
		// 设备端口号
		DevicePort:"8070",
		// 设备mac地址
		MacAddr:"5f-34-54-34-54",
	}
	dtJson, _ := json.Marshal(dt)
	data ,err := protocol.Pack(110, dtJson)
	size,err:=conn.Write(data)
	fmt.Println("write",size,err)
}

func main() {
	//tcpAddr, err := net.ResolveTCPAddr("tcp", "121.42.143.130:8989")
	//if err != nil {
	//	logger.Error("", err)
	//}
	//conn, err := net.DialTCP("tcp", nil, tcpAddr)
	//if err != nil {
	//	logger.Error("连接失败，", err)
	//}
	//defer conn.Close()
	//data ,err := protocol.Pack(888, []byte{'i', 'o', 't', 'h', 'u', 'b', 'a', 's', 'd', 'f', 'd'})
	//fmt.Println(data)
	//go write(conn, data)
	//go read(conn)
	go func() {
		for {
			getBroadcast()
			time.Sleep(time.Millisecond*100)
		}
	}()
	for {
		time.Sleep(time.Second * 10)
	}
}
