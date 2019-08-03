package filter

import (
	"ai-platform/api/dao"
	"ai-platform/api/entity"
	"ai-platform/panda"
	"ai-platform/panda/hret"
	"ai-platform/panda/jwt"
	"ai-platform/panda/route"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type LoggerFilter struct {
}

func (log *LoggerFilter) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
	go writeHandleLogs(w, r)
}

func init() {
	// 开启操作日志监听
	go logSync()
}

var logger = dao.NewPlatUserLoggerDao()
var logBuf = make(chan entity.PlatUserLogger, 40960)

func writeHandleLogs(w http.ResponseWriter, r *http.Request) {
	defer hret.RecvPanic()

	if nw, ok := w.(*route.Response); ok {
		var one entity.PlatUserLogger
		one.RetCode = strconv.Itoa(nw.Status)
		one.ReqUrl = r.URL.Path
		one.ReqParam = formencode(r.Form)
		one.ReqMethod = r.Method
		one.HandleTime = panda.CurTime()

		claim, err := jwt.ParseHttp(r)
		if err != nil {
			one.UserId = "anonymous"
		} else {
			one.UserId = claim.UserId
		}
		logBuf <- one
	}
}

func formencode(form url.Values) string {
	rst := make(map[string]string)
	for key, val := range form {
		if key == "_" {
			continue
		}
		rst[key] = val[0]
	}

	str, _ := json.Marshal(rst)
	return string(str)
}

func savelogs(data []entity.PlatUserLogger) {
	for _, item := range data {
		logger.Insert(item)
	}
}

func logSync() {
	var buf []entity.PlatUserLogger
	for {
		select {
		case <-time.After(time.Second * 10):
			// sync handle logs to database per 5 second.
			if len(buf) == 0 {
				continue
			}
			go savelogs(buf)
			buf = make([]entity.PlatUserLogger, 0)
		case val, ok := <-logBuf:
			if ok {
				buf = append(buf, val)
				if len(buf) > 1000 {
					go savelogs(buf)
					buf = make([]entity.PlatUserLogger, 0)
				}
			}
		}
	}
}
