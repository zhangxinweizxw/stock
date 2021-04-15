package stocks_db

import (
	"fmt"
	. "stock/share/models"
)

const (
	XQ_STOCK = "xq_stock" //
)

type XQ_Stock struct {
	Model      `db:"-" `
	StockName  string `db:"stock_name"`
	StockCode  string `db:"stock_code"`
	CreateTime string `db:"create_time"`
}

func NewXQ_Stock() *XQ_Stock {
	return &XQ_Stock{
		Model: Model{
			TableName: XQ_STOCK,
			Db:        MyCat,
		},
	}
}

func (this *XQ_Stock) DelXqStock() {

	b := this.Db.DeleteFrom(this.TableName)
	_, err := this.DeleteWhere(b, nil).Exec()

	if err != nil {
		fmt.Println("Delete Table TABLE_xq_stock  |  Error   %v", err)
	}

}
