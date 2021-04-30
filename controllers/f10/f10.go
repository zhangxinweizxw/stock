package f10

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"stock/models/stocks_db"
	"stock/share/logging"
)

// 财务分析报告期
type FinancialReports struct {
	Date        string `json:"date"`        // 发布日期
	Jbmgsy      string `json:"jbmgsy"`      // 基本每股收益
	Mgjzc       string `json:"mgjzc"`       // 每股净资产(元)
	Mgwfply     string `json:"mgwfply"`     // 每股未分配利润(元)
	Yyzsr       string `json:"yyzsr"`       // 营业总收入(元)
	Gsjlr       string `json:"gsjlr"`       // 归属净利润(元)
	Kfjlr       string `json:"kfjlr"`       // 扣非净利润(元)
	Yyzsrtbzz   string `json:"yyzsrtbzz"`   // 营业总收入同比增长(%) --
	Gsjlrtbzz   string `json:"gsjlrtbzz"`   // 归属净利润同比增长(%)
	Kfjlrtbzz   string `json:"kfjlrtbzz"`   // 归属净利润同比增长(%)
	Yyzsrgdhbzz string `json:"yyzsrgdhbzz"` // 营业总收入滚动环比增长(%)--
	Gsjlrgdhbzz string `json:"gsjlrgdhbzz"` // 归属净利润滚动环比增长(%)
}

func NewFinancialReports() *FinancialReports {
	return &FinancialReports{}
}

// 保存F10 财务分析数据
func (this *FinancialReports) SaveFinaRepo(c *gin.Context) {

	// 先清空表数据
	stocks_db.NewFinancialReports().DelFinancialReports()
	// 查询最新个股
	sl := stocks_db.NewStockInfo().GetStockInfoList()
	if len(sl) <= 0 {
		logging.Error("Stock_info 表为空！")
		return
	}

	for _, v := range sl {
		sc := ""
		switch v.StockCode[:3] {
		case "600", "601", "603", "605", "688", "689", "608":
			sc = fmt.Sprintf("SH%v", v.StockCode)
		case "300", "002", "000", "001", "003":
			sc = fmt.Sprintf("SZ%v", v.StockCode)
		}
		fd := this.GetDFCFF10(sc) // 获取到的数据插入mysql
		for _, s := range fd {
			i := stocks_db.NewFinancialReports()
			p := map[string]interface{}{
				"stock_code": v.StockCode,
				"date":       s.Date,
				"jbmgsy":     s.Jbmgsy,
				"mgjzc":      s.Mgjzc,
				"mgwfply":    s.Mgwfply,
				"yyzsr":      s.Yyzsr,
				"gsjlr":      s.Gsjlr,
				"yyzsrtbzz":  s.Yyzsrtbzz,
				"gsjlrtbzz":  s.Gsjlrtbzz,
				"kfjlrtbzz":  s.Kfjlrtbzz,
			}
			_, err1 := i.Insert(p)
			if err1 != nil {
				logging.Error("Insert Table Financial Reports | %v", err1)
				break
			}
		}
	}
}

func (this *FinancialReports) GetDFCFF10(sc string) []*FinancialReports {
	// 查询东财F10
	url := fmt.Sprintf("http://f10.eastmoney.com/NewFinanceAnalysis/MainTargetAjax?type=0&code=%v", sc)
	resp, err := http.Get(url)
	if err != nil {
		logging.Error("F10", err)
	}
	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		logging.Error("ioutil.ReadAll", err1)
	}

	var data []*FinancialReports
	if err = json.Unmarshal(body, &data); err != nil {
		logging.Error("解析F10 财报分析  | Error:=", err)
	}
	return data
}

////更新季报 半年报 年报（手动调接口执行）
//func (this *F10) UpdateAnnualReport(c *gin.Context) {
//	// 更新mysql和redis
//
//	redisli, err := redis.Keys("stock:f10:*")
//	if err != nil {
//		fmt.Println("Delete Redis stock:f10:* Error   |%v", err)
//		return
//	}
//	for _, v := range redisli {
//		redis.Del(v)
//	}
//
//	t := fcmysql.NewF10()
//	t.DelTableF10()
//
//	// 查询数据库中所有个股信息 需要排重
//	stockList, err := fcmysql.NewStock_Day_K().GetStockInfoList()
//	if err != nil {
//		lib.WriteString(c, 40001, fmt.Sprintf("Select DB Error | %v", err))
//		return
//	}
//
//	fmt.Println(time.Now(), "     | F10 更新redis和mysql开始")
//	for _, v := range stockList {
//		stockCode := ""
//		switch v.F12[:3] {
//		case "600":
//			stockCode = fmt.Sprintf("SH%v", v.F12)
//		case "603":
//			stockCode = fmt.Sprintf("SH%v", v.F12)
//		case "601":
//			stockCode = fmt.Sprintf("SH%v", v.F12)
//		case "688":
//			stockCode = fmt.Sprintf("SH%v", v.F12)
//		case "300":
//			stockCode = fmt.Sprintf("SZ%v", v.F12)
//		case "002":
//			stockCode = fmt.Sprintf("SZ%v", v.F12)
//		case "000":
//			stockCode = fmt.Sprintf("SZ%v", v.F12)
//		case "001":
//			stockCode = fmt.Sprintf("SZ%v", v.F12)
//		case "003":
//			stockCode = fmt.Sprintf("SZ%v", v.F12)
//		}
//		stockIndexDayKUrl := fmt.Sprintf("http://f10.eastmoney.com/NewFinanceAnalysis/MainTargetAjax?type=0&code=%v", stockCode)
//		resp, err := http.Get(stockIndexDayKUrl)
//		if err != nil {
//			lib.WriteString(c, 40001, fmt.Sprintf("Get DFCF Http Error | %v", err))
//		}
//		defer resp.Body.Close()
//
//		body, err := ioutil.ReadAll(resp.Body)
//		if err != nil {
//			lib.WriteString(c, 40001, fmt.Sprintf("Io ReadAll Error | %v", err))
//		}
//		//fmt.Println("====", string(body))
//
//		var data []*F10
//		if err = json.Unmarshal(body, &data); err != nil {
//
//			fmt.Println("解析F10  Error:=   %v      |stockCode=|  %v", err, v.F12)
//			lib.WriteString(c, 40001, fmt.Sprintf("Json Unmarshal F10 Error | %v", err))
//			return
//		}
//		// 存储redis和mysql
//
//		this.SaveF10(data, stockCode)
//
//	}
//	fmt.Println(time.Now(), "     | F10 更新redis和mysql结束")
//	lib.WriteString(c, 200, "Success")
//}
//
//// 把F10数据存储到redis
//func (this *F10) SaveF10(data []*F10, stockCode string) {
//	redisKey := fmt.Sprintf("stock:f10:%v", stockCode)
//	datas, err := json.Marshal(data)
//	if err != nil {
//		fmt.Println("Redis Set Json marshal Error | %v", err)
//	}
//	rerr := redis.Set(redisKey, datas)
//	if rerr != nil {
//		fmt.Println("Redis Set  Error | %v", err)
//	}
//	//  同时插入到数据库
//	for _, v := range data {
//		t := fcmysql.NewF10()
//		params := map[string]string{
//			"f12":         stockCode[2:],
//			"date":        v.Date,
//			"jbmgsy":      v.Jbmgsy,
//			"mgjzc":       v.Mgjzc,
//			"mggjj":       v.Mggjj,
//			"mgwfply":     v.Mgwfply,
//			"mgjyxjl":     v.Mgjyxjl,
//			"yyzsr":       v.Yyzsr,
//			"mlr":         v.Mlr,
//			"gsjlr":       v.Gsjlr,
//			"kfjlr":       v.Kfjlr,
//			"yyzsrtbzz":   v.Yyzsrtbzz,
//			"gsjlrtbzz":   v.Gsjlrtbzz,
//			"kfjlrtbzz":   v.Kfjlrtbzz,
//			"yyzsrgdhbzz": v.Yyzsrgdhbzz,
//			"gsjlrgdhbzz": v.Gsjlrgdhbzz,
//			"kfjlrgdhbzz": v.Kfjlrgdhbzz,
//			"tbjzcsyl":    v.Tbjzcsyl,
//			"tbzzcsyl":    v.Tbzzcsyl,
//			"mll":         v.Mll,
//			"jll":         v.Jll,
//			"create_time": time.Now().Unix(),
//		}
//		_, err := t.Insert(params)
//		if err != nil {
//			fmt.Println("Insert Stock_F10 | %v", err)
//			return
//		}
//	}
//}
//
//// 查询 F10 返回F10 打分 --- 暂时不用了
//func (this *F10) F10Analysis(stockCode string) int {
//	var f10Data []F10
//	f10Info, _ := redis.Get(fmt.Sprintf("stock:f10:%v", stockCode))
//	err := json.Unmarshal([]byte(f10Info), &f10Data)
//	if err != nil {
//		logging.Error("Json Unmarshal Error | %v", err)
//		return 0
//	}
//	//-----------------------------------------
//	initialScore := 0
//	if len(f10Data) > 0 {
//
//		f10 := f10Data[0]
//
//		if this.interfaceToStringToFloat32(f10, f10.Gsjlrgdhbzz) < 0 || this.interfaceToStringToFloat32(f10, f10.Tbjzcsyl) < 0 { // 如果返回负数 直接跳过
//			return -1
//		}
//		// 每股未分配利润
//		mgwfply := this.interfaceToStringToFloat32(f10, f10.Mgwfply)
//		if mgwfply >= 1 && mgwfply <= 3 { // +3
//			initialScore += 3
//		} else if mgwfply >= 3 && mgwfply <= 6 { // +5
//			initialScore += 6
//		} else if mgwfply > 6 { // +10
//			initialScore += 10
//		}
//		// 每股经营现金流
//		mgjyxjl := this.interfaceToStringToFloat32(f10, f10.Mgjyxjl)
//		if mgjyxjl >= 1 && mgjyxjl <= 3 { // +3
//			initialScore += 3
//		} else if mgjyxjl >= 3 && mgjyxjl <= 6 { // +5
//			initialScore += 6
//		} else if mgjyxjl > 6 { // +10
//			initialScore += 10
//		}
//		// 营业总收入同比增长
//		yyzsrtbzz := this.interfaceToStringToFloat32(f10, f10.Yyzsrtbzz)
//		if yyzsrtbzz >= 5 && yyzsrtbzz <= 15 { // +1
//			initialScore += 1
//		} else if yyzsrtbzz >= 16 && yyzsrtbzz <= 30 { // +3
//			initialScore += 3
//		} else if yyzsrtbzz >= 31 && yyzsrtbzz <= 60 { // +6
//			initialScore += 6
//		} else if yyzsrtbzz > 60 { // +10
//			initialScore += 10
//		}
//		// 营业总收入滚动环比增长
//		yyzsrgdhbzz := this.interfaceToStringToFloat32(f10, f10.Yyzsrgdhbzz)
//		if yyzsrgdhbzz >= 0 && yyzsrgdhbzz <= 1 { // +1
//			initialScore += 1
//		} else if yyzsrgdhbzz > 1 && yyzsrgdhbzz <= 2 { // +3
//			initialScore += 3
//		} else if yyzsrgdhbzz > 2 { // +5
//			initialScore += 5
//		}
//
//		// 归属净利润滚动环比增长
//		gsjlrgdhbzz := this.interfaceToStringToFloat32(f10, f10.Gsjlrgdhbzz)
//		if gsjlrgdhbzz >= 0 && gsjlrgdhbzz <= 10 { // +1
//			initialScore += 1
//		} else if gsjlrgdhbzz > 10 && gsjlrgdhbzz <= 20 { // +3
//			initialScore += 3
//		} else if gsjlrgdhbzz > 20 { // +5
//			initialScore += 5
//		}
//		// 摊薄资产收益率
//		tbjzcsyl := this.interfaceToStringToFloat32(f10, f10.Tbjzcsyl)
//		if tbjzcsyl >= 0 && tbjzcsyl <= 1 { // +1
//			initialScore += 1
//		} else if tbjzcsyl > 1 && tbjzcsyl <= 2 { // +3
//			initialScore += 3
//		} else if tbjzcsyl > 2 { // +5
//			initialScore += 5
//		}
//		// 净利率
//		jll := this.interfaceToStringToFloat32(f10, f10.Jll)
//		if jll >= 0 && jll <= 10 { // +1
//			initialScore += 1
//		} else if jll > 10 && jll <= 20 { // +3
//			initialScore += 3
//		} else if jll > 20 { // +5
//			initialScore += 5
//		}
//	}
//
//	return initialScore
//}
//
//// F10 数据 string -> string -> float32
//func (this *F10) interfaceToStringToFloat32(f10 F10, intf string) float64 {
//
//	float01 := f10.Gsjlrgdhbzz.(string)
//	float02, _ := strconv.ParseFloat(float01, 32)
//	return float02
//}
//
//// 查询出的数据根据F10过滤
//func (this *F10) AnalysisF10List(stockCL []string) ([]string, error) {
//
//	dateL, err := fcmysql.NewF10().GetDate4Str() // 最近四期发布时间
//	if err != nil {
//		logging.Error("Get Date 4 Error | %v", err)
//		return nil, err
//	}
//	// stockCL 拼接成  'XXX','XXX'格式
//	var stockCodeStr = ""
//	for inx, v := range stockCL {
//
//		if len(stockCL) == 1 {
//			stockCodeStr += fmt.Sprintf("'%v'", v)
//			break
//		}
//		if inx == len(stockCL)-1 {
//			stockCodeStr += fmt.Sprintf("'%v'", v)
//		} else {
//			stockCodeStr += fmt.Sprintf("'%v',", v)
//		}
//	}
//	return fcmysql.NewF10().GetAnalysisF10(dateL, stockCodeStr)
//
//}
