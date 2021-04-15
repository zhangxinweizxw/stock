package strategy

/share/models"

	"stock
/share/gocraft/dbr"
)

type StrategyStockProfit struct {
	Model      `db:"-" `
	ID         int64  // GUID
	StrategyID string // 策略GUID(hn_strategys.ID)
	UpdateTime int64  // 修改时间
}

func NewStrategyStockProfit() *StrategyStockProfit {
	return &StrategyStockProfit{
		Model: Model{
			TableName: TABLE_STRATEGY_STOCK_PROFIT,
			Db:        MyCat,
		},
	}
}

func NewStrategyStockProfitTx(tx *dbr.Tx) *StrategyStockProfit {
	return &StrategyStockProfit{
		Model: Model{
			TableName: TABLE_STRATEGY_STOCK_PROFIT,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

func (this *StrategyStockProfit) GetSingle(exps map[string]interface{}) error {
	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps).
		LoadStruct(this)

	return err
}
