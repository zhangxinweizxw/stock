package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"stock/config"
	"stock/share/logging"
	"stock/share/util"
)

type UtilHttps struct {
	C *config.AppConfig
}

func NewUtilHttps(cfg *config.AppConfig) *UtilHttps {
	return &UtilHttps{
		C: cfg,
	}
}

// 雪球判断财务跟最近涨跌
func (this *UtilHttps) GetXqPd(scode string) int {

	sc := this.GetUtilCode(scode)
	if len(sc) <= 0 {
		return 0
	}

	url := `https://xueqiu.com/service/screener/screen?category=CN&exchange=sh_sz&areacode=&indcode=&order_by=symbol&order=desc&page=1&size=30&only_count=0&npay.20210630=1_5000&oiy.20210630=1_5000&npay.20210930=1_5000&oiy.20210930=1_5000&mc=2500000000_150000000000&pct5=-2.8_8.8&pct20=-5.8_13`
	url += fmt.Sprintf(`&symbol=%v`, sc)

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
	//logging.Error("========:", data.XQResuData.Count)
	return data.XQResuData.Count
}

func (this *UtilHttps) GetUtilCode(scode string) string {
	sc := ""
	switch scode[:3] {
	case "600", "601", "603", "605", "688", "689", "608":
		sc = fmt.Sprintf("SH%v", scode)
	case "300", "002", "000", "001", "003", "301":
		sc = fmt.Sprintf("SZ%v", scode)
	}
	return sc
}

func (this *UtilHttps) GetUtilCode1(scode string) string {
	stockCodes := ""
	switch scode[:3] {
	case "600", "601", "603", "605", "688", "689", "608":
		stockCodes = fmt.Sprintf("1.%v", scode)
	case "300", "002", "000", "001", "003", "301":
		stockCodes = fmt.Sprintf("0.%v", scode)
	}
	return stockCodes
}
