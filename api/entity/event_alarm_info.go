package entity

type EventAlarmInfo struct {
	Id                int
	EventTypeCd       int
	OccurrenceTime    string
	DeviceId          string
	DeviceName        string
	DeviceIp          string
	DeviceAttribute   uint8
	DeviceBrightness  uint8
	DeviceTemperature uint8
	HandleStatus      uint8
	DeleteStatus      uint8
}
