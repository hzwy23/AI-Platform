package entity

// 设备分组信息表
type DeviceGroupInfo struct {
	GroupId      int
	GroupName    string
	CreateBy     string
	CreateDate   string
	UpdateBy     string
	UpdateDate   string
	DeleteStatus uint8
}
