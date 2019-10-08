package cron

import (
	"ai-platform/dbobj"
	"fmt"
	"github.com/robfig/cron"
	"time"
)

func init()  {
	c := cron.New(cron.WithSeconds())
	// 每天13：30定时清楚7天前的设备与平台交互日志
	c.AddFunc("0 30 13 * * *", func() {
		fmt.Println("定时任务执行清除历史日志信息：", time.Now().Format("2006-01-02 15:04:05"))
		dbobj.Exec("delete from plat_device_logger where date_format(handle_time,'%Y-%m-%d %H:%i:%s') < DATE_SUB(now(), INTERVAL 7 day)")
		fmt.Println("清除设备与平台7天前交互日志完成")
	})

	c.AddFunc("0 30 14 * * *", func() {
		fmt.Println("定时任务执行清除历史用户操作日志",time.Now().Format("2006-01-02 15:04:05"))
		dbobj.Exec("delete from plat_user_logger where date_format(handle_time,'%Y-%m-%d %H:%i:%s') < DATE_SUB(now(), INTERVAL 1 year)")
		fmt.Println("清除用户1年前操作记录成功")
	})

	c.Start()
}
