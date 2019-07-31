package service

import (
	"ai-platform/api/dao"
	"ai-platform/api/entity"
	"ai-platform/server/service"
)

type DeviceScanService interface {
	FindAll() ([]entity.DeviceScan, int, error)
}

type deviceScanServiceImpl struct {
	deviceManageInfoDao dao.DeviceManageInfoDao
}

func (r *deviceScanServiceImpl) FindAll() ([]entity.DeviceScan, int, error) {
	data := make([]entity.DeviceScan, 0)
	ret := service.GetOnlineDevice()
	idx := 0
	for _, val := range ret {
		item := entity.DeviceScan{
			// 设别序列号
			SerialNumber: val.SerialNumber,
			// 软件版本号
			FirmwareVersion: val.FirmwareVersion,
			// 设备IP地址
			Ip: val.Ip,
			// 设备掩码
			Mask: val.Mask,
			// 网关地址
			GatewayAddr: val.GatewayAddr,
			// 设备端口号
			Port: val.Port,
			// 设备mac地址
			MacAddr: val.MacAddr,
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
