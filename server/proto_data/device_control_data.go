package proto_data

type DeviceControlData struct {
	SerialNumber  string `json:"client_CPUID"`
	LightMode     string `json:"client_AutoFunction"`
	Timer         []DeviceTimerData `json:"timer"`
}

type DeviceTimerData struct {
	AutoStartTime string `json:"AutoTimeStart"`
	AutoEndTime   string `json:"AutoTimeStop"`
}