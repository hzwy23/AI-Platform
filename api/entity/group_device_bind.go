package entity

type GroupDeviceBind struct {
	Id           int
	GroupId      int
	DeviceId     int
	CreateBy     string
	CreateDate   string
	UpdateBy     string
	UpdateDate   string
	DeleteStatus uint8
}
