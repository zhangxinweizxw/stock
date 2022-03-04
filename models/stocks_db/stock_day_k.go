package stocks_db

import (
	"fmt"
	"stock/share/logging"
	. "stock/share/models"
	"stock/share/util"
	"time"
)

const (
	TABLE_STOCK_DAY_K       = "stock_day_k"   // 日k线信息
	TABLE_STOCK_INDEX_DAY_K = "stock_index_k" // 日k线信息
)

type Stock struct {
	F1  float32 `json:"f1"`
	F2  float32 `json:"f2"`  // 最新价
	F3  float32 `json:"f3"`  // 涨跌幅
	F4  float32 `json:"f4"`  // 涨跌额
	F5  float32 `json:"f5"`  // 成交量(手)
	F6  float32 `json:"f6"`  // 成交额
	F7  float32 `json:"f7"`  // 振幅
	F8  float32 `json:"f8"`  // 换手率
	F9  float32 `json:"f9"`  // 市盈率(动态)
	F10 float32 `json:"f10"` // 量比
	F11 float32 `json:"f11"` // 5分钟涨跌
	F12 string  `json:"f12"` // 代码
	F13 float32 `json:"f13"`
	F14 string  `json:"f14"` // 名称
	F15 float32 `json:"f15"` // 最高
	F16 float32 `json:"f16"` // 最低
	F17 float32 `json:"f17"` // 今开
	F18 float32 `json:"f18"` // 昨收
	//F19 float32 `json:"f19"`
	F20  float32 `json:"f20"` // 总市值
	F21  float32 `json:"f21"` // 流通市值
	F22  float32 `json:"f22"` // 涨速
	F23  float32 `json:"f23"` // 市净率
	F24  float32 `json:"f24"` // 60日涨跌幅
	F25  float32 `json:"f25"` // 年初至今涨跌幅
	F62  float32 `json:"f62"` // 主力净流入
	F115 float32 `json:"f115"`
	F128 string  `json:"f128"` // 领涨股
	F140 string  `json:"f140"`
	F141 string  `json:"f141"`
	F136 string  `json:"f136"` // 涨跌幅
	F152 float32 `json:"f152"`
}

type Data struct {
	Total float32 `json:"total"`
	Diff  []Stock `json:"diff"`
}

type Dfcf struct {
	Data01 Data `json:"data"`
}

type Stock_Day_K struct {
	Model `db:"-" `
	F1    string
	F2    string `db:"f2"` // 最新价
	F3    string // 涨跌幅
	F4    string // 涨跌额
	F5    string // 成交量(手)
	F6    string `db:"f6"` // 成交额
	F7    string // 振幅
	F8    string `db:"f8"` // 换手率
	F9    string // 市盈率(动态)
	F10   string `db:"f10"` // 量比
	F11   string // 5分钟涨跌
	F12   string `db:"f12"` // 代码
	F13   string
	F14   string `db:"f14"` // 名称
	F15   string // 最高
	F16   string // 最低
	F17   string `db:"f17"` // 今开
	F18   string `db:"f18"` // 昨收
	//F19 float32 `json:"f19"`
	F20      string // 总市值
	F21      string // 流通市值
	F22      string // 涨速
	F23      string // 市净率
	F24      string // 60日涨跌幅
	F25      string // 年初至今涨跌幅
	F62      string // 主力净流入
	F115     string
	F128     string // 领涨股
	F140     string
	F141     string
	F136     string // 涨跌幅
	F152     string
	DayK5    float64 `db:"dayK5"`
	DayK10   float64 `db:"dayK10"`
	DayK20   float64 `db:"dayK20"`
	DayK30   float64 `db:"dayK30"`
	DayK60   float64 `db:"dayK60"`
	Day5Zdf  float64 `db:"day5zdf"`
	Day10Zdf float64 `db:"day10zdf"`
	Day20Zdf float64 `db:"day20zdf"`
}

func NewStock_Day_K() *Stock_Day_K {
	return &Stock_Day_K{
		Model: Model{
			TableName: TABLE_STOCK_DAY_K,
			Db:        MyCat,
		},
	}
}

func (this *Stock_Day_K) DelStockDayK() {

	b := this.Db.DeleteFrom(this.TableName).
		Where(fmt.Sprintf("create_time='%v'", time.Now().Format("2006-01-02")))
	_, err := this.DeleteWhere(b, nil).Exec()

	if err != nil {
		fmt.Println("Delete Table TABLE_STOCK_Day_K  |  Error   %v", err)
	}

}

func (this *Stock_Day_K) GetXQStockList() []*XQ_Stock {

	// 查询最新日期
	ctime := ""
	bulid := this.Db.Select("create_time").From(XQ_STOCK).
		OrderBy("create_time DESC").Limit(1)
	_, err := this.SelectWhere(bulid, nil).LoadStructs(&ctime)
	if err != nil {
		fmt.Println("Select Table xq_stock  |  Error   %v", err)
		return nil
	}
	var xqStock []*XQ_Stock
	bulid1 := this.Db.Select("*").From(XQ_STOCK).
		Where(fmt.Sprintf("create_time='%v'", ctime))

	_, err1 := this.SelectWhere(bulid1, nil).LoadStructs(&xqStock)
	if err1 != nil {
		fmt.Println("Select Table xq_stock  |  Error   %v", err1)
		return nil
	}
	return xqStock
}

// 查询日K 雪球筛选是否执行
func (this *Stock_Day_K) GetIsZx() int {

	t := time.Now().Format("2006-01-02")

	ct := 0
	bulid := this.Db.Select("count(1) ct").From(TABLE_STOCK_DAY_K).
		Where(fmt.Sprintf("create_time='%v'", t))
	_, err := this.SelectWhere(bulid, nil).LoadStructs(&ct)
	if err != nil {
		fmt.Println("Select Table TABLE_STOCK_DAY_K  |  Error   %v", err)
		return ct
	}
	if ct < 4000 && ct != 0 {
		// 大概率日K异常报警
		go util.NewDdRobot().DdRobotPush("日K更新异常！")
	}

	return ct
}

// 短线 3、4、5、7、8 、11、13、18、21 选股法
func (this *Stock_Day_K) GetDxStockDayKList(sql string) []*Stock_Day_K {

	var sdkl []*Stock_Day_K

	//sql := `SELECT f12 FROM stock_day_k
	//		WHERE f8>3  AND f9 <20 AND f62 >1000000 AND f24<20
	//		AND FROM_UNIXTIME(create_time,'%Y-%m-%d') ='` + dateStr[0]
	//sql += `' AND f12 IN (
	//			SELECT f12 FROM stock_day_k
	//			WHERE f8>2  AND f9 <20 AND f62 >1000000 AND f24<20
	//			AND FROM_UNIXTIME(create_time,'%Y-%m-%d') ='` + dateStr[1]
	//sql += `' AND f12 IN (
	//				SELECT f12 FROM stock_day_k
	//				WHERE f8>1  AND f9 <30 AND f62 >1000000 AND f24<20
	//				AND FROM_UNIXTIME(create_time,'%Y-%m-%d') ='` + dateStr[0]
	//sql += `'
	//			)
	//		)`
	_, err := this.Db.SelectBySql(sql).
		LoadStructs(&sdkl)
	if err != nil {
		logging.Error("Select Stock_day_k Error：%v", err)
		return nil
	}

	return sdkl
}

// 获取日k表中最近23个交易人日期
func (this *Stock_Day_K) GetStockDayKDate() []string {
	// 查询最新日期
	var cm []string
	bulid := this.Db.Select("create_time").
		From(this.TableName).
		GroupBy("create_time").
		OrderBy("create_time DESC").
		Limit(23)
	_, err := this.SelectWhere(bulid, nil).LoadStructs(&cm)
	if err != nil {
		fmt.Println("Select Table stock_day_k create_time  |  Error   %v", err)
		return nil
	}
	return cm
}

// 获取30日K均线价位
func (this *Stock_Day_K) GetStockDayK30Date(sc string) float64 {
	// 查询最新日期
	var dk30 float64
	bulid := this.Db.Select("dayK30 as dk30").
		From(this.TableName).
		Where(fmt.Sprintf("f12='%v'", sc)).
		OrderBy("create_time DESC").
		Limit(1)
	_, err := this.SelectWhere(bulid, nil).LoadStructs(&dk30)
	if err != nil {
		fmt.Println("Select Table stock_day_k create_time  |  Error   %v", err)
		return 0
	}
	return dk30
}

// 获取10日K均线价位
func (this *Stock_Day_K) GetStockDayK10Date(sc string) float64 {
	// 查询最新日期
	var dk10 float64
	bulid := this.Db.Select("dayK10 as dk10").
		From(this.TableName).
		Where(fmt.Sprintf("f12='%v'", sc)).
		OrderBy("create_time DESC").
		Limit(1)
	_, err := this.SelectWhere(bulid, nil).LoadStructs(&dk10)
	if err != nil {
		fmt.Println("Select Table stock_day_k create_time  |  Error   %v", err)
		return 0
	}
	return dk10
}

//抓涨停测试使用
func (this *Stock_Day_K) ZtSelStockDayk() []*Stock_Day_K {
	var sdkl []*Stock_Day_K

	sql := `SELECT f12,f14,dayK5,dayK10,dayK20,dayK30 FROM stock_day_k
			where f14 NOT LIKE 'N%' AND f14 NOT LIKE '*%' 
			AND f14 NOT LIKE 'ST%' AND f12 NOT LIKE '688%'
			AND f2 <58 AND f2 >3.98 AND f23 <6 AND f23 >1
			AND f21 >2500000000 AND f21 < 15000000000
			AND f3 < 1.8 AND f3 > -1.28 AND f10 < 1.28
			AND create_time=(SELECT create_time FROM stock_day_k ORDER BY create_time DESC LIMIT 1)`
	_, err := this.Db.SelectBySql(sql).
		LoadStructs(&sdkl)
	if err != nil {
		logging.Error("Select Stock_day_k -> zt_stock Error：%v", err)
		return nil
	}

	return sdkl
}

// 获取5、10、20、30日K均线价位
func (this *Stock_Day_K) GetStockDayKJJ(sc string) *Stock_Day_K {
	// 查询最新日期
	var d *Stock_Day_K
	bulid := this.Db.Select("dayK5,dayK10,dayK20,dayK30,day5zdf,day10zdf,day20zdf,f8,f10").
		From(this.TableName).
		Where(fmt.Sprintf("f12='%v'", sc)).
		OrderBy("create_time DESC").
		Limit(1)
	_, err := this.SelectWhere(bulid, nil).LoadStructs(&d)
	if err != nil {
		fmt.Println("Select Table stock_day_k create_time  |  Error   %v", err)
		return d
	}
	return d
}

// 筛选上个交易日
func (this *Stock_Day_K) GetDaykJC(sc string) string {

	f12 := ""
	bulid := this.Db.Select("f12").
		From(this.TableName).
		Where(fmt.Sprintf("f12='%v'", sc)).
		Where(" dayK30 > dayK60 AND dayK20 > dayK30 AND day5zdf < 8  AND f8 >0.8 AND f8< 8 AND f10 >0.5 AND f10 <8 AND f3 >0 AND f3 <3 AND f2 < 58 AND f9 <28 AND f9 >5 AND f62 >10000000 AND f2 >f17 AND dayK10 > dayK20 AND dayK5 < dayK10").
		Where("create_time= (SELECT date FROM stock_info ORDER BY date DESC  LIMIT 1)")
	_, err := this.SelectWhere(bulid, nil).LoadStructs(&f12)
	if err != nil {
		fmt.Println("Select Table 上个交易日 create_time  |  Error   %v", err)
		return ""
	}
	return f12
}

//// 返回 今日cjje和前面平均五日成交金额
//func (this *Stock_Day_K) ReturnIsBuy(sc string) (avg5cjje float64, sdkInfo *Stock_Day_K) {
//
//	//t := time.Now().Format("2006-01-02")
//	l := this.GetStockDayKDate()
//	cjje := 0.0
//	//t = "2021-12-07"
//	sql := fmt.Sprintf(`SELECT AVG(a.f6) FROM(
//				SELECT f6,create_time FROM stock_day_k
//				WHERE f12='%v'
//				AND create_time < '%v'
//				ORDER BY create_time DESC
//				LIMIT 18
//				)a`, sc, l[0])
//	_, err := this.Db.SelectBySql(sql).
//		LoadStructs(&cjje)
//	if err != nil {
//		logging.Error("Select Stock_day_k avg(f6) 5day Error：%v", err)
//		return cjje, nil
//	}
//
//	//查询当日日K数据
//	var d *Stock_Day_K
//	bulid := this.Db.Select("f6,f17,f18,f2,create_time").
//		From(this.TableName).
//		Where(fmt.Sprintf("f12='%v' AND create_time = '%v'", sc, l[0]))
//	_, err1 := this.SelectWhere(bulid, nil).LoadStructs(&d)
//	if err != nil {
//		fmt.Println("Select Table stock_day_k   |  Error   %v", err1)
//		return cjje, d
//	}
//	return cjje, d
//}
//
//// 返回 今日cjje和前面平均五日成交金额  zt 专用
//func (this *Stock_Day_K) ReturnIsBuyZt(sc string) (avg5cjje float64, sdkInfo *Stock_Day_K) {
//
//	l := this.GetStockDayKDate()
//
//	cjje := 0.0
//	sql := fmt.Sprintf(`SELECT AVG(a.f6) FROM(
//				SELECT f6,create_time FROM stock_day_k
//				WHERE f12='%v'
//				AND create_time < '%v'
//				ORDER BY create_time DESC
//				LIMIT 8
//				)a`, sc, l[0])
//	_, err := this.Db.SelectBySql(sql).
//		LoadStructs(&cjje)
//	if err != nil {
//		logging.Error("Select Stock_day_k avg(f6) 5day Error：%v", err)
//		return cjje, nil
//	}
//
//	//查询当日日K数据
//	var d *Stock_Day_K
//	bulid := this.Db.Select("f6,f17,f18,f2,create_time").
//		From(this.TableName).
//		Where(fmt.Sprintf("f12='%v' AND create_time = '%v'", sc, l[0]))
//	_, err1 := this.SelectWhere(bulid, nil).LoadStructs(&d)
//	if err != nil {
//		fmt.Println("Select Table stock_day_k   |  Error   %v", err1)
//		return cjje, d
//	}
//	return cjje, d
//}
