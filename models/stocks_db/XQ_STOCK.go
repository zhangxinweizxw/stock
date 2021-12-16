package stocks_db

import (
	"fmt"
	. "stock/share/models"
	"time"
)

const (
	XQ_STOCK = "xq_stock" //
)

type XQ_Stock struct {
	Model      `db:"-" `
	Id         int    `db:"id"`
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

//func (this *XQ_Stock) DelXqStock() {
//
//	b := this.Db.DeleteFrom(this.TableName)
//	_, err := this.DeleteWhere(b, nil).Exec()
//
//	if err != nil {
//		fmt.Println("Delete Table TABLE_xq_stock  |  Error   %v", err)
//	}
//
//}
func (this *XQ_Stock) DelXqStock() {

	// 查询过滤走势比较弱的垃圾
	d := time.Now().Format("2006-01-02")
	var sd []*XQ_Stock
	bulid1 := this.Db.Select("id,f3").From("stock_day_k s").
		Join("xq_stock t", "s.f12=t.stock_code").
		Where(fmt.Sprintf("s.create_time='%v'", d)).
		Where(" ( f3 < 0.28 or f62 < 0 or f8 < 0.58 or f10 < 0.58 ) ")

	_, err1 := this.SelectWhere(bulid1, nil).LoadStructs(&sd)
	if err1 != nil {
		fmt.Println("Select Table xq_stock join stok_day_k |  Error   %v", err1)
	}

	for _, v := range sd {
		b := this.Db.DeleteFrom(this.TableName).Where(fmt.Sprintf("id=%v", v.Id))
		_, err := this.DeleteWhere(b, nil).Exec()

		if err != nil {
			fmt.Println("Delete Table TABLE_xq_stock  |  Error   %v", err)
		}
	}

}
