package util

import (
	"encoding/json"
	"io/ioutil"
	"stock/share/logging"
	"time"
)

/**
stock 公用方法
*/

type StockUtil struct {
}

func NewStockUtil() *StockUtil {
	return &StockUtil{}
}

type Data struct {
	Dates []*Date `json:"data"`
}
type Date struct {
	Zrxh int
	Jybz string `json:"jybz"` // 是否是交易日  0： 否  1：是
	Jyrq string `json:"jyrq"` // 日期
}

// 调用深交所 交易日历判断当天是否是交易日  返回true 是交易人
func (this *StockUtil) GetSjsMonthList() bool {
	//url := this.C.Url.SjsMonthList
	url := "http://www.szse.cn/api/report/exchange/onepersistenthour/monthList"

	err, resp := NewHttpUtil().GetJson(url)
	if err != nil {
		logging.Error("sjs", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.Error("ioutil.ReadAll  Error | %v", err)
	}

	var retBool = false
	var data *Data
	if err = json.Unmarshal(body, &data); err != nil {
		logging.Error("解析是否是交易日Url  | Error:=", err)
	}
	dtDate := time.Now().Format("2006-01-02")
	for _, v := range data.Dates {
		if v.Jyrq == dtDate {
			if v.Jybz == "1" {
				retBool = true
				break
			}
		}
	}
	//logging.Error("=====", retBool)
	return retBool

}
