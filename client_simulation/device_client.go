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

func write(conn *net.TCPConn, message *protocol.Message) {
	sheader := protocol.Header{
		MsgLength:   100,
		MsgSN:       12122,
		MsgID:       12,
		MsgCenterID: 34,
		VersionFlag: [3]byte{0x01, 0x02, 0x03},
		EncryptFlag: 0x01,
		EncryptKey:  45,
	}

	body := protocol.Encrypt(protocol.KEY, []byte{0x5a, 0x01, 'x', 0x5e, 0x01, 'x', 0x5e, 0x02, 'x', 0x5a, 0x02})
	fmt.Println(body)
	crc, _ := protocol.CRC16CCITT(body)

	smessage := &protocol.Message{
		HeaderFlag: protocol.HEADER_FLAG,
		MsgHeader:  sheader,
		MsgBody:    body,
		CrcCode:    crc,
		FooterFlag: protocol.FOOTER_FLAG,
	}

	buf := protocol.ConvertToByte(message)
	buf = append(buf, protocol.ConvertToByte(smessage)...)
	idx := 0
	for {
		_, err := conn.Write(buf[idx : idx+1])
		fmt.Println("发送数据：", buf[idx:idx+1])
		idx += 1
		if len(buf) == idx {
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
		fmt.Println("读取到结果是：", rbuf[:rsize], err)
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
	header := protocol.Header{
		MsgLength:   100,
		MsgSN:       12122,
		MsgID:       12,
		MsgCenterID: 34,
		VersionFlag: [3]byte{0x01, 0x02, 0x03},
		EncryptFlag: 0x01,
		EncryptKey:  45,
	}
	body := protocol.Encrypt(protocol.KEY, []byte{'a', 'b', 'c', 'd', 'e', 'f', 'a', 's', 'd', 'f', 'd'})
	fmt.Println(body)
	crc, _ := protocol.CRC16CCITT(body)

	fmt.Println(crc, 123)

	message := &protocol.Message{
		HeaderFlag: protocol.HEADER_FLAG,
		MsgHeader:  header,
		MsgBody:    body,
		CrcCode:    crc,
		FooterFlag: protocol.FOOTER_FLAG,
	}

	go write(conn, message)
	go read(conn)
	for {
		time.Sleep(time.Second * 10)
	}
}
