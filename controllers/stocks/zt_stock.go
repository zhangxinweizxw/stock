package stocks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"stock/controllers"
	. "stock/models"
	"stock/models/stocks_db"
	"stock/share/logging"
	"stock/share/util"
	"strconv"
	"strings"
	"time"
)

type ZtStock struct {
}

func NewZtStock() *ZtStock {
	return &ZtStock{}
}

func (this *ZtStock) ZtStockFx() {

	if len(ZtStockDb) <= 0 {
		ZtStockDb = stocks_db.NewZtStockDB().GetZtStockList()
		if ZtStockDb == nil {
			return
		}
	}
	name := ""
	defer func() {
		if err := recover(); err != nil {
			logging.Error("Panic Error=======:%v======:%v", name, err)
		}
	}()

	for k, v := range ZtStockDb {
		//sc := ""
		//sci := ""
		//switch v.StockCode[:3] {
		//case "600", "601", "603", "605", "688", "689", "608":
		//	sc = fmt.Sprintf("SH%v", v.StockCode)
		//	sci = fmt.Sprintf("1.%v", v.StockCode)
		//case "300", "002", "000", "001", "003", "301":
		//	sc = fmt.Sprintf("SZ%v", v.StockCode)
		//	sci = fmt.Sprintf("0.%v", v.StockCode)
		//default:
		//	continue
		//}
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

		if i.Zdf < 0 {
			ZtStockDb = append(ZtStockDb[:k], ZtStockDb[k+1:]...)
			continue
		}
		//if i.Zdf < -1.28 || i.Zdf > 3.8 || fs > 4 {
		//	ZtStockDb = append(ZtStockDb[:k], ZtStockDb[k+1:]...)
		//	go stocks_db.NewZtStockDB().DelZtStockTj(v.StockCode)
		//	continue
		//}
		fsd := this.GetFsZjlr(sci).Data.KLines
		if len(fsd) < 5 {
			continue
		}

		if reflect.TypeOf(i.Zljlr).Name() == "string" {
			continue
		}
		//d1 := decimal.NewFromFloat(zljlrv)

		kl1 := fsd[len(fsd)-1]
		s1 := strings.Split(kl1, ",")
		f1, _ := strconv.ParseFloat(s1[1], 64)
		kl2 := fsd[len(fsd)-2]
		s2 := strings.Split(kl2, ",")
		f2, _ := strconv.ParseFloat(s2[1], 64)

		kl3 := fsd[len(fsd)-3]
		s3 := strings.Split(kl3, ",")
		f3, _ := strconv.ParseFloat(s3[1], 64)

		kl4 := fsd[len(fsd)-4]
		s4 := strings.Split(kl4, ",")
		f4, _ := strconv.ParseFloat(s4[1], 64)

		kl5 := fsd[len(fsd)-5]
		s5 := strings.Split(kl5, ",")
		f5, _ := strconv.ParseFloat(s5[1], 64)

		// 计算涨跌幅
		// 最高涨跌幅
		zgzdf := int((i.Zgjg - i.Kpj) / i.Kpj * 100)
		zgzdfv, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", zgzdf), 64)
		// 最低涨跌幅
		zdzdf := int((i.Zdjg - i.Kpj) / i.Kpj * 100)
		zdzdfv, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", zdzdf), 64)
		// 条件1 高开回调 上涨选
		if zgzdfv > 2.8 && zgzdfv < 8 && i.Zxjg > i.Kpj && i.Zljlr.(float64) > 8000000 && i.Zljlr.(float64) > f3 && i.Zxjg < i.Zgjg && i.Zxjg > i.Zdjg {
			// 判断是否已入库
			if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
				continue
			}

			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg, 2)
			ZtStockDb = append(ZtStockDb[:k], ZtStockDb[k+1:]...)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))

		}

		zlf := i.Zljlr.(float64)

		if zdzdfv >= 0.8 && zgzdfv < 6 && zlf > 10000000 && i.Zxjg > i.Zdjg && f1/2 >= f3 && f2/2 >= f4 && f1 >= f2 && f2 >= f3 && f3 >= f4 && f4 >= f5 {
			// 判断是否已入库
			if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
				continue
			}

			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg, 2)
			ZtStockDb = append(ZtStockDb[:k], ZtStockDb[k+1:]...)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))

		}

		////name = i.Gpmc
		//zljlrv := 0.0
		//if reflect.TypeOf(i.Zljlr).String() != "string" {
		//	zljlrv = i.Zljlr.(float64)
		//}
		//d1 := decimal.NewFromFloat(zljlrv)
		////d3 := decimal.NewFromFloat(i.Jdd)
		//d2 := "0"
		//if reflect.TypeOf(i.Jcd).String() != "string" {
		//	d2 = fmt.Sprintf("%v", i.Jcd.(float64))
		//}
		////if i.Zdf > 1.8 && i.Zdf < 5.8 && i.Lb > 1 && i.Lb < 10 && i.Hsl > 1.28 && d1.String() > "10000000" && d2 > "1000000" && d3.String() > "500000" {
		//if i.Zdjg > v.Dayk20 && i.Zxjg > v.Dayk10 && i.Zdf > 0 && i.Zdf < 5.7 && i.Lb > 0.28 && i.Lb < 10 && d1.String() > "10000000" && d2 > "5000000" {
		//	// 判断是否已入库
		//	if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
		//		continue
		//	}
		//
		//	// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
		//	go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg, 2)
		//	ZtStockDb = append(ZtStockDb[:k], ZtStockDb[k+1:]...)
		//	//go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))
		//}

		// 条件2 平开或者低开 然后资金流入 加速
		if i.Zdf > 0.5 && f1 >= 8800000 && f2 > 5000000 && f3 > 3800000 && i.Zdf < 5.8 {
			// 判断是否已入库
			if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
				continue
			}

			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg, 2)
			ZtStockDb = append(ZtStockDb[:k], ZtStockDb[k+1:]...)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))

		}
	}

}

// 抓涨停实验
func (this *ZtStock) GetZTStock() {

	stocks_db.NewZtStockDB().DelZtStock()

	url := "http://73.push2.eastmoney.com/api/qt/clist/get?pn=1&pz=880&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&fid=f3&fs=m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152"
	resp, err := http.Get(url)

	if err != nil {
		logging.Error("ztstock Error", err)
	}

	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		logging.Error("ioutil.ReadAll", err1)
	}

	var data *util.StockDayK
	if err = json.Unmarshal(body, &data); err != nil {
		logging.Error("解析日K  | Error:=", err)
	}
	d := data.Datas.Diff
	//l := len(d)
	i2 := int(len(d) / 2)
	go func() {
		for i, v := range d {
			if i >= 0 && i < i2 {
				if v.F12.(string)[:3] == "688" || v.F12.(string)[:2] == "ST" || v.F12.(string)[:3] == "*ST" {
					continue
				}
				//f62 := decimal.NewFromFloat(v.F62.(float64))
				if v.F3.(float64) < -0.5 || v.F3.(float64) > 8 || v.F2.(float64) > 58 || v.F62.(float64) < 0 || v.F23.(float64) < 1.28 || v.F23.(float64) > 10 {
					continue
				}
				d := stocks_db.NewStock_Day_K().GetStockDayKJJ(v.F12.(string))
				if (d.DayK30 > d.DayK10 || d.DayK20 > d.DayK5 || d.DayK10 > d.DayK5) || d.Day5Zdf > 6 || d.Day5Zdf < -3 || d.Day20Zdf < -10 || d.Day20Zdf > 15 {
					continue
				}

				ntime := time.Now().Format("2006-01-02")

				// 股票信息写入zt_stock表方便使用
				i := stocks_db.NewZtStockDB()
				p := map[string]interface{}{
					"create_time": ntime,
					"stock_code":  v.F12,
					"stock_name":  v.F14,
					"dayk5":       d.DayK5,
					"dayk10":      d.DayK10,
					"dayk20":      d.DayK20,
					"dayk30":      d.DayK30,
				}
				_, err := i.Insert(p)
				if err != nil {
					logging.Error("Insert Table zt_stock | %v", err)
					continue
				}

			}
		}
	}()
	go func() {
		for i, v := range d {
			if i >= i2 && i < len(d)-1 {
				if reflect.TypeOf(v.F3).Name() == "string" {
					continue
				}
				if v.F3.(float64) < -0.5 || v.F3.(float64) > 8 || v.F2.(float64) > 58 || v.F62.(float64) < 0 || v.F23.(float64) < 1.28 || v.F23.(float64) > 10 {
					continue
				}
				d := stocks_db.NewStock_Day_K().GetStockDayKJJ(v.F12.(string))
				if (d.DayK30 > d.DayK10 || d.DayK20 > d.DayK5 || d.DayK10 > d.DayK5) || d.Day5Zdf > 6 || d.Day5Zdf < -3 || d.Day20Zdf < -10 || d.Day20Zdf > 15 {
					continue
				}

				ntime := time.Now().Format("2006-01-02")

				// 股票信息写入zt_stock表方便使用
				i := stocks_db.NewZtStockDB()
				p := map[string]interface{}{
					"create_time": ntime,
					"stock_code":  v.F12,
					"stock_name":  v.F14,
					"dayk5":       d.DayK5,
					"dayk10":      d.DayK10,
					"dayk20":      d.DayK20,
					"dayk30":      d.DayK30,
				}
				_, err := i.Insert(p)
				if err != nil {
					logging.Error("Insert Table zt_stock | %v", err)
					continue
				}

			}
		}
	}()
	ZtStockDb = nil
}

type FRul struct {
	Data *FData `json:"data"`
}
type FData struct {
	KLines []string `json:"klines"`
}

// 分时资金流入
func (this *ZtStock) GetFsZjlr(sc string) *FRul {

	url := "http://push2.eastmoney.com/api/qt/stock/fflow/kline/get?lmt=0&klt=1&fields1=f1%2Cf2%2Cf3%2Cf7&fields2=f51%2Cf52%2Cf53%2Cf54%2Cf55%2Cf56%2Cf57%2Cf58%2Cf59%2Cf60%2Cf61%2Cf62%2Cf63%2Cf64%2Cf65&ut=b2884a393a59ad64002292a3e90d46a5&secid=" + sc
	resp, err := http.Get(url)

	if err != nil {
		logging.Error("ztstock Error", err)
	}

	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		logging.Error("ioutil.ReadAll", err1)
	}

	var d *FRul
	if err = json.Unmarshal(body, &d); err != nil {
		logging.Error("实时资金流向 | Error:=", err)
	}
	return d
	//logging.Error("=====", len(d.Data.KLines))
}
