package service

var RouteService = &RouteServiceImpl{}

type RouteServiceImpl struct {
}

func (this *RouteServiceImpl) CheckUrlAuth(userId string, url, method string) bool {
	return true
	//return returnthis.role.CheckUrlAuth(userId, url, method)
}

func (this *RouteServiceImpl) CheckResIDAuth(userId string, resId string) bool {
	return true
	//return thiss.role.CheckResIDAuth(userId, resId)
}
