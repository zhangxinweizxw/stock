package stocks_db

import (
	"fmt"
	. "stock/share/models"
)

//  千股千评
type QgqpStockDb struct {
	Model      `db:"-" `
	StockCode  string `db:"stock_code"`
	StockName  string `db:"stock_name"`
	CreateTime string `db:"create_time"`
}

const (
	TABLE_QGQP_STOCK = "qgqp_stock"
)

func NewQgqpStockDb() *QgqpStockDb {
	return &QgqpStockDb{
		Model: Model{
			TableName: TABLE_QGQP_STOCK,
			Db:        MyCat,
		},
	}
}

func (this *QgqpStockDb) GetQgqpStockList() []*QgqpStockDb {

	var qgqpStock []*QgqpStockDb
	bulid1 := this.Db.Select("*").From(this.TableName)

	_, err1 := this.SelectWhere(bulid1, nil).LoadStructs(&qgqpStock)
	if err1 != nil {
		fmt.Println("Select Table qgqp_stock |  Error   %v", err1)
		return nil
	}
	return qgqpStock
}

func (this *QgqpStockDb) DelQgqpStock() {

	b := this.Db.DeleteFrom(this.TableName)
	_, err := this.DeleteWhere(b, nil).Exec()

	if err != nil {
		fmt.Println("Delete Table TABLE_qgqp_stock  |  Error   %v", err)
	}

}
