package service

import (
	"ai-platform/panda/logger"
	"regexp"
)

var RouteService = &RouteServiceImpl{}

var permit = map[string]map[string][]string{
	"admin": {
		"GET": {
			"/api/device/install",
			"/api/device/logger",
			"/api/device/manage",
			"/api/device/manage/ungroup",
			"/api/scan/device",
			"/api/device/offline",
			"/api/device/temperature",
			"/api/device/lamp/exception",
			"/api/device/group",
			"/api/homepage/statistics",
			"/api/platform/logger",
			`/api/device/plus/[\w]+`,
			`/api/device/minus/[\w]+`,
			"/api/global/config",
			"/api/account/profiles",
		},
		"POST": {
			"/api/device/install",
			"/api/device/manage",
			"/api/device/manage/group",
			"/api/device/group",
			"/api/login",
			`/api/device/remote/control/[\w]+`,
		},
		"PUT": {
			"/api/device/install",
			"/api/device/group/change",
			"/api/device/event",
			"/api/device/group",
			`/api/device/manage/network/[\w]+`,
			`/api/device/manage/attribute/[\w]+`,
			"/api/global/config",
			"/api/account/password",
			"/api/account/profiles",
			"/api/device/global/setting",
		},
		"DELETE": {
			`/api/device/install/[\w]+`,
			`/api/device/manage/[\w]+`,
			`/api/device/bind/[\w]+`,
			`/api/device/group/[\w]+`,
		},
	},
	"operation": {
		"GET": {
			"/api/device/install",
			"/api/device/logger",
			"/api/device/manage",
			"/api/device/manage/ungroup",
			"/api/scan/device",
			"/api/device/offline",
			"/api/device/temperature",
			"/api/device/lamp/exception",
			"/api/device/group",
			"/api/homepage/statistics",
			"/api/platform/logger",
			`/api/device/plus/[\w]+`,
			`/api/device/minus/[\w]+`,
			"/api/global/config",
			"/api/account/profiles",
		},
		"POST": {
			"/api/device/install",
			"/api/device/manage",
			"/api/device/manage/group",
			"/api/device/group",
			"/api/login",
			`/api/device/remote/control/[\w]+`,
		},
		"PUT": {
			"/api/device/install",
			"/api/device/group/change",
			"/api/device/event",
			"/api/device/group",
			`/api/device/manage/network/[\w]+`,
			`/api/device/manage/attribute/[\w]+`,
			"/api/global/config",
			"/api/account/password",
			"/api/account/profiles",
		},
		"DELETE": {
			`/api/device/install/[\w]+`,
			`/api/device/manage/[\w]+`,
			`/api/device/bind/[\w]+`,
			`/api/device/group/[\w]+`,
		},
	},
	"operation1": {
		"GET": {
			"/api/device/install",
			"/api/device/logger",
			"/api/device/manage",
			"/api/device/manage/ungroup",
			"/api/scan/device",
			"/api/device/offline",
			"/api/device/temperature",
			"/api/device/lamp/exception",
			"/api/device/group",
			"/api/homepage/statistics",
			"/api/platform/logger",
			`/api/device/plus/[\w]+`,
			`/api/device/minus/[\w]+`,
			"/api/global/config",
			"/api/account/profiles",
		},
		"POST": {
			"/api/device/install",
			"/api/device/manage",
			"/api/device/manage/group",
			"/api/device/group",
			"/api/login",
			`/api/device/remote/control/[\w]+`,
		},
		"PUT": {
			"/api/device/install",
			"/api/device/group/change",
			"/api/device/event",
			"/api/device/group",
			`/api/device/manage/network/[\w]+`,
			`/api/device/manage/attribute/[\w]+`,
			"/api/global/config",
			"/api/account/password",
			"/api/account/profiles",
		},
		"DELETE": {
			`/api/device/install/[\w]+`,
			`/api/device/manage/[\w]+`,
			`/api/device/bind/[\w]+`,
			`/api/device/group/[\w]+`,
		},
	},
	"operation2": {
		"GET": {
			"/api/device/install",
			"/api/device/logger",
			"/api/device/manage",
			"/api/device/manage/ungroup",
			"/api/scan/device",
			"/api/device/offline",
			"/api/device/temperature",
			"/api/device/lamp/exception",
			"/api/device/group",
			"/api/homepage/statistics",
			"/api/platform/logger",
			`/api/device/plus/[\w]+`,
			`/api/device/minus/[\w]+`,
			"/api/global/config",
			"/api/account/profiles",
		},
		"POST": {
			"/api/device/install",
			"/api/device/manage",
			"/api/device/manage/group",
			"/api/device/group",
			"/api/login",
			`/api/device/remote/control/[\w]+`,
		},
		"PUT": {
			"/api/device/install",
			"/api/device/group/change",
			"/api/device/event",
			"/api/device/group",
			`/api/device/manage/network/[\w]+`,
			`/api/device/manage/attribute/[\w]+`,
			"/api/global/config",
			"/api/account/password",
			"/api/account/profiles",
		},
		"DELETE": {
			`/api/device/install/[\w]+`,
			`/api/device/manage/[\w]+`,
			`/api/device/bind/[\w]+`,
			`/api/device/group/[\w]+`,
		},
	},
	"inspector": {
		"GET": {
			"/api/device/install",
			"/api/device/logger",
			"/api/device/manage",
			"/api/device/manage/ungroup",
			"/api/scan/device",
			"/api/device/offline",
			"/api/device/temperature",
			"/api/device/lamp/exception",
			"/api/device/group",
			"/api/homepage/statistics",
			"/api/platform/logger",
			"/api/global/config",
			"/api/account/profiles",
		},
		"POST": {
			"/api/login",
		},
		"PUT": {
			"/api/account/password",
			"/api/account/profiles",
		},
		"DELETE": {},
	},
	"inspector1": {
		"GET": {
			"/api/device/install",
			"/api/device/logger",
			"/api/device/manage",
			"/api/device/manage/ungroup",
			"/api/scan/device",
			"/api/device/offline",
			"/api/device/temperature",
			"/api/device/lamp/exception",
			"/api/device/group",
			"/api/homepage/statistics",
			"/api/platform/logger",
			"/api/global/config",
			"/api/account/profiles",
		},
		"POST": {
			"/api/login",
		},
		"PUT": {
			"/api/account/password",
			"/api/account/profiles",
		},
		"DELETE": {},
	},
	"inspector2": {
		"GET": {
			"/api/device/install",
			"/api/device/logger",
			"/api/device/manage",
			"/api/device/manage/ungroup",
			"/api/scan/device",
			"/api/device/offline",
			"/api/device/temperature",
			"/api/device/lamp/exception",
			"/api/device/group",
			"/api/homepage/statistics",
			"/api/platform/logger",
			"/api/global/config",
			"/api/account/profiles",
		},
		"POST": {
			"/api/login",
		},
		"PUT": {
			"/api/account/password",
			"/api/account/profiles",
		},
		"DELETE": {},
	},
}

type RouteServiceImpl struct {
}

func (this *RouteServiceImpl) CheckUrlAuth(userId string, url, method string) bool {
	if permits, ok := permit[userId]; ok {
		if urls, yes := permits[method]; yes && len(urls) > 0 {
			for _, val := range urls {
				if m, _ := regexp.MatchString(val, url); m {
					return true
				} else {
					continue
				}
			}
		}
	}
	logger.Info("请求信息是：", userId, url, method)
	return false
}
