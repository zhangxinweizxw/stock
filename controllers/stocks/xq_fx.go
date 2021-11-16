package stocks

import (
	"fmt"
	"github.com/shopspring/decimal"
	"reflect"
	"stock/controllers"
	. "stock/models"
	"stock/models/stocks_db"
	"stock/share/logging"
	"stock/share/util"
)

type XqFxStock struct {
	StockCode     string      `json:"Code"`          // 股票代码
	StockName     string      `json:"Name"`          // 股票名称
	New           interface{} `json:"New"`           // 最新股价
	ChangePercent interface{} `json:"ChangePercent"` // 涨跌幅
	PERation      float64     `json:"PERation"`      // 市盈率
	TurnoverRate  interface{} `json:"TurnoverRate"`  // 换手率
	TotalScore    interface{} `json:"TotalScore"`    // 评分
	ZLCB          interface{} `json:"ZLCB"`          // 主力成本
	ZLJLR         interface{} `json:"ZLJLR"`         // 主力净利润
}

func NewXqFxStock() *XqFxStock {
	return &XqFxStock{}
}

func (this *XqFxStock) XqFxTs() {

	if len(XqFxStockDb) <= 0 {
		XqFxStockDb = stocks_db.NewXQ_Stock_FX().GetXqFxStockList()
	}
	name := ""
	defer func() {
		if err := recover(); err != nil {
			logging.Error("Panic Error=======:%v======:%v", name, err)
		}
	}()

	for _, v := range XqFxStockDb {

		sc := controllers.NewUtilHttps(nil).GetUtilCode(v.StockCode)
		if len(sc) <= 0 {
			name = v.StockName
			continue
		}
		sci := controllers.NewUtilHttps(nil).GetUtilCode1(v.StockCode)
		if len(sci) <= 6 {
			name = v.StockName
			continue
		}

		i := NewStockDayk(nil).StockInfoSS(sc).StockDate
		if i == nil {
			continue
		}

		if i.Zdf < -2.8 {
			continue
		}

		//fsd := NewZtStock().GetFsZjlr(sci).Data.KLines
		//if len(fsd) < 5 {
		//	continue
		//}

		if reflect.TypeOf(i.Zljlr).Name() == "string" {
			continue
		}
		//d1 := decimal.NewFromFloat(zljlrv)

		//kl1 := fsd[len(fsd)-1]
		//s1 := strings.Split(kl1, ",")
		////f1, _ := strconv.ParseFloat(s1[1], 64)
		//f1 := fmt.Sprintf("%v", s1[1])[:len(s1[1])-2]
		//kl2 := fsd[len(fsd)-2]
		//s2 := strings.Split(kl2, ",")
		////f2, _ := strconv.ParseFloat(s2[1], 64)
		//f2 := fmt.Sprintf("%v", s2[1])[:len(s2[1])-2]
		//kl3 := fsd[len(fsd)-3]
		//s3 := strings.Split(kl3, ",")
		////f3, _ := strconv.ParseFloat(s3[1], 64)
		//f3 := fmt.Sprintf("%v", s3[1])[:len(s3[1])-2]
		//
		//kl4 := fsd[len(fsd)-4]
		//s4 := strings.Split(kl4, ",")
		////f4, _ := strconv.ParseFloat(s4[1], 64)
		//f4 := fmt.Sprintf("%v", s4[1])[:len(s4[1])-2]
		//kl5 := fsd[len(fsd)-5]
		//s5 := strings.Split(kl5, ",")
		//f5 := fmt.Sprintf("%v", s5[1])[:len(s5[1])-2]
		//
		//kl6 := fsd[0]
		//s6 := strings.Split(kl6, ",")
		//f6 := fmt.Sprintf("%v", s6[1])[:len(s6[1])-2]

		//// 计算涨跌幅
		//// 最高涨跌幅
		//zgzdf := (i.Zgjg - i.Kpj) / i.Kpj
		//
		//zgzdfv, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", zgzdf*100), 64)
		//// 最低涨跌幅
		//zdzdf := (i.Zdjg - i.Kpj) / i.Kpj
		//zdzdfv, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", zdzdf*100), 64)

		// 条件1 高开回调 上涨选
		if reflect.TypeOf(i.Zxjg).Name() == "string" {
			continue
		}
		dzljlr := decimal.NewFromFloat(i.Zljlr.(float64)).String()
		//logging.Debug("name:", v.StockName, "zgzdf:", zgzdfv, "zdzdf:", zdzdfv, "zxjg:", i.Zxjg, "zgjg:", i.Zgjg, "zdjg:", i.Zdjg, "kpj:", i.Kpj, "fffff:", i.Hsl, v.Dayk10)

		// 根据不同市值筛选条件做出改变
		dzljlr01 := ""
		jdd01 := 0.0

		if i.Zsz < 3000000000 { // 市值30亿以内公司 净流入 1千万就很多了
			dzljlr01 = "3800000"
			jdd01 = 1880000
		}
		if i.Zsz > 3000000000 && i.Zsz < 5000000000 { //
			dzljlr01 = "5800000"
			jdd01 = 1880000
		}
		if i.Zsz > 5000000000 && i.Zsz < 15000000000 { //
			dzljlr01 = "8800000"
			jdd01 = 5880000

		}
		if i.Zsz > 15000000000 { //
			dzljlr01 = "12800000"
			jdd01 = 8880000

		}

		if i.Zgjg > i.Kpj && dzljlr > dzljlr01 && i.Jdd.(float64) > jdd01 && i.Zxjg.(float64) > i.Kpj && i.Zdf < 5.8 && i.Hsl > 1.28 && i.Lb > 1.28 {
			// 判断是否已入库
			if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
				continue
			}

			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg.(float64), 9)
			//XqFxStockDb = append(XqFxStockDb[:k], XqFxStockDb[k+1:]...)
			logging.Debug("=xq_fx")
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))

		}
	}

}
