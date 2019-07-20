package platform

import "net"

// 设备管理
type Device struct {
	// 设备号
	DeviceId string
	//
	Conn net.Conn

	// 设备IP地址
	IpAddr string

	// 设备分组信息
	GroupId string
}
