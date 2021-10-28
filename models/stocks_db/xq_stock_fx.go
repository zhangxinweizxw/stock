package stocks_db

import (
	. "stock/share/models"
)

const (
	XQ_STOCK_FX = "xq_fx" //
)

type XQ_Stock_FX struct {
	Model      `db:"-" `
	StockName  string `db:"stock_name"`
	StockCode  string `db:"stock_code"`
	CreateTime string `db:"create_time"`
}

func NewXQ_Stock_FX() *XQ_Stock_FX {
	return &XQ_Stock_FX{
		Model: Model{
			TableName: XQ_STOCK_FX,
			Db:        MyCat,
		},
	}
}
