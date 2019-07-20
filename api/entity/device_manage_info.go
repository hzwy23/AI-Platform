package entity

type DeviceManageInfo struct {
	DeviceId             int
	SerialNumber         string
	DeviceName           string
	DhcpFlag             uint8
	DeviceIp             string
	DevicePort           string
	DeviceStatus         uint8
	DeviceAttribute      uint8
	DevicePower          int
	DeviceLightThreshold uint8
	DeviceBrightness     uint8
	DeviceTemperature    uint8
	AutoStartTime        string
	AutoEndTime          string
	LightMode            uint8
	MacAddress           string
	FirmwareVersion      string
	Longitude            string
	Latitude             string
	Address              string
	Mask                 string
	Gateway              string
	Pin                  string
	CreateBy             string
	CreateDate           string
	UpdateBy             string
	UpdateData           string
	DeleteStatus         uint8
	PowerTotal           int
	StrobeCount          int
}
