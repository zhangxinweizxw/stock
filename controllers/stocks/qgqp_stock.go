package stocks

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"reflect"
	"stock/controllers"
	"stock/share/util"

	. "stock/models"
	"stock/models/stocks_db"
	"stock/share/logging"
	"time"
)

type QgqpStock struct {
	StockCode     string      `json:"Code"`          // 股票代码
	StockName     string      `json:"Name"`          // 股票名称
	New           interface{} `json:"New"`           // 最新股价
	ChangePercent interface{} `json:"ChangePercent"` // 涨跌幅
	PERation      float64     `json:"PERation"`      // 市盈率
	TurnoverRate  interface{} `json:"TurnoverRate"`  // 换手率
}

func NewQgqpStock() *QgqpStock {
	return &QgqpStock{}
}

// 千股千评保存个股
func (this *QgqpStock) QgqpStockSave() {

	stocks_db.NewQgqpStockDb().DelQgqpStock()

	url := "http://dcfm.eastmoney.com/em_mutisvcexpandinterface/api/js/get?type=QGQP_LB&token=70f12f2f4f091e459a279469fe49eca5&cmd=&st=TotalScore&sr=-1&p=1&ps=200&filter=&pageNo=1&pageNum=1"
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
	ntime := time.Now().Format("2006-01-02")
	for _, v := range data {
		if reflect.TypeOf(v.New).String() == "string" || reflect.TypeOf(v.ChangePercent).String() == "string" {
			continue
		}
		if v.New.(float64) > 58 || v.ChangePercent.(float64) > 5.8 || v.ChangePercent.(float64) < 1.8 || v.PERation > 58 || v.TurnoverRate.(float64) < 1.8 || v.TurnoverRate.(float64) > 8 {
			continue
		}

		// 筛选通过   需要判断下最近涨跌和财务数据
		if controllers.NewUtilHttps(nil).GetXqPd(v.StockCode) <= 0 {
			continue
		}

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
		d3 := decimal.NewFromFloat(i.Jdd)
		d2 := "0"
		if reflect.TypeOf(i.Jcd).String() != "string" {
			d2 = fmt.Sprintf("%v", i.Jcd.(float64))
		}
		//  判断最近 涨跌幅 和财务数据
		if controllers.NewUtilHttps(nil).GetXqPd(v.StockCode) <= 0 {
			continue
		}
		if i.Zdf > 0.8 && i.Zdf < 5.8 && i.Lb > 0.8 && i.Lb < 10 && i.Hsl > 1.28 && d1.String() > "10000000" && d2 > "5000000" && d3.String() > "1000000" {
			// 判断是否以入库
			if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
				continue
			}

			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg, 4)
			QgqpStockDb = append(QgqpStockDb[:k], QgqpStockDb[k+1:]...)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))
		}
	}

}
