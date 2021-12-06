package stocks

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"reflect"
	"stock/controllers"
	. "stock/models"
	"stock/models/stocks_db"
	"stock/share/logging"
	"stock/share/util"
	"time"
)

type QgqpStock struct {
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

func NewQgqpStock() *QgqpStock {
	return &QgqpStock{}
}

// 千股千评保存个股
func (this *QgqpStock) QgqpStockSave() {

	stocks_db.NewQgqpStockDb().DelQgqpStock()

	defer func() {
		if err := recover(); err != nil {
			logging.Error("Panic Error=======:%v======:%v", "保存千股千评数据", err)
		}
	}()
	url := "http://dcfm.eastmoney.com/em_mutisvcexpandinterface/api/js/get?type=QGQP_LB&token=70f12f2f4f091e459a279469fe49eca5&cmd=&st=RankingUp&sr=-1&p=1&ps=1880"
	//url := "http://dcfm.eastmoney.com/em_mutisvcexpandinterface/api/js/get?type=QGQP_LB&token=70f12f2f4f091e459a279469fe49eca5&cmd=&st=TotalScore&sr=-1&p=1&ps=3000"
	resp, err := http.Get(url)
	if err != nil {
		logging.Error("qgqpStock:", err)
	}
	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		logging.Error("ioutil.ReadAll", err1)
	}
	var data []*QgqpStock
	if err = json.Unmarshal(body, &data); err != nil {
		logging.Error("千股千评  | Error:=", err)
		return
	}
	if len(data) <= 0 {
		return
	}

	var d []*QgqpStock
	for _, v := range data {
		if reflect.TypeOf(v.New).Name() == "string" || reflect.TypeOf(v.ZLCB).Name() == "string" {
			continue
		}

		f := (v.New.(float64) - v.ZLCB.(float64)) / v.ZLCB.(float64)

		if f < 0.088 && f >= 0.0128 {
			d = append(d, v)
		}

	}

	ntime := time.Now().Format("2006-01-02")

	for _, v := range d {

		if reflect.TypeOf(v.New).String() == "string" || reflect.TypeOf(v.ChangePercent).String() == "string" || reflect.TypeOf(v.TotalScore).String() == "string" {
			continue
		}
		if NewStockDayk(nil).GetReturnIsBuy(v.StockCode) == false {
			continue
		}

		if v.New.(float64) > 88 || v.ChangePercent.(float64) > 5.8 || v.ChangePercent.(float64) < 1.28 || v.PERation > 80 || v.TurnoverRate.(float64) < 2.8 || v.TurnoverRate.(float64) > 10 || v.TotalScore.(float64) < 68 {
			continue
		}

		// 筛选通过   需要判断下最近涨跌和财务数据
		//if controllers.NewUtilHttps(nil).GetXqPd(v.StockCode) <= 0 {
		//	continue
		//}

		// 股票信息写入qgqp_stock表方便使用
		i := stocks_db.NewQgqpStockDb()
		p := map[string]interface{}{
			"create_time": ntime,
			"stock_code":  v.StockCode,
			"stock_name":  v.StockName,
		}
		_, err1 := i.Insert(p)
		if err1 != nil {
			logging.Error("Insert Table qgqp_stock | %v", err)
			continue
		}

	}
	QgqpStockDb = nil
}

// 需求个股分析监控 9：15 - 11：30   13：00-15：00  资金流向
func (this *QgqpStock) QgqpStockFx() {

	if len(QgqpStockDb) <= 0 {
		QgqpStockDb = stocks_db.NewQgqpStockDb().GetQgqpStockList()
	}
	name := ""
	defer func() {
		if err := recover(); err != nil {
			logging.Error("Panic Error=======:%v======:%v", name, err)
		}
	}()

	for k, v := range QgqpStockDb {

		sc := controllers.NewUtilHttps(nil).GetUtilCode(v.StockCode)
		if len(sc) <= 0 {
			continue
		}

		i := NewStockDayk(nil).StockInfoSS(sc).StockDate
		if i == nil {
			continue
		}
		name = i.Gpmc
		zljlrv := 0.0
		if reflect.TypeOf(i.Zljlr).String() != "string" {
			zljlrv = i.Zljlr.(float64)
		}
		d1 := decimal.NewFromFloat(zljlrv)

		d2 := "0"
		if reflect.TypeOf(i.Jcd).String() != "string" {
			d2 = fmt.Sprintf("%v", decimal.NewFromFloat(i.Jcd.(float64)))
		}
		d3 := "0"
		if reflect.TypeOf(i.Jdd).String() != "string" {
			d3 = fmt.Sprintf("%v", decimal.NewFromFloat(i.Jdd.(float64)))
		}
		//  判断最近 涨跌幅 和财务数据
		if controllers.NewUtilHttps(nil).GetXqPd(v.StockCode) <= 0 {
			continue
		}

		d101, d201, d301 := "", "", ""
		if i.Zsz < 3000000000 { // 市值30亿以内公司 净流入 1千万就很多了
			d101 = "5880000"
			d201 = "1280000"
			d301 = "880000"
		}
		if i.Zsz > 3000000000 && i.Zsz < 5000000000 { //
			d101 = "8880000"
			d201 = "3288000"
			d301 = "1280000"
		}
		if i.Zsz > 5000000000 && i.Zsz < 15000000000 { //
			d101 = "12880000"
			d201 = "5880000"
			d301 = "3288000"
		}
		if i.Zsz > 15000000000 { //
			d101 = "32880000"
			d201 = "8880000"
			d301 = "3288000"
		}
		if i.Zdf > 0.5 && i.Zdf < 3.8 && i.Lb > 1.28 && i.Lb < 10 && i.Hsl > 1.28 && d1.String() > d101 && d2 > d201 && d3 > d301 {
			// 判断是否以入库
			if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
				continue
			}
			if reflect.TypeOf(i.Zxjg).Name() == "string" {
				continue
			}
			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg.(float64), 4)
			QgqpStockDb = append(QgqpStockDb[:k], QgqpStockDb[k+1:]...)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))
		}
	}

}
