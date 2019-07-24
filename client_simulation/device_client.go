package main

import (
	"ai-platform/panda/logger"
	"ai-platform/protocol"
	"fmt"
	"net"
	"time"
)

func reconnect() (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "118.31.46.174:8989")
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

		buf,_ := protocol.UnPack(rbuf[:rsize])
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
	conn, err := net.Dial("udp", "192.168.6.255:9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	data ,err := protocol.Pack(110, []byte{'a', 'b', 'c', 'd', 'e', 'f', 'a', 's', 'd', 'f', 'd'})
	conn.Write(data)
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "118.31.46.174:8989")
	if err != nil {
		logger.Error("", err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logger.Error("连接失败，", err)
	}
	defer conn.Close()
	data ,err := protocol.Pack(888, []byte{'a', 'b', 'c', 'd', 'e', 'f', 'a', 's', 'd', 'f', 'd'})
	fmt.Println(data)
	go write(conn, data)
	go read(conn)
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
