package strategy

import (
    "fmt"
    "time"

/share/models"

	"stock
/share/gocraft/dbr"
)

type Strategy struct {
	Model           `db:"-" `
	ID              int64          // GUID
	StrategyID      string         // 策略GUID
	Assort          int            // 分类
	BeginTime       int64          // 开始时间
	CreateTime      int64          // 创建时间
	EndTime         int64          // 结束时间
	Description     dbr.NullString // 简介
	IsRecommend     int            // 推荐
	Order           int64          // 排序
	Price           float32        // 价格
	Status          int            // 状态（0禁用，1启用）
	SuccessRate     float32        // 成功率
	TotalProfitRate float32        // 总收益率
	YearProfitRate  float32        // 年收益率
	Tags            string         // 策略风格
	Title           string         // 标题
	UpdateTime      int64          // 修改时间
}

type StrategyBriefJson struct {
	StrategyID      string  `json:"_id"`
	ID              string  `json:"ref_id"`
	Title           string  `json:"title"`
	BeginTime       int64   `json:"begin_time"`
	EndTime         int64   `json:"end_time"`
	YearProfitRate  float32 `json:"profit_year"`
	TotalProfitRate float32 `json:"profit_total"`
	SuccessRate     float32 `json:"success"`
}

// 策略的过往收益情况
type StrategyProfitBrief struct {
	StrategyID       string
	ID               int64
	Title            string
	YearProfitRate   float32
	ProfitMaxCode    string
	ProfitMaxName    string
	ProfitMaxSetCode string
	ProfitMaxDays    int
	ProfitMax        float32
	IsRecommend      int
}

// 策略的过往收益情况JSON
type StrategyProfitBriefJson struct {
	StrategyID       string  `json:"_id"`
	ID               string  `json:"ref_id"`
	Title            string  `json:"title"`
	YearProfitRate   float32 `json:"profit_year"`
	ProfitMaxCode    string  `json:"profit_max_code"`
	ProfitMaxName    string  `json:"profit_max_name"`
	ProfitMaxSetCode string  `json:"profit_max_setcode"`
	ProfitMaxDays    int     `json:"profit_max_days"`
	ProfitMax        float32 `json:"profit_max"`
	IsRecommend      int     `json:"is_recommend"`
}

func NewStrategy() *Strategy {
	return &Strategy{
		Model: Model{
			TableName: TABLE_STRATEGYS,
			Db:        MyCat,
		},
	}
}

func NewStrategyTx(tx *dbr.Tx) *Strategy {
	return &Strategy{
		Model: Model{
			TableName: TABLE_STRATEGYS,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

func (this *Strategy) GetSingle(exps map[string]interface{}) error {
	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps).
		LoadStruct(this)

	return err
}

func (this *Strategy) GetSubscribeList(memberID int64, limit uint64, page uint64) ([]*Strategy, error) {
	var data []*Strategy

	exps := map[string]interface{}{
		"MemberID=?": memberID,
	}
	builder := this.Db.Select("s.*").From(this.TableName+" AS s").
		Join(TABLE_NOTIFY_SUBSCRIBE+" AS ns", fmt.Sprintf("ns.RefType=%v AND ns.RefID=s.ID", REFTYPE_STRATEGY))
	_, err := this.SelectWhere(builder, exps).
		OrderBy("ns.CreateTime DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		LoadStructs(&data)

	return data, err
}

// 获取策略列表（过往收益）
func (this *Strategy) GetProfitList(exps map[string]interface{}, limit uint64, page uint64, orderBy string) ([]*StrategyProfitBriefJson, error) {
	var briefList []*StrategyProfitBrief

	builder := this.Db.Select("s.ID,s.StrategyID,s.Title,s.YearProfitRate,IFNULL(sp.ProfitMaxCode,\"\") AS ProfitMaxCode,IFNULL(sp.ProfitMaxName,\"\") AS ProfitMaxName,IFNULL(sp.ProfitMaxSetCode,\"\") AS ProfitMaxSetCode,IFNULL(sp.ProfitMaxDays,0) AS ProfitMaxDays,IFNULL(sp.ProfitMax,0.0000) AS ProfitMax,s.IsRecommend").
		From(this.TableName+" AS s").
		LeftJoin(TABLE_STRATEGY_STOCK_PROFIT+" AS sp", "s.StrategyID=sp.StrategyID AND sp.Status=1 AND sp.IsDelete=0")
	_, err := this.SelectWhere(builder, exps).
		OrderBy(orderBy).
		Limit(limit).
		Offset(limit * (page - 1)).
		LoadStructs(&briefList)
	if err != nil {
		return nil, err
	}

	data := make([]*StrategyProfitBriefJson, len(briefList))
	for i, v := range briefList {
		data[i] = this.getProfitBriefJson(v)
	}
	return data, nil
}

func (this *Strategy) getProfitBriefJson(profit *StrategyProfitBrief) *StrategyProfitBriefJson {
	var data StrategyProfitBriefJson

	data.ID = IDEncrypt(profit.ID)
	data.StrategyID = profit.StrategyID
	data.Title = profit.Title
	data.YearProfitRate = profit.YearProfitRate
	data.ProfitMaxCode = profit.ProfitMaxCode
	data.ProfitMaxName = profit.ProfitMaxName
	data.ProfitMaxDays = profit.ProfitMaxDays
	data.ProfitMaxSetCode = profit.ProfitMaxSetCode
	data.ProfitMax = profit.ProfitMax
	data.IsRecommend = profit.IsRecommend

	return &data
}

//
func (this *Strategy) GetStateByID(id int64) int {

	exps := map[string]interface{}{
		"ID=?": id,
	}
	builder := this.Db.Select("*").
		From(this.TableName)
	this.SelectWhere(builder, exps).
		LoadStruct(this)

	var state int
	now := time.Now().Unix()
	if now < this.BeginTime {
		state = TREASURE_RUN_STATE_BEFORE_SELLING
	} else if now >= this.BeginTime && now <= this.EndTime {
		state = TREASURE_RUN_STATE_RUNNING
	} else {
		state = TREASURE_RUN_STATE_ENDED
	}
	return state
}
