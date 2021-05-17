package main

import (
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"stock/config"
	"stock/controllers/stocks"
	. "stock/models"
	"stock/models/stocks_db"
	"stock/routes"
	. "stock/share/app"
	"stock/share/logging"
	"stock/share/util"
	"time"
)

func main() {

	cfg := config.Default(APP_PID)

	// 项目初始化
	a := NewApp(APP_NAME, APP_VERSION)
	a.PidName = APP_PID
	a.WSPort = cfg.Serve.Port
	a.LogPort = cfg.Log.Port
	a.LogAddr = cfg.Log.Addr
	a.LogOn = cfg.Log.On
	a.SessionOn = cfg.Session.On
	a.SessionProviderName = cfg.Session.ProviderName
	a.SessionConfig = cfg.Session.Config
	a.DisableGzip = true
	a.Cors = cfg.Cors.AllowOrigin

	r := a.Init()

	// 路由注册
	routes.Register(r)
	//go util.NewDdRobot().DdRobotPush("test 数据 更新异常预警信息")

	//stocks.NewStockDayk(cfg).GetStockDayK()
	//f10.NewFinancialReports().SaveFinaRepo()
	//stocks.NewStockDayk(nil).GetDayK("002555")
	//logging.Error("================:", stocks_db.NewTransactionHistory().GetTranHist("九安医疗"))
	//controllers.NewUtilHttps(nil).GetXqPd()

	go func() {
		for { // 1
			//判断当天是否是交易日
			b := util.NewStockUtil(cfg).GetSjsMonthList()

			if b {
				c1 := stocks_db.NewStock_Day_K().GetIsZx()

				if c1 == time.Now().Format("2006-01-02") {
					time.Sleep(8 * time.Hour)
					continue
				}

				//雪球筛选
				if time.Now().Hour() >= 17 {

					//每天下午跑日K数据
					stocks.NewStockDayk(cfg).GetStockDayK()

					stocks.NewStockDayk(cfg).GetXueqiu()

					stocks.NewAvsHStock(cfg).SaveAvsHStock()

					stocks.NewZjlxStock().ZjlxStockSave()

					stocks.NewQgqpStock().QgqpStockSave()

					stocks.NewDxStock().SaveDxstock()
				}
				time.Sleep(1 * time.Hour)
			}
		}
	}()

	go func() {
		for { // 1
			//判断当天是否是交易日
			b := util.NewStockUtil(cfg).GetSjsMonthList()
			if b {

				time1 := fmt.Sprintf("%v 09:29", time.Now().Format("2006-01-02"))
				time2 := fmt.Sprintf("%v 14:59", time.Now().Format("2006-01-02"))
				time3 := fmt.Sprintf("%v 08:01", time.Now().AddDate(0, 0, 1).Format("2006-01-02"))
				//先把时间字符串格式化成相同的时间类型
				t1, err := time.Parse("2006-01-02 15:04", time.Now().Format("2006-01-02 15:04"))
				t2, err := time.Parse("2006-01-02 15:04", time1)
				t3, err := time.Parse("2006-01-02 15:04", time2)
				t4, err := time.Parse("2006-01-02 15:04", time3)

				//logging.Error("======：", t1, "=====:", t3, "======", t4, "=====", err)
				if err == nil && t1.After(t3) && t1.Before(t4) {
					//处理逻辑
					logging.Debug("============:15点以后休眠1小时")
					time.Sleep(1 * time.Hour)
					continue
				}

				if err == nil && t1.After(t2) && t1.Before(t3) {
					// 雪球筛选处理逻辑
					go stocks.NewStockDayk(nil).XQStockFx()
					go stocks.NewAvsHStock(nil).AvsHStockFx()
					//go stocks.NewZjlxStock().ZjlxtockFx()
					go stocks.NewQgqpStock().QgqpStockFx()
					go stocks.NewDxStock().DxStockFx()
					go stocks.NewZjlxStock().PkydStockFx()

					go stocks.NewZjlxStock().ZjlxStockSellFx()

				}

				time.Sleep(30 * time.Second)
			}
		}
	}()

	// stock 逻辑处理-------------------

	// 监控性能
	ginpprof.Wrapper(r)

	logging.Error("%s", r.Run(cfg.Serve.Port))

}
