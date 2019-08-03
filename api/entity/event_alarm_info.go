package entity

type EventAlarmInfo struct {
	Id                int
	EventTypeCd       int
	OccurrenceTime    string
	SerialNumber      string
	DeviceName        string
	DeviceIp          string
	DeviceAttribute   uint8
	DeviceBrightness  uint8
	DeviceTemperature string
	HandleStatus      uint8
	DeleteStatus      uint8
}
