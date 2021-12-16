package stocks_db

import (
	"fmt"
	. "stock/share/models"
	"time"
)

type DxStockDb struct {
	Model      `db:"-" `
	Id         int    `db:"id"`
	StockCode  string `db:"stock_code"`
	StockName  string `db:"stock_name"`
	CreateTime string `db:"create_time"`
	//Status     int     `db:"status"`
	DayK5  float64 `db:"dayk5"`
	DayK10 float64 `db:"dayk10"`
	DayK20 float64 `db:"dayk20"`
	DayK30 float64 `db:"dayk30"`
}

const (
	TABLE_DX_STOCK = "dx_stock" // 短线选股
)

func NewDxStockDb() *DxStockDb {
	return &DxStockDb{
		Model: Model{
			TableName: TABLE_DX_STOCK,
			Db:        MyCat,
		},
	}
}

// 最新交易日短线数据
func (this *DxStockDb) GetDxStockList() []*DxStockDb {

	// 查询最新日期
	ctime := ""
	bulid := this.Db.Select("create_time").From(this.TableName).
		OrderBy("create_time DESC").Limit(1)
	_, err := this.SelectWhere(bulid, nil).LoadStructs(&ctime)
	if err != nil {
		fmt.Println("Select Table dx_stock  |  Error   %v", err)
		return nil
	}

	var dxStock []*DxStockDb
	bulid1 := this.Db.Select("*").From(this.TableName).
		Where(fmt.Sprintf("create_time='%v'", ctime))

	_, err1 := this.SelectWhere(bulid1, nil).LoadStructs(&dxStock)
	if err1 != nil {
		fmt.Println("Select Table dx_stock |  Error   %v", err1)
		return nil
	}
	return dxStock
}

//func (this *DxStockDb) DelDxStock() {
//
//	b := this.Db.DeleteFrom(this.TableName)
//	_, err := this.DeleteWhere(b, nil).Exec()
//
//	if err != nil {
//		fmt.Println("Delete Table TABLE_dx_stock  |  Error   %v", err)
//	}
//
//}
func (this *DxStockDb) DelDxStock() {

	// 查询过滤走势比较弱的垃圾
	d := time.Now().Format("2006-01-02")
	var sd []*DxStockDb
	bulid1 := this.Db.Select("id,f3").From("stock_day_k s").
		Join("dx_stock t", "s.f12=t.stock_code").
		Where(fmt.Sprintf("s.create_time='%v'", d)).
		Where(" ( f3 < 0.28 or f62 < 0 or f8 < 0.58 or f10 < 0.58 ) ")

	_, err1 := this.SelectWhere(bulid1, nil).LoadStructs(&sd)
	if err1 != nil {
		fmt.Println("Select Table dx_stock join stok_day_k |  Error   %v", err1)
	}

	for _, v := range sd {
		b := this.Db.DeleteFrom(this.TableName).Where(fmt.Sprintf("id=%v", v.Id))
		_, err := this.DeleteWhere(b, nil).Exec()

		if err != nil {
			fmt.Println("Delete Table TABLE_dx_stock  |  Error   %v", err)
		}
	}

}
