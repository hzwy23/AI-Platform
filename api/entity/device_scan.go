package entity

type DeviceScan struct {
	// 设别序列号
	SerialNumber string
	// 软件版本号
	FirmwareVersion string
	// 设备IP地址
	Ip string
	// 设备掩码
	Mask string
	// 网关地址
	GatewayAddr string
	// 设备端口号
	Port string
	// 设备mac地址
	MacAddr string
	// 是否添加
	IsAdded bool
}
