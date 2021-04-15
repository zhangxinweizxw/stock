package stocks_db

import (
	"fmt"
	. "stock/share/models"
)

//  AH 比价
type AvsHStock struct {
	Model      `db:"-" `
	StockCode  string `db:"stock_code"`
	StockName  string `db:"stock_name"`
	CreateTime string `db:"create_time"`
}

const (
	TABLE_AVSH_STOCK = "avsh_stock" // ah 股比价
)

func NewAvsHStock() *AvsHStock {
	return &AvsHStock{
		Model: Model{
			TableName: TABLE_AVSH_STOCK,
			Db:        MyCat,
		},
	}
}

func (this *AvsHStock) GetAvsHStockList() []*AvsHStock {

	// 查询最新日期
	ctime := ""
	bulid := this.Db.Select("create_time").From(this.TableName).
		OrderBy("create_time DESC").Limit(1)
	_, err := this.SelectWhere(bulid, nil).LoadStructs(&ctime)
	if err != nil {
		fmt.Println("Select Table avsh_stock  |  Error   %v", err)
		return nil
	}
	var avshStock []*AvsHStock
	bulid1 := this.Db.Select("*").From(this.TableName).
		Where(fmt.Sprintf("create_time='%v'", ctime))

	_, err1 := this.SelectWhere(bulid1, nil).LoadStructs(&avshStock)
	if err1 != nil {
		fmt.Println("Select Table avsh_stock |  Error   %v", err1)
		return nil
	}
	return avshStock
}

func (this *AvsHStock) DelAvshStock() {

	b := this.Db.DeleteFrom(this.TableName)
	_, err := this.DeleteWhere(b, nil).Exec()

	if err != nil {
		fmt.Println("Delete Table TABLE_zjlx_stock  |  Error   %v", err)
	}

}
