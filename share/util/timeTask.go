package util

import (
    "fmt"
    "stock/share/logging"
	"time"
)

type TimedTask struct {
}

func NewTimedTask() *TimedTask {
	return &TimedTask{}
}

//
//golang 定时器，启动的时候执行一次，
func (this *TimedTask) StartTimerStockDayK(f func()) {

	logging.Debug(fmt.Sprintf("%v  |======执行定时任务了|", time.Now()))
	go func() {
		for {
			f()
			now := time.Now()

			next := now.Add(time.Hour * 24)
			//next := now.Add(time.Minute * 1)
			next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()

}

// 定时任务精确到 小时分钟
func (this *TimedTask) StartTimer(f func(), hour int, minutes int) {
	logging.Debug(fmt.Sprintf("%v  |======执行定时任务了|", time.Now()))
	next := time.Now()

	go func() {
		logging.Debug("测试重复调用:111111111")
		for {
			now := time.Now()
			// 计算下一个零点
			next = time.Date(next.Year(), next.Month(), next.Day(), hour, minutes, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
			logging.Debug("测试重复调用:222222222")
			f()
			next = now.Add(time.Hour * 24)
		}
	}()
}

func GetTimeArr(start, end string) int64 {
	timeLayout := "2016-01-02"
	loc, _ := time.LoadLocation("Local")
	// 转成时间戳
	startUnix, _ := time.ParseInLocation(timeLayout, start, loc)
	endUnix, _ := time.ParseInLocation(timeLayout, end, loc)
	startTime := startUnix.Unix()
	endTime := endUnix.Unix()
	// 求相差天数
	date := (endTime - startTime) / 86400
	return date
}
