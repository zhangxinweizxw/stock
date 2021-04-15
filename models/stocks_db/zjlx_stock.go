package stocks_db

import (
	"fmt"
	. "stock/share/models"
)

//  主力流入
type ZjlxStockDb struct {
	Model      `db:"-" `
	StockCode  string `db:"stock_code"`
	StockName  string `db:"stock_name"`
	CreateTime string `db:"create_time"`
}

const (
	TABLE_ZJLX_STOCK = "zjlx_stock" // 个股资金流向
)

func NewZjlxStockDb() *ZjlxStockDb {
	return &ZjlxStockDb{
		Model: Model{
			TableName: TABLE_ZJLX_STOCK,
			Db:        MyCat,
		},
	}
}

func (this *ZjlxStockDb) GetZjlxStockList() []*ZjlxStockDb {

	var zjlxStock []*ZjlxStockDb
	bulid1 := this.Db.Select("*").From(this.TableName)

	_, err1 := this.SelectWhere(bulid1, nil).LoadStructs(&zjlxStock)
	if err1 != nil {
		fmt.Println("Select Table zjlx_stock |  Error   %v", err1)
		return nil
	}
	return zjlxStock
}

func (this *ZjlxStockDb) DelZjlxStock() {

	b := this.Db.DeleteFrom(this.TableName)
	_, err := this.DeleteWhere(b, nil).Exec()

	if err != nil {
		fmt.Println("Delete Table TABLE_zjlx_stock  |  Error   %v", err)
	}

}
