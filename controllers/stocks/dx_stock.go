package stocks

import (
	"fmt"
	"github.com/shopspring/decimal"
	"reflect"
	. "stock/models"
	"stock/models/stocks_db"
	"stock/share/logging"
	"stock/share/util"
	"time"
)

type DxStock struct {
}

func NewDxStock() *DxStock {
	return &DxStock{}
}

// 短线选股保存到dx_stock 表中
func (this *DxStock) SaveDxstock() {
	defer func() {
		if err := recover(); err != nil {
			logging.Error("Panic Error=============:%v", err)
		}
	}()

	stocks_db.NewDxStockDb().DelDxStock()

	// 获取 stock_day_k 表最近23个交易日日期
	d := stocks_db.NewStock_Day_K().GetStockDayKDate()
	if len(d) < 15 {
		logging.Error("日K查询时间 Error：%v", len(d))
		return
	}
	ntime := time.Now().Format("2006-01-02")

	//{
	//	sql := `SELECT f12,f14,dayK5,dayK10,dayK20,dayK30 FROM  stock_day_k
	//	                      WHERE create_time='` + d[0]
	//	sql += `' AND dayK30 > dayK60 AND dayK20 > dayK30
	//			AND day5zdf < 8
	//			AND f8 >0.8 AND f8< 8 AND f10 >0.5 AND f10 <8
	//			AND f3 >0 AND f3 <3 AND f2 < 58 AND f9 <28 AND f9 >5
	//			AND f62 >5880000 AND f2 >f17 AND dayK10 > dayK20
	//			AND dayK5 < dayK10 `
	//
	//	sdkl := stocks_db.NewStock_Day_K().GetDxStockDayKList(sql)
	//
	//	if len(sdkl) > 0 {
	//		for _, v := range sdkl {
	//			//logging.Error("=========", v.F12, v.F14)
	//			//if NewStockDayk(nil).GetReturnIsBuy(v.F12) == false {
	//			//	continue
	//			//}
	//
	//			this.Save(v.F12, v.F14, ntime, v.DayK5, v.DayK10, v.DayK20, v.DayK30, 2)
	//		}
	//	}
	//}

	{
		sql := `SELECT f12,f14,dayK5,dayK10,dayK20,dayK30 FROM  stock_day_k
		                      WHERE create_time='` + d[0]
		sql += `' AND dayK20 > dayK60 AND dayK5 > dayK10
				AND dayK5 > dayK20 AND dayK10 > dayK20
				AND f2 > dayK5 AND f16 < dayK10
				AND f3 >0 AND f2 >f17 AND f2 < 18
				AND f12 IN(
					SELECT f12 FROM stock_day_k
					WHERE create_time='` + d[1]
		sql += `' AND f2 < f17 )`

		sdkl := stocks_db.NewStock_Day_K().GetDxStockDayKList(sql)

		if len(sdkl) > 0 {
			for _, v := range sdkl {
				//logging.Error("=========", v.F12, v.F14)
				//if NewStockDayk(nil).GetReturnIsBuy(v.F12) == false {
				//	continue
				//}

				this.Save(v.F12, v.F14, ntime, v.DayK5, v.DayK10, v.DayK20, v.DayK30, 2)
			}
		}
	}

	DxStockDb = nil
}

// 短线筛选保存
func (this *DxStock) Save(sc, sn, ntime string, dk5, dk10, dk20, dk30 float64, sta int) {

	i := stocks_db.NewDxStockDb()
	p := map[string]interface{}{
		"create_time": ntime,
		"stock_code":  sc,
		"stock_name":  sn,
		"status":      sta,
		"dayk5":       dk5,
		"dayk10":      dk10,
		"dayk20":      dk20,
		"dayk30":      dk30,
	}
	_, err1 := i.Insert(p)
	if err1 != nil {
		logging.Error("Insert Table dx_stock | %v", err1)
	}

}

// 需求个股分析监控 9：15 - 11：30   13：00-15：00  短线分析
func (this *DxStock) DxStockFx() {
	if len(DxStockDb) <= 0 {
		DxStockDb = stocks_db.NewDxStockDb().GetDxStockList()
	}
	name := ""
	defer func() {
		if err := recover(); err != nil {
			logging.Error("Panic Error=======:%v======:%v", name, err)
		}
	}()

	for k, v := range DxStockDb {

		//sc := controllers.NewUtilHttps(nil).GetUtilCode(v.StockCode)
		//if len(sc) <= 0 {
		//	continue
		//}

		i := NewStockDayk(nil).StockInfoSS(v.StockCode).StockDate
		if i == nil {
			continue
		}
		name = i.Gpmc
		zljlrv := 0.0
		if reflect.TypeOf(i.Zljlr).String() != "string" {
			zljlrv = i.Zljlr.(float64)
		}
		d1 := decimal.NewFromFloat(zljlrv)

		//d2 := "0"
		//if reflect.TypeOf(i.Jcd).String() != "string" {
		//	d2 = fmt.Sprintf("%v", i.Jcd.(float64))
		//}
		// 判断是否以入库
		if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
			continue
		}
		if reflect.TypeOf(i.Zxjg).Name() == "string" {
			continue
		}
		zxjgf := i.Zxjg.(float64)
		// 最新交易日判断 最低价最好是 回探 跌破五日 10日之上。然后 当前价 >= 5日的时候选出
		if i.Zdjg <= v.DayK5 && i.Zdjg >= v.DayK10 && zxjgf > v.DayK10 && zxjgf > i.Zdjg && d1.String() > "1888880" && i.Lb > 0.8 && i.Hsl > 1.28 {
			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, zxjgf, 5)
			DxStockDb = append(DxStockDb[:k], DxStockDb[k+1:]...)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))
		}
		// 开盘 最低价格 >= 五日K线 涨跌幅 不大于 3.8 量比 > 0.5  主力净流入 >0
		if i.Zdjg >= v.DayK5 && i.Zdf < 3.8 && i.Zdf > -0.8 && d1.String() > "1888880" && i.Lb > 1.28 && i.Hsl > 1.28 {
			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, zxjgf, 5)
			DxStockDb = append(DxStockDb[:k], DxStockDb[k+1:]...)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))
		}
	}
}
