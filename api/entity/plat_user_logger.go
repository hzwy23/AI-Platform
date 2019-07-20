package entity

type PlatUserLogger struct {
	Id         int
	UserId     string
	HandleTime string
	ReqMethod  string
	ReqUrl     string
	ReqParam   string
	RetMsg     string
	RetCode    string
}
