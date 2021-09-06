package model

import "github.com/robfig/cron"

func Croninit() {
	go func() {
		crontab := cron.New()
		crontab.AddFunc("0 */3 * * * *", CheckAgentStatus) // 每隔3分钟 定时执行 CheckAgentStatus 函数
		crontab.AddFunc("0 */4 * * * * ", WritePrometheus)
		crontab.AddFunc("0 */50 * * * * ", ScanHardWareInfo)
		//crontab.AddFunc("0 */4 * * * * ", GenerateAnsibleHosts)
		crontab.Start()
	}()
}
