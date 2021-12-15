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
	//f10.NewFinancialReports().SaveFinaRepo(nil)
	//stocks.NewStockDayk(nil).GetDayK("002555")
	//logging.Error("================:", stocks_db.NewTransactionHistory().GetTranHist("九安医疗"))
	//controllers.NewUtilHttps(nil).GetXqPd()
	//stocks_db.NewStock_Day_K().GetSStockInfo("000225")
	//stocks.NewStockDayk(nil).GetReturnIsBuy()

	var status = 1
	var status01 = 1

	go func() {
		for { // 1
			//判断当天是否是交易日
			b := util.NewStockUtil().GetSjsMonthList()

			if b {
				c1 := stocks_db.NewStock_Day_K().GetIsZx()
				if c1 != 0 {
					time.Sleep(8 * time.Hour)
					continue
				}
				//雪球筛选
				if time.Now().Hour() >= 18 {
					//每天下午跑日K数据
					stocks.NewStockDayk(cfg).GetStockDayK()

					time.Sleep(15 * time.Minute)
					stocks.NewStockDayk(cfg).GetXueqiu()
					stocks.NewStockDayk(cfg).SaveXueqiuFx()

					//stocks.NewZjlxStock().ZjlxStockSave()

					stocks.NewQgqpStock().QgqpStockSave()

					stocks.NewDxStock().SaveDxstock()
					stocks_db.NewZtStockDB().DelZtStock()
					status, status01 = 1, 1
				}
				time.Sleep(1 * time.Hour)
			} else {
				logging.Info("非交易时间")
				time.Sleep(8 * time.Hour)
			}
		}
	}()

	go func() {
		for { // 1
			//判断当天是否是交易日
			b := util.NewStockUtil().GetSjsMonthList()
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
					if time.Now().Hour() == 9 && time.Now().Minute() == 33 && status == 1 {
						stocks.NewZtStock().GetZTStock()
						status = 0
					}

					if time.Now().Hour() == 14 && time.Now().Minute() == 29 && status01 == 1 {
						stocks.NewZtStock().GetZTStock01()
						status01 = 0
					}

					zt1, _ := time.Parse("2006-01-02 15:04", fmt.Sprintf("%v 09:34", time.Now().Format("2006-01-02")))
					zt2, _ := time.Parse("2006-01-02 15:04", fmt.Sprintf("%v 11:28", time.Now().Format("2006-01-02")))
					zt3, _ := time.Parse("2006-01-02 15:04", fmt.Sprintf("%v 13:28", time.Now().Format("2006-01-02")))
					zt4, _ := time.Parse("2006-01-02 15:04", fmt.Sprintf("%v 14:56", time.Now().Format("2006-01-02")))

					if (t1.After(zt1) && t1.Before(zt2)) || (t1.After(zt3) && t1.Before(zt4)) {
						stocks.NewZtStock().ZtStockFx()
					}

					go stocks.NewXqFxStock().XqFxTs()
					go stocks.NewStockDayk(nil).XQStockFx()
					//go stocks.NewZjlxStock().ZjlxtockFx()
					go stocks.NewQgqpStock().QgqpStockFx()
					go stocks.NewDxStock().DxStockFx()
					go stocks.NewZjlxStock().PkydStockFx()

					go stocks.NewZjlxStock().ZjlxStockSellFx()

				}

				time.Sleep(30 * time.Second)
			} else {
				logging.Info("非交易时间")
				time.Sleep(8 * time.Hour)
			}
		}
	}()

	// stock 逻辑处理-------------------

	// 监控性能
	ginpprof.Wrapper(r)

	logging.Error("%s", r.Run(cfg.Serve.Port))

}
