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
	//         <Option value="3">爆闪</Option>
	DeviceAttribute string `json:"client_Mode"`
}

var LightValue = map[string]string{
	"1":  "2",
	"2":  "4",
	"3":  "6",
	"4":  "8",
	"5":  "10",
	"6":  "12",
	"7":  "14",
	"8":  "16",
	"9":  "18",
	"10": "20",
	"11": "22",
	"12": "24",
	"13": "26",
	"14": "28",
	"15": "30",
	"16": "32",
	"17": "34",
	"18": "36",
	"19": "38",
	"20": "40",
}

var LightValueReserve = map[string]string{
	"2":  "1",
	"4":  "2",
	"6":  "3",
	"8":  "4",
	"10": "5",
	"12": "6",
	"14": "7",
	"16": "8",
	"18": "9",
	"20": "10",
	"22": "11",
	"24": "12",
	"26": "13",
	"28": "14",
	"30": "15",
	"32": "16",
	"34": "17",
	"36": "18",
	"38": "19",
	"40": "20",
}

var ThresholdValue = map[string]string{
	"1": "3998",
	"2": "3552",
	"3": "2698",
	"4": "2175",
	"5": "1747",
	"6": "1395",
}

var ThresholdValueReserve = map[string]string{
	"3998": "1",
	"3552": "2",
	"2698": "3",
	"2175": "4",
	"1747": "5",
	"1395": "6",
}
