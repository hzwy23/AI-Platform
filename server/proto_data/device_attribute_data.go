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

var LightValue = map[string]string{
	"1":  "1",
	"2":  "2",
	"3":  "3",
	"4":  "5",
	"5":  "8",
	"6":  "15",
	"7":  "21",
	"8":  "36",
	"9":  "58",
	"10": "66",
	"11": "79",
	"12": "93",
	"13": "110",
	"14": "120",
	"15": "128",
	"16": "139",
	"17": "169",
	"18": "181",
	"19": "190",
	"20": "200",
}

var LightValueReserve = map[string]string{
	"1":   "1",
	"2":   "2",
	"3":   "3",
	"5":   "4",
	"8":   "5",
	"15":  "6",
	"21":  "7",
	"36":  "8",
	"58":  "9",
	"66":  "10",
	"79":  "11",
	"93":  "12",
	"110": "13",
	"120": "14",
	"128": "15",
	"139": "16",
	"169": "17",
	"181": "18",
	"190": "19",
	"200": "20",
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
