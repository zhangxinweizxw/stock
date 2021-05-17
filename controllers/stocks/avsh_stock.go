package stocks

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"reflect"
	"stock/config"
	. "stock/models"
	"stock/models/stocks_db"
	"stock/share/logging"
	"stock/share/util"
	"time"
)

type AvsHStock struct {
	C *config.AppConfig
}

func NewAvsHStock(cfg *config.AppConfig) *AvsHStock {
	return &AvsHStock{
		C: cfg,
	}
}

// 保存av对比数据
func (this *AvsHStock) SaveAvsHStock() {

	stocks_db.NewAvsHStock().DelAvshStock()

	url := "http://87.push2.eastmoney.com/api/qt/clist/get?pn=1&pz=50&po=0&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&fid=f189&fs=b:DLMK0101&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152,f191,f192,f193,f186,f185,f187,f189,f188"
	resp, err := http.Get(url)
	if err != nil {
		logging.Error("AvsHStock", err)
	}
	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		logging.Error("ioutil.ReadAll", err1)
	}
	var data *util.StockDayK
	if err = json.Unmarshal(body, &data); err != nil {
		logging.Error("AH 比价  | Error:=", err)
		return
	}
	if data.Datas.Total <= 0 {
		return
	}
	ntime := time.Now().Format("2006-01-02")
	for _, v := range data.Datas.Diff {
		b := stocks_db.NewFinancialReports().GetTbzz(v.F191)
		if b { // 插入数据

			// 股票信息写入stock_info表方便使用
			i := stocks_db.NewAvsHStock()
			p := map[string]interface{}{
				"create_time": ntime,
				"stock_code":  v.F191,
				"stock_name":  v.F193,
			}
			_, err1 := i.Insert(p)
			if err1 != nil {
				logging.Error("Insert Table avsh_stock | %v", err)
				continue
			}
		}
	}
	AvsHStockl = nil
}

// 需求个股分析监控 9：15 - 11：30   13：00-15：00  AH股比价
func (this *AvsHStock) AvsHStockFx() {

	if len(AvsHStockl) <= 0 {
		AvsHStockl = stocks_db.NewAvsHStock().GetAvsHStockList()
	}
	name := ""
	defer func() {
		if err := recover(); err != nil {
			logging.Error("Panic Error=======:%v======:%v", name, err)
		}
	}()

	for k, v := range AvsHStockl {
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
		if i.Zdf > 1.8 && i.Zdf < 5.8 && i.Lb > 1 && i.Lb < 10 && i.Hsl > 1.28 && d1.String() > "10000000" && d2.String() > "1000000" && d3.String() > "500000" {
			// 判断是否以入库
			if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
				continue
			}

			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg, 2)
			AvsHStockl = append(AvsHStockl[:k], AvsHStockl[k+1:]...)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))
		}

	}

}
