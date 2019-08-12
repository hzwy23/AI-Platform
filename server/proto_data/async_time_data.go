package proto_data

type AsyncTimeData struct {
	SerialNumber string `json:"client_CPUID"`
	TimeZone     string `json:"TimeZone"`
	Year         string `json:"Year"`
	Month        string `json:"Month"`
	Day          string `json:"Day"`
	Hour         string `json:"Hour"`
	Minute       string `json:"Minute"`
	Second       string `json:"Second"`
	Week         string `json:"Week"`
}
