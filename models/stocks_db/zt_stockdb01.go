package stocks_db

import (
	"fmt"
	. "stock/share/models"
)

//  AH 比价
type ZtStockDB01 struct {
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
	TABLE_ZT_STOCK01 = "zt_stock_01" //
)

func NewZtStockDB01() *ZtStockDB01 {
	return &ZtStockDB01{
		Model: Model{
			TableName: TABLE_ZT_STOCK01,
			Db:        MyCat,
		},
	}
}

//func (this *ZtStockDB) GetZtStockList() []*ZtStockDB {
//
//	var avshStock []*ZtStockDB
//	bulid1 := this.Db.Select("*").From(this.TableName)
//
//	_, err1 := this.SelectWhere(bulid1, nil).LoadStructs(&avshStock)
//	if err1 != nil {
//		fmt.Println("Select Table Zt_stock |  Error   %v", err1)
//		return nil
//	}
//	return avshStock
//}
//
func (this *ZtStockDB01) DelZtStock01() {

	b := this.Db.DeleteFrom(this.TableName)
	_, err := this.DeleteWhere(b, nil).Exec()

	if err != nil {
		fmt.Println("Delete Table TABLE_zt_stock01  |  Error   %v", err)
	}
}

func (this *ZtStockDB01) GetZtStockList01() []*ZtStockDB01 {

	var avshStock []*ZtStockDB01
	bulid1 := this.Db.Select("*").From(this.TableName)

	_, err1 := this.SelectWhere(bulid1, nil).LoadStructs(&avshStock)
	if err1 != nil {
		fmt.Println("Select Table Zt_stock01 |  Error   %v", err1)
		return nil
	}
	return avshStock
}
