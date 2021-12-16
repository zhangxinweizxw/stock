package stocks_db

import (
	"fmt"
	. "stock/share/models"
	"time"
)

//  千股千评
type QgqpStockDb struct {
	Model      `db:"-" `
	Id         int    `db:"id"`
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

//func (this *QgqpStockDb) DelQgqpStock() {
//
//	b := this.Db.DeleteFrom(this.TableName)
//	_, err := this.DeleteWhere(b, nil).Exec()
//
//	if err != nil {
//		fmt.Println("Delete Table TABLE_qgqp_stock  |  Error   %v", err)
//	}
//
//}
func (this *QgqpStockDb) DelQgqpStock() {

	// 查询过滤走势比较弱的垃圾
	d := time.Now().Format("2006-01-02")
	var sd []*QgqpStockDb
	bulid1 := this.Db.Select("id,f3").From("stock_day_k s").
		Join("qgqp_stock t", "s.f12=t.stock_code").
		Where(fmt.Sprintf("s.create_time='%v'", d)).
		Where(" ( f3 < 0.28 or f62 < 0 or f8 < 0.58 or f10 < 0.58 ) ")

	_, err1 := this.SelectWhere(bulid1, nil).LoadStructs(&sd)
	if err1 != nil {
		fmt.Println("Select Table qgqp_stock join stok_day_k |  Error   %v", err1)
	}

	for _, v := range sd {
		b := this.Db.DeleteFrom(this.TableName).Where(fmt.Sprintf("id=%v", v.Id))
		_, err := this.DeleteWhere(b, nil).Exec()

		if err != nil {
			fmt.Println("Delete Table TABLE_qgqp_stock  |  Error   %v", err)
		}
	}

}
