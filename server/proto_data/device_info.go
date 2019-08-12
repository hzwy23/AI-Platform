package proto_data

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
