package stocks_db

import (
	"fmt"
	. "stock/share/models"
)

//  F10数据 财务分析 报告期
type FinancialReports struct {
	Model     `db:"-" `
	Date      string `db:"date"`
	StockCode string `db:"stock_code"`
	Jbmgsy    string `db:"jbmgsy"`
	Mgjzc     string `db:"mgjzc"`
	Mgwfply   string `db:"mgwfply"`
	Yyzsr     string `db:"yyzsr"`
	Gsjlr     string `db:"gsjlr"`
	Yyzsrtbzz string `db:"yyzsrtbzz"`
	Gsjlrtbzz string `db:"gsjlrtbzz"`
	Kfjlrtbzz string `db:"kfjlrtbzz"`
}

const (
	TABLE_FINANNCIAL_REPORTS = "financial_reports" // 简单个股信息
)

func NewFinancialReports() *FinancialReports {
	return &FinancialReports{
		Model: Model{
			TableName: TABLE_FINANNCIAL_REPORTS,
			Db:        MyCat,
		},
	}
}

func (this *FinancialReports) DelFinancialReports() {

	b := this.Db.DeleteFrom(this.TableName)
	_, err := this.DeleteWhere(b, nil).Exec()

	if err != nil {
		fmt.Println("Delete TABLE_FINANNCIAL_REPORTS  |  Error   %v", err)
	}

}

// 营业总收入同比增长、扣非净利润同比增长
func (this *FinancialReports) GetTbzz(cd string) bool {

	type Tb struct {
		Zsr float64 `db:"zsr"`
		Jlr float64 `db:"jlr"`
	}
	var tb *Tb
	bulid := this.Db.Select("avg(yyzsrtbzz) zsr,AVG(kfjlrtbzz) jlr").From("`financial_reports`").
		Where(fmt.Sprintf("stock_code='%v'", cd))
	_, err := this.SelectWhere(bulid, nil).LoadStructs(&tb)
	if err != nil {
		fmt.Println("Select Table Financial_reports  |  Error   %v", err)
		return false
	}

	if tb.Zsr > 10 && tb.Jlr > 8 {
		return true
	}
	return false
}
