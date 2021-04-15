package stocks_db

import (
	"fmt"
	. "stock/share/models"
)

type StockInfo struct {
	Model     `db:"-" `
	Date      string `db:"date"`
	StockName string `db:"stock_name"`
	StockCode string `db:"stock_code"`
}

const (
	TABLE_STOCK_INFO = "stock_info" // 简单个股信息
)

func NewStockInfo() *StockInfo {
	return &StockInfo{
		Model: Model{
			TableName: TABLE_STOCK_INFO,
			Db:        MyCat,
		},
	}
}

func (this *StockInfo) DelStockInfo() {

	b := this.Db.DeleteFrom(this.TableName)
	_, err := this.DeleteWhere(b, nil).Exec()

	if err != nil {
		fmt.Println("Delete Table TABLE_STOCK_INFO  |  Error   %v", err)
	}

}
func (this *StockInfo) GetStockInfoList() []*StockInfo {
	var s []*StockInfo
	bulid1 := this.Db.Select("*").From(this.TableName)

	_, err1 := this.SelectWhere(bulid1, nil).LoadStructs(&s)
	if err1 != nil {
		fmt.Println("Select Table_STOCK_INFO |  Error   %v", err1)
		return nil
	}
	return s
}
