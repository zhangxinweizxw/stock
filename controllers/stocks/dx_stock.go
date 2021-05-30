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
	// 获取 stock_day_k 表最近23个交易日日期
	d := stocks_db.NewStock_Day_K().GetStockDayKDate()
	if len(d) < 15 {
		logging.Error("日K查询时间 Error：%v", len(d))
		return
	}
	ntime := time.Now().Format("2006-01-02")

	{
		sql := `SELECT f12,f14,dayK5,dayK10,dayK20,dayK30 FROM  stock_day_k
				WHERE create_time='` + d[0]
		sql += `' AND dayK5 >= dayK10  AND f3 > 0 AND f3 < 3.8  AND f7 < 5
				AND f2 > dayK20 AND f16 > dayK20 AND f16 < dayK5 AND day20zdf < 10 AND day20zdf > -5
				AND f12 NOT LIKE '688%' AND f14 NOT LIKE '*%'  AND f14 NOT LIKE 'ST%' 
				AND f12 IN (
				SELECT f12 FROM  stock_day_k
				WHERE create_time='` + d[1]
		sql += `' AND dayK5 >= dayK10  AND f3 > -1.8 AND f3 < 3.8 AND f7 < 5
				AND f2 > dayK30 AND f16 > dayK20 AND f16 < dayK5
				AND f12 IN (
				SELECT f12 FROM  stock_day_k
				WHERE create_time='` + d[2]
		sql += `' AND f3 > -1.8 AND f3 < 3.8  AND f7 < 5
				AND f2 > dayK30 AND f16 > dayK30 AND f16 < dayK5
				AND (dayK20 <= dayK30 OR dayK5 <= dayK10 OR dayK10 <= dayK20)
				AND f12 IN (
				SELECT f12 FROM  stock_day_k
				WHERE create_time='` + d[3]
		sql += `' AND f3 > -1.8 AND f3 < 3.8  AND f7 < 5
				AND f12 IN (
				SELECT f12 FROM  stock_day_k
				WHERE create_time='` + d[4]
		sql += `' AND f3 > -1.8 AND f3 < 3.8  AND f7 < 5
				AND f12 IN (
				SELECT f12 FROM  stock_day_k
				WHERE create_time='` + d[5]
		sql += `' AND f3 > -1.8 AND f3 < 3.8  AND f7 < 5
				AND f12 IN (
				SELECT f12 FROM  stock_day_k
				WHERE create_time='` + d[6]
		sql += `' AND f3 > -1.8 AND f3 < 3.8  AND f7 < 5
				AND f12 IN (
				SELECT f12 FROM  stock_day_k
				WHERE create_time='` + d[7] + "' ) ) ) ) ) ) )"
		//sql += `' AND f3 > -1.8 AND f3 < 3.8  AND f7 < 5
		//		AND f12 IN (
		//		SELECT f12 FROM  stock_day_k
		//		WHERE create_time='` + d[8]
		//sql += `' AND f3 > -1.8 AND f3 < 3.8  AND f7 < 5
		//		AND f12 IN (
		//		SELECT f12 FROM  stock_day_k
		//		WHERE create_time='` + d[9]
		sdkl := stocks_db.NewStock_Day_K().GetDxStockDayKList(sql)
		if len(sdkl) > 0 {
			for _, v := range sdkl {
				//logging.Error("=========", v.F12, v.F14)

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

		sc := ""
		switch v.StockCode[:3] {
		case "600", "601", "603", "605", "688", "689", "608":
			sc = fmt.Sprintf("SH%v", v.StockCode)
		case "300", "002", "000", "001", "003", "301":
			sc = fmt.Sprintf("SZ%v", v.StockCode)
		default:
			continue
		}

		i := NewStockDayk(nil).StockInfoSS(sc).StockDate

		name = i.Gpmc
		zljlrv := 0.0
		if reflect.TypeOf(i.Zljlr).String() != "string" {
			zljlrv = i.Zljlr.(float64)
		}
		d1 := decimal.NewFromFloat(zljlrv)
		//d2 := decimal.NewFromFloat(i.Jcd)
		//d3 := decimal.NewFromFloat(i.Jdd)
		d2 := "0"
		if reflect.TypeOf(i.Jcd).String() != "string" {
			d2 = fmt.Sprintf("%v", i.Jcd.(float64))
		}
		// 判断是否以入库
		if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
			continue
		}

		// 最新交易日判断 最低价最好是 回探 跌破五日 10日之上。然后 当前价 >= 5日的时候选出
		if i.Zdjg < v.DayK5 && (i.Zdjg >= v.DayK10 || i.Zdjg >= v.DayK20) && i.Zxjg >= v.DayK5 && d1.String() > "0" && d2 > "0" {
			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg, 5)
			DxStockDb = append(DxStockDb[:k], DxStockDb[k+1:]...)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))
		}
		// 开盘 最低价格 >= 五日K线 涨跌幅 不大于 3.8 量比 > 0.5  主力净流入 >0
		if i.Zdjg >= v.DayK5 && i.Zdf < 3.8 && i.Zdf > 0.18 && d1.String() > "5000000" && d2 > "0" && i.Lb > 0.5 {
			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg, 5)
			DxStockDb = append(DxStockDb[:k], DxStockDb[k+1:]...)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))
		}
	}
}
