package filter

import (
	service2 "ai-platform/api/auth/service"
	"ai-platform/panda/route"
	"net/http"
	"strings"
)

type AuthFilter struct {
}

func (this *AuthFilter) staticFile(url string) bool {
	if !strings.HasPrefix(url, "/api") {
		return true
	}
	return url == "/favicon.ico"
}

func (this *AuthFilter) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if this.staticFile(r.URL.Path) || service2.Identify(w, r) {
		nw := route.NewResponse(w)
		next(nw, r)
	}
}

func init() {
	// 设置白名单，免认证请求
	service2.AddConnUrl("/")
	service2.AddConnUrl("/v1/auth/login")

	/// 设置白名单，免授权请求
	service2.AddAuthUrl("/v1/auth/logout")
	service2.AddAuthUrl("/v1/auth/theme/update")
	service2.AddAuthUrl("/v1/auth/user/query")
	service2.AddAuthUrl("/v1/auth/HomePage")
	service2.AddAuthUrl("/v1/auth/main/menu")
	service2.AddAuthUrl("/v1/auth/index/entry")
	service2.AddAuthUrl("/v1/auth/privilege/user/domain")
	service2.AddAuthUrl("/v1/auth/menu/all/except/button")
}
