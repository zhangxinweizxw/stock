package stocks_db

import (
	"fmt"
	. "stock/share/models"
	"time"
)

//  AH 比价
type ZtStockDB struct {
	Model      `db:"-" `
	Id         int     `db:"id"`
	StockCode  string  `db:"stock_code"`
	StockName  string  `db:"stock_name"`
	CreateTime string  `db:"create_time"`
	Dayk5      float64 `db:"dayk5"`
	Dayk10     float64 `db:"dayk10"`
	Dayk20     float64 `db:"dayk20"`
	Dayk30     float64 `db:"dayk30"`
}

const (
	TABLE_ZT_STOCK = "zt_stock" //
)

func NewZtStockDB() *ZtStockDB {
	return &ZtStockDB{
		Model: Model{
			TableName: TABLE_ZT_STOCK,
			Db:        MyCat,
		},
	}
}

func (this *ZtStockDB) GetZtStockList() []*ZtStockDB {

	var avshStock []*ZtStockDB
	bulid1 := this.Db.Select("*").From(this.TableName)

	_, err1 := this.SelectWhere(bulid1, nil).LoadStructs(&avshStock)
	if err1 != nil {
		fmt.Println("Select Table Zt_stock |  Error   %v", err1)
		return nil
	}
	return avshStock
}

func (this *ZtStockDB) DelZtStock() {

	// 查询过滤走势比较弱的垃圾
	d := time.Now().Format("2006-01-02")
	var sd []*ZtStockDB
	bulid1 := this.Db.Select("id,f3").From("stock_day_k s").
		Join("zt_stock t", "s.f12=t.stock_code").
		Where(fmt.Sprintf("s.create_time='%v'", d)).
		Where("(f3 <1 OR f3 >6)")

	_, err1 := this.SelectWhere(bulid1, nil).LoadStructs(&sd)
	if err1 != nil {
		fmt.Println("Select Table Zt_stock join stok_day_k |  Error   %v", err1)
	}

	for _, v := range sd {
		b := this.Db.DeleteFrom(this.TableName).Where(fmt.Sprintf("id=%v", v.Id))
		_, err := this.DeleteWhere(b, nil).Exec()

		if err != nil {
			fmt.Println("Delete Table TABLE_zt_stock  |  Error   %v", err)
		}
	}

}

// 带条件删除
func (this *ZtStockDB) DelZtStockTj(sc string) {

	b := this.Db.DeleteFrom(this.TableName).Where(fmt.Sprintf("stock_code='%v'", sc))
	_, err := this.DeleteWhere(b, nil).Exec()

	if err != nil {
		fmt.Println("Delete Table TABLE_zt_stock  |  Error   %v", err)
	}

}
