package filter

import (
	"ai-platform/api/auth/service"
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
	if this.staticFile(r.URL.Path) || service.Identify(w, r) {
		nw := route.NewResponse(w)
		next(nw, r)
	}
}
