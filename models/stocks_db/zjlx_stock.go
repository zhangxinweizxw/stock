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

	//// 查询最新日期
	//var ctime []string
	//bulid := this.Db.Select("create_time").From(this.TableName).
	//	GroupBy("create_time").
	//	OrderBy("create_time DESC").Limit(3)
	//_, err := this.SelectWhere(bulid, nil).LoadStructs(&ctime)
	//if err != nil {
	//	fmt.Println("Select Table dx_stock  |  Error   %v", err)
	//	return nil
	//}

	var zjlxStock []*ZjlxStockDb
	bulid1 := this.Db.Select("stock_name,stock_code ").From(this.TableName)

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
