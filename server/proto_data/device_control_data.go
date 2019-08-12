package proto_data

type DeviceControlData struct {
	SerialNumber  string `json:"client_CPUID"`
	LightMode     string `json:"client_AutoFunction"`
	AutoStartTime string `json:"AutoTimeStart"`
	AutoEndTime   string `json:"AutoTimeStop"`
}
