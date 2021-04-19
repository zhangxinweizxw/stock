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
	if len(d) < 23 {
		logging.Error("日K查询时间 Error：%v", len(d))
		return
	}
	ntime := time.Now().Format("2006-01-02")
	{ // 3
		sql := `SELECT f12,f14 FROM stock_day_k
				WHERE create_time='` + d[3]
		sql += `' AND f3<-5 AND f10>1
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[2]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[1]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[0]
		sql += `' AND f3>3 AND f3 <5 AND f10 > 1
			)))`
		sdkl := stocks_db.NewStock_Day_K().GetDxStockDayKList(sql)
		if len(sdkl) > 0 {
			for _, v := range sdkl {
				//logging.Error("=========", v.F12, v.F14)
				this.Save(v.F12, v.F14, ntime, 3)
			}
		}

	}

	{ // 4
		sql := `SELECT f12,f14 FROM stock_day_k
				WHERE create_time='` + d[4]
		sql += `' AND f3<-5 AND f10>1
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[3]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[2]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[1]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[0]
		sql += `' AND f3>3 AND f3 <5 AND f10 > 1
			))))`
		sdkl := stocks_db.NewStock_Day_K().GetDxStockDayKList(sql)
		if len(sdkl) > 0 {
			for _, v := range sdkl {
				//logging.Error("=========", v.F12, v.F14)
				this.Save(v.F12, v.F14, ntime, 4)
			}
		}

	}

	{ // 5
		sql := `SELECT f12,f14 FROM stock_day_k
				WHERE create_time='` + d[5]
		sql += `' AND f3<-5 AND f10>1
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[4]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[3]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[2]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[1]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[0]
		sql += `' AND f3>3 AND f3 <5 AND f10 > 1
			)))))`
		sdkl := stocks_db.NewStock_Day_K().GetDxStockDayKList(sql)
		if len(sdkl) > 0 {
			for _, v := range sdkl {
				//logging.Error("=========", v.F12, v.F14)
				this.Save(v.F12, v.F14, ntime, 5)
			}
		}

	}

	{ // 7
		sql := `SELECT f12,f14 FROM stock_day_k
				WHERE create_time='` + d[7]
		sql += `' AND f3<-5 AND f10>1
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[6]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[5]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[4]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[3]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[2]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[1]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[0]
		sql += `' AND f3>3 AND f3 <5 AND f10 > 1
			)))))))`
		sdkl := stocks_db.NewStock_Day_K().GetDxStockDayKList(sql)
		if len(sdkl) > 0 {
			for _, v := range sdkl {
				//logging.Error("=========", v.F12, v.F14)
				this.Save(v.F12, v.F14, ntime, 7)
			}
		}

	}

	{ // 8
		sql := `SELECT f12,f14 FROM stock_day_k
				WHERE create_time='` + d[8]
		sql += `' AND f3<-5 AND f10>1
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[7]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[6]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[5]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[4]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[3]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[2]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[1]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[0]
		sql += `' AND f3>3 AND f3 <5 AND f10 > 1
			))))))))`
		sdkl := stocks_db.NewStock_Day_K().GetDxStockDayKList(sql)
		if len(sdkl) > 0 {
			for _, v := range sdkl {
				//logging.Error("=========", v.F12, v.F14)
				this.Save(v.F12, v.F14, ntime, 8)
			}
		}
	}

	{ // 11
		sql := `SELECT f12,f14 FROM stock_day_k
				WHERE create_time='` + d[11]
		sql += `' AND f3<-5 AND f10>1
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[10]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[9]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[8]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[7]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[6]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[5]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[4]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[3]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[2]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[1]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[0]
		sql += `' AND f3>3 AND f3 <5 AND f10 > 1
			)))))))))))`
		sdkl := stocks_db.NewStock_Day_K().GetDxStockDayKList(sql)
		if len(sdkl) > 0 {
			for _, v := range sdkl {
				//logging.Error("=========", v.F12, v.F14)
				this.Save(v.F12, v.F14, ntime, 11)
			}
		}

	}

	{ // 13
		sql := `SELECT f12,f14 FROM stock_day_k 
				WHERE create_time='` + d[13]
		sql += `' AND f3<-5 AND f10>1
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[12]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[11]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[10]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[9]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[8]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[7]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[6]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[5]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[4]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[3]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[2]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[1]
		sql += `' AND ( f3<2 OR f3>-2)
				  AND f12 IN(
				SELECT f12 FROM stock_day_k
				WHERE create_time='` + d[0]
		sql += `' AND f3>3 AND f3 <5 AND f10 > 1 
			)))))))))))))`
		sdkl := stocks_db.NewStock_Day_K().GetDxStockDayKList(sql)
		if len(sdkl) > 0 {
			for _, v := range sdkl {
				//logging.Error("=========", v.F12, v.F14)
				this.Save(v.F12, v.F14, ntime, 11)
			}
		}

	}

	DxStockDb = nil
}

// 短线筛选保存
func (this *DxStock) Save(sc, sn, ntime string, sta int) {

	i := stocks_db.NewDxStockDb()
	p := map[string]interface{}{
		"create_time": ntime,
		"stock_code":  sc,
		"stock_name":  sn,
		"status":      sta,
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
		case "300", "002", "000", "001", "003":
			sc = fmt.Sprintf("SZ%v", v.StockCode)
		}

		i := NewStockDayk(nil).StockInfoSS(sc).StockDate

		name = i.Gpmc
		zljlrv := 0.0
		if reflect.TypeOf(i.Zljlr).String() != "string" {
			zljlrv = i.Zljlr.(float64)
		}
		d1 := decimal.NewFromFloat(zljlrv)
		d2 := decimal.NewFromFloat(i.Jcd)
		d3 := decimal.NewFromFloat(i.Jdd)

		if i.Zdf > 1.28 && i.Zdf < 5 && i.Lb > 1.25 && i.Hsl > 3 && d1.String() > "10000000" && d2.String() > "0" && d3.String() > "0" {
			// 判断是否以入库
			if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
				continue
			}

			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg, 5)
			DxStockDb = append(DxStockDb[:k], DxStockDb[k+1:]...)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))
		}
	}
}
