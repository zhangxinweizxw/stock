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

type ZjlxStock struct {
}

func NewZjlxStock() *ZjlxStock {
	return &ZjlxStock{}
}

// 资金流向保存个股
func (this *ZjlxStock) ZjlxStockSave() {

	stocks_db.NewZjlxStockDb().DelZjlxStock()

	url := "http://push2.eastmoney.com/api/qt/clist/get?fid=f62&po=1&pz=568&pn=1&np=1&fltt=2&invt=2&ut=b2884a393a59ad64002292a3e90d46a5&fs=m%3A0%2Bt%3A6%2Bf%3A!2%2Cm%3A0%2Bt%3A13%2Bf%3A!2%2Cm%3A0%2Bt%3A80%2Bf%3A!2%2Cm%3A1%2Bt%3A2%2Bf%3A!2%2Cm%3A1%2Bt%3A23%2Bf%3A!2%2Cm%3A0%2Bt%3A7%2Bf%3A!2%2Cm%3A1%2Bt%3A3%2Bf%3A!2&fields=f12%2Cf14%2Cf2%2Cf3%2Cf62%2Cf184%2Cf66%2Cf69%2Cf72%2Cf75%2Cf78%2Cf81%2Cf84%2Cf87%2Cf204%2Cf205%2Cf124%2Cf10"
	resp, err := http.Get(url)
	if err != nil {
		logging.Error("ZjlxStock:", err)
	}
	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		logging.Error("ioutil.ReadAll", err1)
	}
	var data *util.StockDayK
	if err = json.Unmarshal(body, &data); err != nil {
		logging.Error("个股资金流向  | Error:=", err)
		return
	}
	if data.Datas.Total <= 0 {
		return
	}
	ntime := time.Now().Format("2006-01-02")
	for _, v := range data.Datas.Diff {

		if v.F3.(float64) > 5.8 || v.F3.(float64) < 1.8 || v.F2.(float64) > 58 || v.F2.(float64) < 5 || v.F184.(float64) < 5 || v.F69.(float64) < 5 || v.F62.(float64) < 10000000 || v.F66.(float64) < 5000000 || v.F10.(float64) < 1 {
			continue
		}
		// 筛选通过   需要判断下最近涨跌和财务数据
		if controllers.NewUtilHttps(nil).GetXqPd(v.F12.(string)) <= 0 {
			continue
		}

		// 股票信息写入stock_info表方便使用
		i := stocks_db.NewZjlxStockDb()
		p := map[string]interface{}{
			"create_time": ntime,
			"stock_code":  v.F12,
			"stock_name":  v.F14,
		}
		_, err1 := i.Insert(p)
		if err1 != nil {
			logging.Error("Insert Table zjlx_stock | %v", err)
			continue
		}

	}
	ZjlxStockDb = nil
}

// 个股资金实时流向 分析是否卖出
func (this *ZjlxStock) ZjlxStockSellFx() {

	// 查询表中数据
	scl := stocks_db.NewTransactionHistory().GetTranHistWmcList()
	if len(scl) <= 0 {
		return
	}

	for _, v := range scl {
		//logging.Error("=====", v)
		stockCodes := ""
		switch v.StockCode[:3] {
		case "600", "601", "603", "605", "688", "689", "608":
			stockCodes = fmt.Sprintf("1.%v", v.StockCode)
		case "300", "002", "000", "001", "003", "301":
			stockCodes = fmt.Sprintf("0.%v", v.StockCode)
		default:
			continue
		}
		s1 := this.ZjlxStockInfo(stockCodes)

		np := this.StockInfo(stockCodes)

		bfb := (np - v.BuyPrice) / v.BuyPrice

		if bfb < -0.03 { // 跌 3% 卖出
			stocks_db.NewTransactionHistory().UpdateTranHist(v.StockCode, np, bfb*100)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议卖出：%v   |   股票代码：%v    卖出价：%v", v.StockName, v.StockCode, np))
			continue
		}
		df := decimal.NewFromFloat(s1.F62.(float64))

		//if (df.String() < "-5000000") && (s1.F66.(float64) < 0) && s1.F10.(float64) > 0.5 {
		//	stocks_db.NewTransactionHistory().UpdateTranHist(v.StockCode, np, bfb*100)
		//	go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议卖出：%v   |   股票代码：%v    卖出价：%v", v.StockName, v.StockCode, np))
		//}

		// 放量下跌  跌破10均线 主力资金流出5000000 大资金流出 <0
		dk10 := stocks_db.NewStock_Day_K().GetStockDayK10Date(v.StockCode)
		if (df.String() < "-5000000" && np < dk10) || (s1.F3.(float64) < -3.0 && s1.F10.(float64) > 1) {
			stocks_db.NewTransactionHistory().UpdateTranHist(v.StockCode, np, bfb*100)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议卖出：%v   |   股票代码：%v    卖出价：%v", v.StockName, v.StockCode, np))
		}

	}
}

// 个股资金实时流向 获取
func (this *ZjlxStock) ZjlxStockInfo(stockCode string) *util.StockInfo {

	url := "http://push2.eastmoney.com/api/qt/ulist.np/get?fltt=2&secids=" + stockCode + "&fields=f62%2Cf184%2Cf66%2Cf69%2Cf72%2Cf75%2Cf78%2Cf81%2Cf10&ut=b2884a393a59ad64002292a3e90d46a5"
	//logging.Error("=======%v", url)
	resp, err := http.Get(url)
	if err != nil {
		logging.Error("个股资金流向", err)
	}
	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		logging.Error("ioutil.ReadAll", err1)
	}
	//logging.Error("=======%v", string(body))
	var data *util.StockDayK
	if err = json.Unmarshal(body, &data); err != nil {
		logging.Error("个股实时资金流向  | Error:=", err)
		return nil
	}

	if len(data.Datas.Diff) <= 0 {
		return nil
	}
	return data.Datas.Diff[0]
}

type StockData struct {
	DataI *D `json:"data"`
}
type D struct {
	F43 float64 `json:"f43"`
}

// 个股最近价格
func (this *ZjlxStock) StockInfo(stockCode string) float64 {
	url := fmt.Sprintf("http://push2.eastmoney.com/api/qt/stock/get?ut=fa5fd1943c7b386f172d6893dbfba10b&invt=2&fltt=2&fields=f43&secid=%v", stockCode)
	resp, err := http.Get(url)
	if err != nil {
		logging.Error("个股资金流向", err)
	}
	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		logging.Error("ioutil.ReadAll", err1)
	}
	var data *StockData
	if err = json.Unmarshal(body, &data); err != nil {
		logging.Error("个股最新价格  | Error:=", err)
		return 0
	}

	return data.DataI.F43
}

// 需求个股分析监控 9：15 - 11：30   13：00-15：00  资金流向
func (this *ZjlxStock) ZjlxtockFx() {

	if len(ZjlxStockDb) <= 0 {
		ZjlxStockDb = stocks_db.NewZjlxStockDb().GetZjlxStockList()
	}
	name := ""
	defer func() {
		if err := recover(); err != nil {
			logging.Error("Panic Error=======:%v======:%v", name, err)
		}
	}()

	for k, v := range ZjlxStockDb {

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
			d2 = fmt.Sprintf("%v", decimal.NewFromFloat(i.Jcd.(float64)))
		}
		//  判断最近 涨跌幅 和财务数据
		if controllers.NewUtilHttps(nil).GetXqPd(v.StockCode) <= 0 {
			continue
		}

		if i.Zdf > 0.8 && i.Zdf < 5.8 && i.Lb > 0.8 && i.Lb < 8 && i.Hsl > 1.28 && d1.String() > "5000000" && d2 > "1000000" && d3.String() > "500000" {
			// 判断是否以入库
			if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
				continue
			}
			// 判断是否以入库
			if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
				break
			}
			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(v.StockCode, v.StockName, i.Zxjg, 3)
			ZjlxStockDb = append(ZjlxStockDb[:k], ZjlxStockDb[k+1:]...)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", v.StockName, v.StockCode, i.Zxjg))

		}

	}

}

type PkydData struct {
	Data AlData `json:"data"`
}
type AlData struct {
	AData []*AllStock `json:"allstock"`
}
type AllStock struct {
	Time        int32  `json:"tm"` // 小时分钟秒
	StockCode   string `json:"c"`  // 股票代码
	StockMarket int    `json:"m"`  // 1、沪市  0、深市
	StockName   string `json:"n"`  // 股票名称
	Type        int    `json:"t"`  //  大笔买入：8193  火箭发射：8201   快速反弹：8202
	Value       string `json:"i"`  // t=8193 手   t=8201  百分比  t=8220 百分比
}

// 盘口异动
func (this *ZjlxStock) PkydStockFx() {

	url := "http://push2ex.eastmoney.com/getAllStockChanges?type=8201,8202,8193&pageindex=0&pagesize=3&ut=7eea3edcaed734bea9cbfc24409ed989&dpt=wzchanges"
	resp, err := http.Get(url)
	if err != nil {
		logging.Error("pkydStock:", err)
	}
	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		logging.Error("ioutil.ReadAll", err1)
	}
	var data *PkydData
	if err = json.Unmarshal(body, &data); err != nil {
		logging.Error("盘口异动  | Error:=", err)
		return
	}
	//logging.Error("======", len(data.Data.AData))
	if len(data.Data.AData) <= 0 {
		return
	}
	for _, v := range data.Data.AData {
		// 判断是否以入库
		if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
			break
		}

		d := NewStockDayk(nil).GetZJlxDFCF(v.StockCode).Datas.Diff[0]

		df62 := decimal.NewFromFloat(d.F62.(float64)).String()
		df66 := decimal.NewFromFloat(d.F66.(float64)).String()
		df72 := decimal.NewFromFloat(d.F72.(float64)).String()

		//logging.Error("=========:", df62, "====:", d.F184, "=====:", df66, "====:", d.F69, "====:", df72, "====:", d.F75, "=====:", d.F2, "=====:", d.F8, "====:", d.F9, "====:", d.F10)
		if (df62 < "5800000") || (df66 < "1880000") || (df72 < "580000") || d.F2.(float64) > 58 || (d.F8.(float64) < 1.28 || d.F8.(float64) > 8) || (d.F9.(float64) < 5.8 || d.F9.(float64) > 68) || d.F10.(float64) < 0.8 || d.F3.(float64) > 3.8 || d.F3.(float64) < 0.28 {
			continue
		}
		// 筛选通过   需要判断下最近涨跌和财务数据
		if controllers.NewUtilHttps(nil).GetXqPd(v.StockCode) <= 0 {
			continue
		}
		// 满足条件   mysql transaction_history 表中添加数据 // 发送叮叮实时消息
		go NewStockDayk(nil).SaveStock(v.StockCode, v.StockName, d.F2.(float64), 6)

		go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", v.StockCode, v.StockName, d.F2))
	}
}
