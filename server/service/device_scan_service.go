package service

import (
	"ai-platform/api/dao"
	"ai-platform/api/entity"
	"ai-platform/dbobj"
	"ai-platform/server/listen"
)

type DeviceScanService interface {
	FindAll() ([]entity.DeviceScan, int, error)
}

type deviceScanServiceImpl struct {
	deviceManageInfoDao dao.DeviceManageInfoDao
}

func (r *deviceScanServiceImpl) FindAll() ([]entity.DeviceScan, int, error) {
	data := make([]entity.DeviceScan, 0)
	ret,_ := listen.GetOnlineDevice()
	idx := 0
	for _, val := range ret {

		address := ""
		dbobj.QueryForObject("select device_address from device_install_info where delete_status = 0 and serial_number = ?", dbobj.PackArgs(val.SerialNumber), &address)

		item := entity.DeviceScan{
			// 设别序列号
			SerialNumber: val.SerialNumber,
			// 软件版本号
			FirmwareVersion: val.FirmwareVersion,
			// 设备IP地址
			DeviceIp: val.DeviceIp,
			// 设备掩码
			Mask: val.Mask,
			// 网关地址
			GatewayAddr: val.GatewayAddr,
			// 设备端口号
			DevicePort: val.DevicePort,
			// 设备mac地址
			MacAddr: val.MacAddr,
			// 安装位置
			DeviceAddress: address,
		}
		element, err := r.deviceManageInfoDao.FindBySerialNumber(val.SerialNumber)

		if err == nil && element.SerialNumber == val.SerialNumber {
			item.IsAdded = true
			idx += 1
		}
		data = append(data, item)
	}
	return data, idx, nil
}

func NewDeviceScanService() DeviceScanService {
	r := &deviceScanServiceImpl{
		deviceManageInfoDao: dao.NewDeviceManageInfoDao(),
	}
	return r
}
