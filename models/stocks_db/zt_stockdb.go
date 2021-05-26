package stocks_db

import (
	"fmt"
	. "stock/share/models"
)

//  AH 比价
type ZtStockDB struct {
	Model      `db:"-" `
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

	b := this.Db.DeleteFrom(this.TableName)
	_, err := this.DeleteWhere(b, nil).Exec()

	if err != nil {
		fmt.Println("Delete Table TABLE_zt_stock  |  Error   %v", err)
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
