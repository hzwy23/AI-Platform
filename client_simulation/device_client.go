package main

import (
	"fmt"
	"ai-platform/protocol"
	"ai-platform/panda/logger"
	"net"
	"time"
)

func reconnect() (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:8989")
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

		fmt.Println("读取到结果是：",msg, string(msg.MsgBody), err)
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

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:8989")
	if err != nil {
		logger.Error("", err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logger.Error("连接失败，", err)
	}
	defer conn.Close()
	data ,err := protocol.Pack(12, []byte{'a', 'b', 'c', 'd', 'e', 'f', 'a', 's', 'd', 'f', 'd'})
	fmt.Println(data)
	go write(conn, data)
	go read(conn)
	for {
		time.Sleep(time.Second * 10)
	}
}
