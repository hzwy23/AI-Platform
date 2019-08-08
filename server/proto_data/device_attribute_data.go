package proto_data

type DeviceAttribute struct {
	// 序列号
	SerialNumber string `json:"client_CPUID"`
	// 功率
	DevicePower string `json:"client_Power"`
	// 温度
	DeviceTemperature string `json:"client_Temp"`
	// 光敏阀值
	DeviceLightThreshold string `json:"client_CDSThreshold"`
	// 设备亮度
	DeviceBrightness string `json:"client_LightLevel"`
	// 总功耗
	PowerTotal string `json:"client_Consumption"`
	// 爆闪次数
	StrobeCount string `json:"client_FlashCount"`
	// 光照强度
	IntensityLight string `json:"client_CDS"`
	// 设备属性  <Option value="1">常亮</Option>
	//         <Option value="2">频闪</Option>
	//         <Option value="3">爆灯</Option>
	DeviceAttribute string `json:"client_Mode"`
}
