package util

type StockDayK struct {
	Datas *Datar `json:"data"`
}

type Datar struct {
	Total float64      `json:"total"`
	Diff  []*StockInfo `json:"diff"`
}

type StockInfo struct {
	F1  interface{} `json:"f1"`
	F2  interface{} `json:"f2"`  // 最新价
	F3  interface{} `json:"f3"`  // 涨跌幅
	F4  interface{} `json:"f4"`  // 涨跌额
	F5  interface{} `json:"f5"`  // 成交量(手)
	F6  interface{} `json:"f6"`  // 成交额
	F7  interface{} `json:"f7"`  // 振幅
	F8  interface{} `json:"f8"`  // 换手率
	F9  interface{} `json:"f9"`  // 市盈率(动态)
	F10 interface{} `json:"f10"` // 量比
	F11 interface{} `json:"f11"` // 5分钟涨跌
	F12 interface{} `json:"f12"` // 代码
	F13 interface{} `json:"f13"`
	F14 interface{} `json:"f14"` // 名称
	F15 interface{} `json:"f15"` // 最高
	F16 interface{} `json:"f16"` // 最低
	F17 interface{} `json:"f17"` // 今开
	F18 interface{} `json:"f18"` // 昨收
	//F19 float32 `json:"f19"`
	F20  interface{} `json:"f20"` // 总市值
	F21  interface{} `json:"f21"` // 流通市值
	F22  interface{} `json:"f22"` // 涨速
	F23  interface{} `json:"f23"` // 市净率
	F24  interface{} `json:"f24"` // 60日涨跌幅
	F25  interface{} `json:"f25"` // 年初至今涨跌幅
	F62  interface{} `json:"f62"` // 主力净流入
	F115 interface{} `json:"f115"`

	F128 interface{} `json:"f128"` // 领涨股
	F140 interface{} `json:"f140"`
	F141 interface{} `json:"f141"`
	F136 interface{} `json:"f136"` // 涨跌幅
	F152 interface{} `json:"f152"`

	F186 interface{} `json:"f186"` // AH比价最新价
	F187 interface{} `json:"f187"` // AH比价涨跌幅
	F191 string      `json:"f191"` // AH比价股票代码
	F193 string      `json:"f193"` // AH比价股票名称

	F127 interface{} `json:"f127"` // 3日涨跌幅
	F268 float64     `json:"f268"` // 净占比

	F184 interface{} `json:"f184"`
	F66  interface{} `json:"f66"`
	F72  interface{} `json:"f72"`
	F69  interface{} `json:"f69"`
	F172 interface{} `json:"f172"`
	F75  interface{} `json:"f75"`
	F178 interface{} `json:"f178"`
	F81  interface{} `json:"f81"`
}

//f62  今日主力净流入        主力净比  f184
//f66  今日超大单净流入      超大单净比 f69
//f72  今日大单净流入        大单净比  f75
//f78  今日中单净流入        中单净比  f81

type XQResult struct {
	XQResuData XQData `json:"data"`
}
type XQData struct {
	Count int       `json:"count"`
	List  []*XQList `json:"list"`
}
type XQList struct {
	StockName string `json:"name"`
	StockCode string `json:"symbol"`
}
