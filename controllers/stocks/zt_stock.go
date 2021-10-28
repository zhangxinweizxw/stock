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
		if k >= k+1 {
			continue
		}
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

		if i.Zdf < 0.28 {
			ZtStockDb = append(ZtStockDb[:k], ZtStockDb[k+1:]...)
			continue
		}

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
		//f1, _ := strconv.ParseFloat(s1[1], 64)
		f1 := fmt.Sprintf("%v", s1[1])[:len(s1[1])-2]
		kl2 := fsd[len(fsd)-2]
		s2 := strings.Split(kl2, ",")
		//f2, _ := strconv.ParseFloat(s2[1], 64)
		f2 := fmt.Sprintf("%v", s2[1])[:len(s2[1])-2]
		kl3 := fsd[len(fsd)-3]
		s3 := strings.Split(kl3, ",")
		//f3, _ := strconv.ParseFloat(s3[1], 64)
		f3 := fmt.Sprintf("%v", s3[1])[:len(s3[1])-2]

		kl4 := fsd[len(fsd)-4]
		s4 := strings.Split(kl4, ",")
		//f4, _ := strconv.ParseFloat(s4[1], 64)
		f4 := fmt.Sprintf("%v", s4[1])[:len(s4[1])-2]
		kl5 := fsd[len(fsd)-5]
		s5 := strings.Split(kl5, ",")
		f5 := fmt.Sprintf("%v", s5[1])[:len(s5[1])-2]

		kl6 := fsd[0]
		s6 := strings.Split(kl6, ",")
		f6 := fmt.Sprintf("%v", s6[1])[:len(s6[1])-2]

		// 计算涨跌幅
		// 最高涨跌幅
		zgzdf := (i.Zgjg - i.Kpj) / i.Kpj

		zgzdfv, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", zgzdf*100), 64)
		// 最低涨跌幅
		zdzdf := (i.Zdjg - i.Kpj) / i.Kpj
		zdzdfv, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", zdzdf*100), 64)

		// 条件1 高开回调 上涨选
		if reflect.TypeOf(i.Zxjg).Name() == "string" {
			continue
		}
		dzljlr := decimal.NewFromFloat(i.Zljlr.(float64)).String()
		logging.Debug("name:", v.StockName, "zgzdf:", zgzdfv, "zdzdf:", zdzdfv, "zxjg:", i.Zxjg, "zgjg:", i.Zgjg, "zdjg:", i.Zdjg, "kpj:", i.Kpj, "fffff:", i.Hsl, v.Dayk10)

		if i.Zgjg > i.Kpj && dzljlr > "50000000" && i.Jdd.(float64) > 10000000 && i.Zxjg.(float64) > i.Kpj && i.Zdf < 5.8 && i.Hsl > 1 && i.Zxjg.(float64) >= v.Dayk10 {
			// 判断是否已入库
			if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
				continue
			}

			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg.(float64), 2)
			ZtStockDb = append(ZtStockDb[:k], ZtStockDb[k+1:]...)
			logging.Debug("=55555")
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))

		}

		if zgzdfv > 0.28 && zgzdfv < 3.8 && i.Zxjg.(float64) > i.Kpj && dzljlr > "38800000" && dzljlr >= f1 && f1 > f2 && f2 > f3 && f3 > f4 && i.Zxjg.(float64) < i.Zgjg && i.Zxjg.(float64) > i.Zdjg && i.Lb > 1.8 && f6 > "3880000" && f6 < f1 && f6 < f3 && i.Zxjg.(float64) >= v.Dayk10 {
			// 判断是否已入库
			if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
				continue
			}

			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg.(float64), 2)
			ZtStockDb = append(ZtStockDb[:k], ZtStockDb[k+1:]...)
			logging.Debug("=11111")
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))

		}
		f1s1, _ := strconv.ParseFloat(f1, 64)
		f1s2 := decimal.NewFromFloat(f1s1 / 2).String()
		f2s1, _ := strconv.ParseFloat(f2, 64)
		f2s2 := decimal.NewFromFloat(f2s1 / 2).String()
		f3s1, _ := strconv.ParseFloat(f3, 64)
		f3s2 := decimal.NewFromFloat(f3s1 / 2).String()
		if zdzdfv >= 0.28 && zgzdfv < 5.8 && dzljlr > "12800000" && i.Zxjg.(float64) > i.Zdjg && f1s2 >= f3 && f2s2 >= f4 && f3s2 >= f5 && f1 >= f2 && f2 >= f3 && f3 >= f4 && f4 >= f5 && i.Lb > 1.5 {
			// 判断是否已入库
			if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
				continue
			}

			logging.Debug("=22222")
			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg.(float64), 2)
			ZtStockDb = append(ZtStockDb[:k], ZtStockDb[k+1:]...)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))

		}

		// 条件2 平开或者低开 然后资金流入 加速

		if i.Zdf > 0.8 && f1 >= "12800000" && f2 >= "8800000" && f3 >= "5800000" && f4 >= "3800000" && f5 > "0" && i.Zdf < 3.8 && i.Lb > 2 && (zgzdfv-i.Zdf) < 1.8 && i.Zxjg.(float64) > i.Zdjg && i.Zxjg.(float64) >= v.Dayk10 && i.Zdjg < v.Dayk5 && f6 > "0" && f6 < f1 {
			// 判断是否已入库
			if stocks_db.NewTransactionHistory().GetTranHist(v.StockCode) > 0 {
				continue
			}
			logging.Debug("=33333")

			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息
			go NewStockDayk(nil).SaveStock(i.Gpdm, i.Gpmc, i.Zxjg.(float64), 2)
			ZtStockDb = append(ZtStockDb[:k], ZtStockDb[k+1:]...)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))

		}

	}

}

// 抓涨停实验
func (this *ZtStock) GetZTStock() {

	//stocks_db.NewZtStockDB().DelZtStock()

	url := "http://73.push2.eastmoney.com/api/qt/clist/get?pn=1&pz=1280&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&fid=f3&fs=m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152"
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
				if reflect.TypeOf(v.F3).Name() == "string" {
					continue
				}
				if reflect.TypeOf(v.F2).Name() == "string" {
					continue
				}
				if reflect.TypeOf(v.F62).Name() == "string" {
					continue
				}
				if v.F3.(float64) < 0.28 || v.F3.(float64) > 7 || v.F2.(float64) > 68 || v.F62.(float64) < 1000000 {
					continue
				}
				d := stocks_db.NewStock_Day_K().GetStockDayKJJ(v.F12.(string))
				if (d.DayK30 > d.DayK20 && d.DayK20 > d.DayK10) || d.Day5Zdf > 8 || d.Day5Zdf < 1 || d.Day20Zdf < 0 || d.Day20Zdf > 13 || d.F8 < "1.28" {
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
				if v.F12.(string)[:3] == "688" || v.F12.(string)[:2] == "ST" || v.F12.(string)[:3] == "*ST" {
					continue
				}
				//f62 := decimal.NewFromFloat(v.F62.(float64))
				if reflect.TypeOf(v.F3).Name() == "string" {
					continue
				}
				if reflect.TypeOf(v.F2).Name() == "string" {
					continue
				}
				if reflect.TypeOf(v.F62).Name() == "string" {
					continue
				}
				if v.F3.(float64) < 0 || v.F3.(float64) > 7 || v.F2.(float64) > 68 || v.F62.(float64) < 0 {
					continue
				}
				d := stocks_db.NewStock_Day_K().GetStockDayKJJ(v.F12.(string))
				if (d.DayK30 > d.DayK20 && d.DayK20 > d.DayK10) || d.Day5Zdf > 8 || d.Day5Zdf < 1 || d.Day20Zdf < 0 || d.Day20Zdf > 13 || d.F8 < "1.28" {
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
