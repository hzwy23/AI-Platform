package entity

type DeviceInstallInfo struct {
	Id            int
	SerialNumber  string
	DeviceAddress string
	Lat           string
	Lon           string
	CreateDate    string
	CreateBy      string
	UpdateDate    string
	UpdateBy      string
	DeleteStatus  int
}
