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

func (this *XQ_Stock_FX) DelStockFx() {

	var create_time string
	bulid := this.Db.SelectBySql("SELECT create_time FROM xq_fx GROUP BY create_time ORDER BY create_time DESC LIMIT 4,1")
	_, err := this.SelectWhere(bulid, nil).LoadStructs(&create_time)
	if err != nil {
		fmt.Println("Select Table xq_fx table  create_time |  Error   %v", err)
	}

	b := this.Db.DeleteFrom(this.TableName).
		Where(fmt.Sprintf("create_time < '%v'", create_time))
	_, err1 := this.DeleteWhere(b, nil).Exec()

	if err1 != nil {
		fmt.Println("Delete Table xq_fx   |  Error   %v", err1)
	}
}

// xq_fx 表中数据筛选
func (this *XQ_Stock_FX) GetXqFxStockList() []*XQ_Stock_FX {
	var xqfxL []*XQ_Stock_FX
	bulid := this.Db.SelectBySql("SELECT * from ( SELECT COUNT(1) o,stock_code,stock_name FROM xq_fx GROUP BY stock_code ORDER BY o DESC ) b WHERE b.o >1 AND b.o <4 ")
	_, err := this.SelectWhere(bulid, nil).LoadStructs(&xqfxL)
	if err != nil {
		fmt.Println("Select Table xq_fx table  |  Error   %v", err)
		return nil
	}

	return xqfxL
}
