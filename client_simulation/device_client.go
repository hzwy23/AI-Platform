package main

import (
 "ai-platform/panda/logger"
 "ai-platform/protocol"
 "fmt"
 "net"
 "time"
)

type DeviceInfo struct {
 // 设别序列号
 SerialNumber string `json:"client_CPUID"`
 // 软件版本号
 FirmwareVersion string `json:"client_FrameworkVersion"`
 // 设备IP地址
 DeviceIp string `json:"client_IP"`
 // 设备掩码
 Mask string `json:"client_MASK"`
 // 网关地址
 GatewayAddr string `json:"client_GATEWAY"`
 // 设备端口号
 DevicePort string `json:"client_PORT"`
 // 设备mac地址
 MacAddr string `json:"client_MAC"`
}

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
 //dt := DeviceInfo{
 // // 设别序列号
 // SerialNumber: "DTP00001",
 // // 软件版本号
 // FirmwareVersion: "V0.0.1",
 // // 设备IP地址
 // DeviceIp: "192.168.1.1",
 // // 设备掩码
 // Mask:"255.255.255.0",
 // // 网关地址
 // GatewayAddr:"192.168.1.1",
 // // 设备端口号
 // DevicePort:"8070",
 // // 设备mac地址
 // MacAddr:"5f-34-54-34-54",
 //}
 //_, _ := json.Marshal(dt)

 data ,err := protocol.Pack(0x7fff, []byte(`{"client_IP":"192.168.2.100","client_PORT":"8989","client_CPUID":"1111111111111111","client_FrameworkVersion":"1.2.3","client_GATEWAY":"192.168.2.1","client_MAC":"01-df-de-sa-df","client_Mode":"Auto","client_MASK":"255.255.255.0","client_Temp":"47.84"}`))
 conn.Write(data)
 //for i:=0; i < len(data); i++ {
 //  conn.Write(data[i:i+1])
 //  fmt.Println(data[i:i+1])
 //  time.Sleep(time.Millisecond * 1)
 //}
 conn.Write(data)
 fmt.Println("write",err)
}

func main() {
 //tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:8989")
 //if err != nil {
 //logger.Error("", err)
 //}
 //conn, err := net.DialTCP("tcp", nil, tcpAddr)
 //if err != nil {
 //logger.Error("连接失败，", err)
 //}
 //defer conn.Close()
 //data ,err := protocol.Pack(0x7fff, []byte{'i', 'o', 't', 'h', 'u', 'b', 'a', 's', 'd', 'f', 'd'})
 //fmt.Println(data)
 //go write(conn, data)
 //go read(conn)
 go func() {
  for {
   getBroadcast()
   time.Sleep(time.Second*5)
  }
 }()
 for {
  time.Sleep(time.Second * 10)
 }
}