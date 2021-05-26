package stocks

import (
	"fmt"
	"github.com/shopspring/decimal"
	"reflect"
	. "stock/models"
	"stock/models/stocks_db"
	"stock/share/logging"
	"time"
)

type ZtStock struct {
}

func NewZtStock() *ZtStock {
	return &ZtStock{}
}

// 抓涨停实验
func (this *ZtStock) GetZTStock() {

	stocks_db.NewZtStockDB().DelZtStock()

	sdk := stocks_db.NewStock_Day_K().ZtSelStockDayk()

	ntime := time.Now().Format("2006-01-02")
	for _, v := range sdk {
		// 股票信息写入zt_stock表方便使用
		i := stocks_db.NewZtStockDB()
		p := map[string]interface{}{
			"create_time": ntime,
			"stock_code":  v.F12,
			"stock_name":  v.F14,
			"dayk5":       v.DayK5,
			"dayk10":      v.DayK10,
			"dayk20":      v.DayK20,
			"dayk30":      v.DayK30,
		}
		_, err := i.Insert(p)
		if err != nil {
			logging.Error("Insert Table avsh_stock | %v", err)
			continue
		}
	}

	ZtStockDb = nil
}

func (this *ZtStock) ZtStockFx() {

	if len(ZtStockDb) <= 0 {
		ZtStockDb = stocks_db.NewZtStockDB().GetZtStockList()
	}
	name := ""
	defer func() {
		if err := recover(); err != nil {
			logging.Error("Panic Error=======:%v======:%v", name, err)
		}
	}()

	for k, v := range ZtStockDb {
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

		if i.Zdf < -2 || i.Zdf > 5.8 {
			ZtStockDb = append(ZtStockDb[:k], ZtStockDb[k+1:]...)
			go stocks_db.NewZtStockDB().DelZtStockTj(v.StockCode)
			continue
		}

		name = i.Gpmc
		zljlrv := 0.0
		if reflect.TypeOf(i.Zljlr).String() != "string" {
			zljlrv = i.Zljlr.(float64)
		}
		d1 := decimal.NewFromFloat(zljlrv)
		//d3 := decimal.NewFromFloat(i.Jdd)
		d2 := "0"
		if reflect.TypeOf(i.Jcd).String() != "string" {
			d2 = fmt.Sprintf("%v", i.Jcd.(float64))
		}
		//if i.Zdf > 1.8 && i.Zdf < 5.8 && i.Lb > 1 && i.Lb < 10 && i.Hsl > 1.28 && d1.String() > "10000000" && d2 > "1000000" && d3.String() > "500000" {
		if i.Zxjg > v.Dayk5 && i.Zdf > 0 && i.Zdf < 5.8 && i.Lb < 10 && d1.String() > "10000000" && d2 > "5000000" {
			// 判断是否已入库
			if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
				continue
			}

			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg, 2)
			ZtStockDb = append(ZtStockDb[:k], ZtStockDb[k+1:]...)
			//go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))
		}

	}

}
