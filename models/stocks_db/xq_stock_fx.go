package stocks_db

import (
	"fmt"
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

// xq_fx 表中数据筛选
func (this *XQ_Stock_FX) GetXqFxStockList() []*XQ_Stock_FX {
	var xqfxL []*XQ_Stock_FX
	bulid := this.Db.Select("*").
		From(" ( SELECT COUNT(1) o,stock_code,stock_name FROM xq_fx GROUP BY stock_code ORDER BY o DESC ) b WHERE b.o >1 AND b.o <4")
	_, err := this.SelectWhere(bulid, nil).LoadStructs(&xqfxL)
	if err != nil {
		fmt.Println("Select Table xq_fx table  |  Error   %v", err)
		return nil
	}

	return xqfxL
}
