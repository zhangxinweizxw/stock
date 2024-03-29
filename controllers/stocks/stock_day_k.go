package stocks

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"stock/config"
	"stock/controllers"
	. "stock/models"
	"stock/models/stocks_db"
	"stock/share/logging"
	"stock/share/util"
	"strconv"
	"strings"
	"time"
)

type Date struct {
	StockDate *StockDayk `json:"data"`
}

//f43     最新价格
//f44  		最高价格
//f45			最低价格
//f47			成交量
//f48			成交额
//f50     量比
//f51     涨停价
//f52     跌停价
//f57     代码
//f58			名称
//f168    换手率
//f169    涨跌额
//f170    涨跌幅

type StockDayk struct {
	C     *config.AppConfig
	Zxjg  interface{} `json:"f43"`
	Zgjg  float64     `json:"f44"`
	Zdjg  float64     `json:"f45"`
	Cjj   float64     `json:"f47"`
	Cje   float64     `json:"f48"`
	Lb    float64     `json:"f50"`
	Ztj   interface{} `json:"f51"`
	Dtj   interface{} `json:"f52"`
	Gpdm  string      `json:"f57"`
	Gpmc  string      `json:"f58"`
	Hsl   float64     `json:"f168"`
	Zde   float64     `json:"f169"`
	Zdf   float64     `json:"f170"`
	Zljlr interface{} `json:"f137"`
	Jcd   interface{} `json:"f140"`
	Jdd   interface{} `json:"f143"`
	//Jzd   float64     `json:"f146"`
	Kpj float64 `json:"f46"`
	Zsz float64 `json:"f116"`
}

func NewStockDayk(cfg *config.AppConfig) *StockDayk {
	return &StockDayk{
		C: cfg,
	}
}

// 所有股票日K
func (this *StockDayk) GetStockDayK() {
	param := "?pn=1&pz=5000&po=1&np=1&fltt=2&invt=2&fid=f3&fs=m:0+t:6,m:0+t:13,m:0+t:80,m:1+t:2,m:1+t:23&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152"
	url := fmt.Sprintf("%s%s", this.C.Url.DfcfStockDayK, param)
	resp, err := http.Get(url)

	if err != nil {
		logging.Error("", err)
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

	stocks_db.NewStockInfo().DelStockInfo() // 清空stock_info 表
	stocks_db.NewStock_Day_K().DelStockDayK()

	i := int(len(data.Datas.Diff) / 2)
	go this.GoFuncFor(data, 0, i)
	go this.GoFuncFor(data, i, len(data.Datas.Diff)-1)
}

func (this *StockDayk) GoFuncFor(data *util.StockDayK, s, e int) {
	// 日K行情写入mysql
	ntime := time.Now().Format("2006-01-02")
	for i, v := range data.Datas.Diff {
		if i >= s && i < e {

			//获取资金流入数据
			if len(v.F12.(string)) != 6 {
				continue
			}
			d := this.GetZJlxDFCF(v.F12.(string)).Datas.Diff[0]

			//日K对应 5日 10日 20日 30日均价
			df := this.GetDayK(v.F12.(string))

			t := stocks_db.NewStock_Day_K()
			params := map[string]interface{}{
				"f1":  v.F1,
				"f2":  v.F2,
				"f3":  v.F3,
				"f4":  v.F4,
				"f5":  v.F5,
				"f6":  v.F6,
				"f7":  v.F7,
				"f8":  v.F8,
				"f9":  v.F9,
				"f10": v.F10,
				"f11": v.F11,
				"f12": v.F12,
				"f13": v.F13,
				"f14": v.F14,
				"f15": v.F15,
				"f16": v.F16,
				"f17": v.F17,
				"f18": v.F18,

				"f20":    v.F20,
				"f21":    v.F21,
				"f22":    v.F22,
				"f23":    v.F23,
				"f24":    v.F24,
				"f25":    v.F25,
				"f62":    d.F62,
				"f66":    d.F66,
				"f69":    d.F69,
				"f72":    d.F72,
				"f75":    d.F75,
				"f184":   d.F184,
				"f136":   v.F136,
				"f128":   v.F128,
				"dayK5":  df[0],
				"dayK10": df[1],
				"dayK20": df[2],
				"dayK30": df[3],
				"dayK60": df[7],

				"day5zdf":     df[4],
				"day10zdf":    df[5],
				"day20zdf":    df[6],
				"create_time": time.Now().Format("2006-01-02"),
				"update_time": time.Now().Format("2006-01-02"),
			}
			_, err := t.Insert(params)
			if err != nil {
				logging.Error("Insert Stock_day_k | %v", err)
				continue
			}

			// 股票信息写入stock_info表方便使用
			i := stocks_db.NewStockInfo()
			p := map[string]interface{}{
				"date":       ntime,
				"stock_code": v.F12,
				"stock_name": v.F14,
			}
			_, err1 := i.Insert(p)
			if err1 != nil {
				logging.Error("Insert Stock_info | %v", err)
				continue
			}
		}
	}
}

// 雪球个股筛选判断
func (this *StockDayk) GetXueqiu() {
	stocks_db.NewXQ_Stock().DelXqStock()

	// 为了简单手动改报告期
	//url := "https://xueqiu.com/service/screener/screen?category=CN&exchange=sh_sz&areacode=&indcode=&order_by=symbol&order=desc&page=1&size=30&only_count=0&current=&pct=1.28_5.8&netprofit.20210630=50000000_61150000000&fmc=2500000000_15000000000&npay.20210630=5_17594.51&oiy.20210630=5_151223.7&volume_ratio=1.8_10&tr=3_10&pct5=0_8&pct20=-5_12"
	//url := "https://xueqiu.com/service/screener/screen?category=CN&exchange=sh_sz&areacode=&indcode=&order_by=symbol&order=desc&page=1&size=30&only_count=0&current=&pct=1.28_3.8&pettm=10_58&oiy.20210930=5_118173.56&npay.20210930=5_58997.1&follow7d=300_17574&tr=2.8_8&volume_ratio=1_9.48&pct5=-1.8_8&oiy.20210630=5_118173.56&npay.20210630=5_118173.56&oiy.20201231=5_118173.56&npay.20201231=5_118173.56"
	//url := "https://xueqiu.com/service/screener/screen?category=CN&exchange=sh_sz&areacode&indcode&order_by=symbol&order=desc&page=1&size=30&only_count=0&current=1.5_58&pct&pettm=5_58&pelyr=5_58&mc=2500000000_18000000000&oiy.20200630=5_2295.9&npay.20200630=5_142839.72&tr=1.28_8&volume_ratio=1.28_8&oiy.20200930=5_2295.9&npay.20200930=5_142839.72&pct5=-5_6"
	url := "https://xueqiu.com/service/screener/screen?category=CN&exchange=sh_sz&areacode&indcode&order_by=symbol&order=desc&page=1&size=30&only_count=0&current=1.8_58&pct=-0.8_2.8&volume_ratio=0.8_8&pettm=5_28&pct5=-3_5&oiy.20210930=5_118173.56&tr=1.28_8&pct20=-8_12&pct120=-38_-20&npay.20210930=5_58997.1&pct10=-5_8&oiy.20210630=5_118173.56&npay.20210630=5_58997.1"
	resp, err := http.Get(url)
	if err != nil {
		logging.Error("", err)
	}
	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		logging.Error("ioutil.ReadAll", err1)
	}
	//logging.Error("=====", string(body))
	var data *util.XQResult
	if err = json.Unmarshal(body, &data); err != nil {
		logging.Error("解析雪球筛选 | Error:=", err)
	}
	//logging.Error("=====", len(data.XQResuData.List))
	// 财务过滤一下

	// 写入mysql
	for _, v := range data.XQResuData.List {

		if len(stocks_db.NewStock_Day_K().GetDaykJC(v.StockCode[2:])) <= 0 {
			continue //基础较差
		}

		t := stocks_db.NewXQ_Stock()
		params := map[string]interface{}{
			"stock_code":  v.StockCode[2:],
			"stock_name":  v.StockName,
			"create_time": time.Now().Format("2006-01-02"),
		}
		_, err := t.Insert(params)
		if err != nil {
			logging.Error("Insert xq_stock | %v", err)
			return
		}
	}
	// 清空 缓存
	XQStock = nil

}

// 雪球最新业绩分析筛选
func (this *StockDayk) SaveXueqiuFx() {

	// 为了简单手动改报告期
	url := "https://xueqiu.com/service/screener/screen?category=CN&exchange=sh_sz&areacode=&indcode=&order_by=symbol&order=desc&page=1&size=30&only_count=0&current=&pct=0_3&pettm=5_60&pb=5_10&oiy.20210630=5_151223.7&npay.20210630=5_74261.79&pct120=-35_30&pct5=-5_8"
	resp, err := http.Get(url)
	if err != nil {
		logging.Error("", err)
	}
	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		logging.Error("ioutil.ReadAll", err1)
	}
	//logging.Error("=====", string(body))
	var data *util.XQResult
	if err = json.Unmarshal(body, &data); err != nil {
		logging.Error("解析雪球筛选 | Error:=", err)
	}
	//logging.Error("=====", len(data.XQResuData.List))

	// 操作mysql  不符合要求剔除 符合加入不清库
	for _, v := range data.XQResuData.List {

		//if NewStockDayk(nil).GetReturnIsBuyZt(v.StockCode[2:]) == false {
		//	continue
		//}

		t := stocks_db.NewXQ_Stock_FX()

		params := map[string]interface{}{
			"stock_code":  v.StockCode[2:],
			"stock_name":  v.StockName,
			"create_time": time.Now().Format("2006-01-02"),
		}
		_, err := t.Insert(params)
		if err != nil {
			logging.Error("Insert xq_stock | %v", err)
			return
		}
	}
	stocks_db.NewXQ_Stock_FX().DelStockFx()
	XqFxStockDb = nil
}

// 需求个股分析监控 9：15 - 11：30   13：00-15：00 XQ
func (this *StockDayk) XQStockFx() {
	//logging.Error("=======", len(XQStock))
	if len(XQStock) <= 0 {
		XQStock = stocks_db.NewStock_Day_K().GetXQStockList()
	}
	name := ""
	defer func() {
		if err := recover(); err != nil {
			logging.Error("Panic Error=======:%v======:%v", name, err)
		}
	}()

	for k, v := range XQStock {
		i := this.StockInfoSS(v.StockCode).StockDate

		if i == nil {
			continue
		}
		name = i.Gpmc
		zljlrv := 0.0
		if reflect.TypeOf(i.Zljlr).String() != "string" {
			zljlrv = i.Zljlr.(float64)
		}
		d1 := decimal.NewFromFloat(zljlrv)

		//d2 := "0"
		//if reflect.TypeOf(i.Jcd).String() != "string" {
		//	d2 = fmt.Sprintf("%v", i.Jcd.(float64))
		//}
		//
		////d3 := decimal.NewFromFloat(i.Jdd.(float64))
		//d3 := "0"
		//if reflect.TypeOf(i.Jdd).String() != "string" {
		//	d3 = fmt.Sprintf("%v", decimal.NewFromFloat(i.Jdd.(float64)))
		//}

		d101 := ""
		if i.Zsz < 3000000000 { // 市值30亿以内公司
			d101 = "1880000"
		}
		if i.Zsz > 3000000000 && i.Zsz < 5000000000 { //
			d101 = "2880000"
		}
		if i.Zsz > 5000000000 && i.Zsz < 15000000000 { //
			d101 = "3880000"
		}
		if i.Zsz > 15000000000 { //
			d101 = "5880000"
		}

		// 最高涨跌幅
		zgzdf := (i.Zgjg - i.Kpj) / i.Kpj

		zgzdfv, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", zgzdf*100), 64)

		if i.Zdf > -0.8 && i.Zdf < 3.8 && i.Lb > 0.5 && i.Lb < 8 && i.Hsl > 0.5 && i.Hsl < 10 && d1.String() > d101 && (zgzdfv-i.Zdf) < 2.8 {
			// 判断是否以入库
			sc := v.StockCode[2:]
			if stocks_db.NewTransactionHistory().GetTranHist(sc) > 0 {
				continue
			}

			// 满足条件从 List 中 去掉    mysql transaction_history 表中添加数据 // 发送叮叮实时消息

			if reflect.TypeOf(i.Zxjg).Name() == "string" {
				continue
			}
			go this.SaveStock(i.Gpdm, i.Gpmc, i.Zxjg.(float64), 1)
			XQStock = append(XQStock[:k], XQStock[k+1:]...)
			go util.NewDdRobot().DdRobotPush(fmt.Sprintf("建议买入：%v   |   股票代码：%v    买入价：%v", i.Gpmc, i.Gpdm, i.Zxjg))
		}
	}

}

// 实时拉去个股信息
func (this *StockDayk) StockInfoSS(sc string) *Date {
	defer func() {
		if err := recover(); err != nil {
			logging.Error("拉取 个股信息 Panic Error=============:%v", err)
		}
	}()
	if len(sc) < 2 {
		return nil
	}
	code := ""

	switch sc[:3] {
	case "600", "601", "603", "605", "688", "689", "608":
		code = fmt.Sprintf("1.%v", sc)
	case "300", "002", "000", "001", "003", "301":
		sc = fmt.Sprintf("0.%v", sc)
	default:
		return nil
	}

	url := fmt.Sprintf("http://push2.eastmoney.com/api/qt/stock/get?ut=fa5fd1943c7b386f172d6893dbfba10b&invt=2&fltt=2&fields=f43,f44,f45,f47,f48,f50,f51,f52,f57,f58,f168,f169,f170,f137,f140,f143,f146,f46,f116&secid=%v", code)
	//logging.Info("===============:", url)
	resp, err := http.Get(url)
	if err != nil {
		logging.Error("", err)
	}
	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		logging.Error("ioutil.ReadAll", err1)
	}

	var data *Date
	if err := json.Unmarshal(body, &data); err != nil {
		logging.Error("个股数据 | Error:=", err)
	}
	return data
}

// 满足选股条件的个股写入mysql
func (this *StockDayk) SaveStock(c, n string, p float64, s int) {

	t := stocks_db.NewTransactionHistory()
	params := map[string]interface{}{
		"stock_code": c,
		"stock_name": n,
		"buy_price":  p,
		"status":     s,
		"buy_time":   time.Now().Format("2006-01-02 15:04"),
	}
	_, err := t.Insert(params)
	if err != nil {
		logging.Error("Insert transaction_history | %v", err)
		return
	}
}

// 保存日K时查询 资金流入数据
func (this *StockDayk) GetZJlxDFCF(stockC string) *util.StockDayK {
	//stockCodes := ""
	//switch stockC[:3] {
	//case "600", "601", "603", "605", "688", "689", "608":
	//	stockCodes = fmt.Sprintf("1.%v", stockC)
	//case "300", "002", "000", "001", "003", "301":
	//	stockCodes = fmt.Sprintf("0.%v", stockC)
	//default:
	//	return nil
	//}
	stockCodes := controllers.NewUtilHttps(nil).GetUtilCode1(stockC)
	if len(stockCodes) <= 0 {
		return nil
	}
	if stockCodes[:2] != "1." && stockCodes[:2] != "0." {
		return nil
	}
	url := "http://push2.eastmoney.com/api/qt/ulist.np/get?fltt=2&secids=" + stockCodes + "&fields=f20%2Cf62%2Cf66%2Cf69%2Cf72%2Cf75%2Cf184%2Cf2%2Cf8%2Cf9%2Cf10%2Cf3"

	resp, err := http.Get(url)
	if err != nil {
		logging.Error("", err)
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
	return data
}

type Data struct {
	Zdata *Kl `json:"data"`
}
type Kl struct {
	Klines []string `json:"klines"`
}

// 返回日K对应的 5、10、20、30 均价
func (this *StockDayk) GetDayK(stockC string) [8]float64 {

	var dk = [8]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}

	stockCodes := controllers.NewUtilHttps(nil).GetUtilCode1(stockC)
	if len(stockCodes) <= 0 {
		return dk
	}

	if stockCodes[:2] != "1." && stockCodes[:2] != "0." {
		return dk
	}
	url := "http://push2his.eastmoney.com/api/qt/stock/kline/get?fields1=f1%2Cf2%2Cf3%2Cf4%2Cf5%2Cf6&fields2=f51%2Cf52%2Cf53%2Cf54%2Cf55%2Cf56%2Cf57%2Cf58%2Cf59%2Cf60%2Cf61&ut=7eea3edcaed734bea9cbfc24409ed989&klt=101&fqt=1&secid=" + stockCodes + "&beg=0&end=20500000"

	resp, err := http.Get(url)
	if err != nil {
		logging.Error("", err)
	}
	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		logging.Error("ioutil.ReadAll", err1)
	}

	var data *Data
	if err = json.Unmarshal(body, &data); err != nil {
		logging.Error("单个个股解析日K  | Error:=", err)
	}
	//logging.Error("==================", data.Zdata.Klines[0])
	if len(data.Zdata.Klines) < 60 {
		return dk
	}
	fl := len(data.Zdata.Klines) - 1

	for i := fl; i >= 0; i-- {
		s := strings.Split(data.Zdata.Klines[i], ",")
		if i >= fl-4 {
			f, _ := strconv.ParseFloat(s[2], 64)
			dk[0] += f
			r, _ := strconv.ParseFloat(s[8], 64)
			dk[4] += r
		}
		if i >= fl-9 {
			f, _ := strconv.ParseFloat(s[2], 64)
			dk[1] += f
			r, _ := strconv.ParseFloat(s[8], 64)
			dk[5] += r
		}
		if i >= fl-19 {
			f, _ := strconv.ParseFloat(s[2], 64)
			dk[2] += f
			r, _ := strconv.ParseFloat(s[8], 64)
			dk[6] += r
		}
		if i >= fl-29 {
			f, _ := strconv.ParseFloat(s[2], 64)
			dk[3] += f
		}
		if i >= fl-59 {
			f, _ := strconv.ParseFloat(s[2], 64)
			dk[7] += f
		}
		if i < fl-60 {
			break
		}

	}
	//logging.Error("==================", dk[2]/20, "==============", dk[3]/30)

	f0, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", dk[0]/5), 64)
	dk[0] = f0
	f1, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", dk[1]/10), 64)
	dk[1] = f1
	f2, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", dk[2]/20), 64)
	dk[2] = f2
	f3, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", dk[3]/30), 64)
	dk[3] = f3

	f7, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", dk[7]/60), 64)
	dk[7] = f7
	//logging.Error("==========", dk)
	return dk
}

func (this *StockDayk) GetStock7p() {

	stocks_db.NewZjlxStockDb().DelZjlxStock()

	s := stocks_db.NewStock_Day_K().Stock7p()
	//st := []string{"2022-02-24", "2022-02-25", "2022-02-28", "2022-03-01", "2022-03-02", "2022-03-03", "2022-03-04"}
	st := stocks_db.NewStock_Day_K().GetCreateTime()
	//logging.Debug("========", st)
	sort.Strings(st)
	//logging.Debug("========", st)
	//return
	for _, v := range s {
		zdj := v.F16 // 第一天的最低价
		k := 0
		if v.F12[:3] == "688" {
			continue
		}
		for i := 0; i < 7; i++ {
			//sv := stocks_db.NewStock_Day_K().Stock7ps(v.F12, v.CreateTime, zdj)
			isis := false
			if i == 6 {
				isis = true
			}
			sv := stocks_db.NewStock_Day_K().Stock7ps(v.F12, st[i], zdj, isis)
			if len(sv) <= 0 {
				break
			}
			k = k + 1
		}
		if k != 7 {
			continue
		}
		//logging.Debug("==========", v.F12)
		// 股票信息写入stock_info表方便使用
		i := stocks_db.NewZjlxStockDb()
		p := map[string]interface{}{
			"create_time": time.Now().Format("2006-01-02"),
			"stock_code":  v.F12,
			"stock_name":  v.F14,
		}
		_, err1 := i.Insert(p)
		if err1 != nil {
			logging.Error("Insert Table zjlx_stock | %v", err1)
			continue
		}
	}
	ZjlxStockDb = nil
}
