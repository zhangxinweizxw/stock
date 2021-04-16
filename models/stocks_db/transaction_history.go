package stocks_db

import (
	"fmt"
	. "stock/share/models"
	"time"
)

type TransactionHistory struct {
	Model      `db:"-" `
	StockCode  string      `db:"stock_code"`
	StockName  string      `db:"stock_name"`
	BuyTime    string      `db:"buy_time"`
	BuyPrice   float64     `db:"buy_price"`
	SellTime   interface{} `db:"sell_time"`
	SellPrice  interface{} `db:"sell_price"`
	Percentage interface{} `db:"percentage"`
}

const (
	TABLE_TRANSACTION_HISTORY = "transaction_history" // 交易历史
)

func NewTransactionHistory() *TransactionHistory {
	return &TransactionHistory{
		Model: Model{
			TableName: TABLE_TRANSACTION_HISTORY,
			Db:        MyCat,
		},
	}
}

// 是否已入库
func (this *TransactionHistory) GetTranHist(sname string) int {

	stock_code := ""
	bulid := this.Db.Select(" stock_code ").From(this.TableName).
		Where(fmt.Sprintf(" stock_name='%v' ", sname)).
		Where("ISNULL(sell_price) ")
	_, err := this.SelectWhere(bulid, nil).LoadStructs(&stock_code)
	if err != nil {
		fmt.Println("Select Table TABLE_TRANSACTION_HISTORY  |  Error   %v", err)
		return 0
	}
	//	logging.Error("=========:", len(sn))
	return len(stock_code)
}

// 未卖出数据
func (this *TransactionHistory) GetTranHistWmcList() []*TransactionHistory {

	var sc []*TransactionHistory
	bulid := this.Db.Select("*").From(this.TableName).
		Where("sell_time IS NULL AND sell_price IS NULL").
		Where(fmt.Sprintf("buy_time !='%v'", time.Now().Format("2006-01-02")))
	_, err := this.SelectWhere(bulid, nil).LoadStructs(&sc)
	if err != nil {
		fmt.Println("Select Table TABLE_TRANSACTION_HISTORY  |  Error   %v", err)
		return sc
	}
	//logging.Error("=========:", len(sc))
	return sc
}

func (this *TransactionHistory) UpdateTranHist(sc string, sp, p float64) {

	b := this.Db.Update(this.TableName).
		Set("sell_time", time.Now().Format("2006-01-02 15:04")).
		Set("sell_price", sp).
		Set("percentage", p).
		Where(fmt.Sprintf("stock_code='%v'", sc))
	_, err := this.UpdateWhere(b, nil).Exec()
	if err != nil {
		fmt.Println("Delete TABLE_TRANSACTION_HISTORY  |  Error   %v", err)
	}

}
