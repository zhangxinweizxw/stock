package common

import (
    "fmt"

/share/models"

	"stock
/share/gocraft/dbr"
	"stock
/share/lib"
	"stock
/share/store/redis"
)

type Product struct {
	Model         `db:"-" `
	ID            int            // ID
	BeginTime     int64          // 启用时间
	CreateTime    int64          // 创建时间
	Description   dbr.NullString // 描述
	ExpireTime    int64          // 过期时间
	ServicePeriod int            // 服务周期
	IsDelete      int            // 删除标记
	ProductName   string         // 产品名称
	Price         float32        // 价格
	RefType       int            // 关联类型
	RefID         int64          // 关联ID
	RiskLevel     int            // 风险等级
	UpdateTime    int64          // 更新时间

}

type ProductJson struct {
	ID          string  `json:"_id"`          // GUID
	BeginTime   int64   `json:"begin_time"`   // 启用时间
	CreateTime  int64   `json:"create_time"`  // 创建时间
	Description string  `json:"description"`  // 描述
	ExpireTime  int64   `json:"expire_time"`  // 过期时间
	RefType     int     `json:"ref_type"`     // 关联类型
	RefID       string  `json:"ref_id"`       // 关联ID
	RiskLevel   int     `json:"risk_level"`   // 风险等级
	Price       float32 `json:"price"`        // 价格
	ProductName string  `json:"product_name"` // 产品名称
	UpdateTime  int64   `json:"update_time"`  // 更新时间
}

type BuyProduct struct {
	ID             int64
	BuyerID        int64
	ProductName    string
	ExpireTime     int64
	ServiceEndTime int64
	RefID          int64
	RefType        int
}

func NewProduct() *Product {
	return &Product{
		Model: Model{
			TableName: TABLE_PRODUCTS,
			Db:        MyCat,
		},
	}
}

func NewProductTx(tx *dbr.Tx) *Product {
	return &Product{
		Model: Model{
			TableName: TABLE_PRODUCTS,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

func (this *Product) GetSingleByExps(exps map[string]interface{}) error {
	builder := this.Db.Select("*").From(this.TableName)

	return this.SelectWhere(builder, exps).Limit(1).LoadStruct(this)
}

func (this *Product) GetListByExps(exps map[string]interface{}, limit int, orderBy string) ([]*Product, error) {
	var data []*Product

	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps).OrderBy(orderBy).Limit(uint64(limit + 1)).LoadStruct(&data)

	return data, err
}

func (this *Product) GetListJsonByExps(exps map[string]interface{}, limit int, orderBy string) ([]*ProductJson, error) {
	data, err := this.GetListByExps(exps, limit, orderBy)
	if err != nil {
		if err != dbr.ErrNotFound {
			return nil, err
		}

		return []*ProductJson{}, nil
	}

	jsns := make([]*ProductJson, len(data))
	for i, v := range data {
		jsns[i] = this.GetJson(v)
	}

	return jsns, nil
}

func (this *Product) GetCreatorById(id int) (int64, error) {
	exps := map[string]interface{}{
		"p.ID=?": id,
	}
	var memberId int64

	builder := this.Db.Select("t.MemberID").From(TABLE_TREASURE_BOX+" AS t").Join(this.TableName+" AS p", "t.ID=p.RefID AND t.Type=p.RefType")

	err := this.SelectWhere(builder, exps).Limit(1).LoadValue(&memberId)

	return memberId, err
}

func (this *Product) GetJson(p *Product) *ProductJson {
	return &ProductJson{
		ID:          IDEncrypt(int64(p.ID)),
		BeginTime:   p.BeginTime,
		ProductName: p.ProductName,
		RiskLevel:   p.RiskLevel,
		ExpireTime:  p.ExpireTime,
		Price:       p.Price,
		Description: p.Description.String,
		CreateTime:  p.CreateTime,
		UpdateTime:  p.UpdateTime,
		RefType:     p.RefType,
		RefID:       IDEncrypt(p.RefID),
	}
}

func (this *Product) UpdateData(params map[string]interface{}, exps map[string]interface{}) error {
	builder := this.Db.Update(this.TableName)
	this.UpdateParams(builder, params)
	r, err := this.UpdateWhere(builder, exps).Exec()
	if err != nil {
		return err
	}
	n, _ := r.RowsAffected()
	if n == 0 {
		return dbr.ErrNotFound
	}
	return nil
}

func (this *Product) IsPaid(refId int64, refType int, memberId int64) (bool, error) {
	if memberId == 0 {
		return false, nil
	}

	var cacheKey string

	switch refType {
	case REFTYPE_TACTIC:
		fallthrough
	case REFTYPE_REPORT:
		cacheKey = fmt.Sprintf(REDIS_MEMBER_ORDER_TREASURE_BOX, memberId)
	case REFTYPE_COURSE:
		cacheKey = fmt.Sprintf(REDIS_MEMBER_ORDER_ACTIVITY, memberId)
	case REFTYPE_ASSEMBLE:
		cacheKey = fmt.Sprintf(REDIS_MEMBER_ORDER_ASSEMBLE, memberId)
	case REFTYPE_STRATEGY:
		cacheKey = fmt.Sprintf(REDIS_MEMBER_ORDER_STRATEGY, memberId)
	default:
		return false, fmt.Errorf("RefType [%v] Not Found", refType)
	}

	paidArray, err := redis.Smembers(cacheKey)

	if err != nil {
		return false, err
	}

	if lib.InArray(paidArray, FormatInt(refId)) == false {
		return false, nil
	}

	return true, nil
}

func (this *Product) Pay(refId int64, refType int, memberId int64) error {
	var cacheKey string

	switch refType {
	case REFTYPE_TACTIC:
		fallthrough
	case REFTYPE_REPORT:
		cacheKey = fmt.Sprintf(REDIS_MEMBER_ORDER_TREASURE_BOX, memberId)
	case REFTYPE_COURSE:
		cacheKey = fmt.Sprintf(REDIS_MEMBER_ORDER_ACTIVITY, memberId)
	default:
		return fmt.Errorf("RefType [%v] Not Found", refType)
	}

	_, err := redis.Sadd(cacheKey, []byte(FormatInt(refId)))
	return err
}

func (this *Product) UnPay(refId int64, refType int, memberId int64) error {
	var cacheKey string

	switch refType {
	case REFTYPE_TACTIC:
		fallthrough
	case REFTYPE_REPORT:
		cacheKey = fmt.Sprintf(REDIS_MEMBER_ORDER_TREASURE_BOX, memberId)
	case REFTYPE_COURSE:
		cacheKey = fmt.Sprintf(REDIS_MEMBER_ORDER_ACTIVITY, memberId)
	default:
		return fmt.Errorf("RefType [%v] Not Found", refType)
	}

	return redis.Srem(cacheKey, []byte(FormatInt(refId)))
}

func (this *Product) GetState(beginTime int64, expireTime int64, nowTime int64, status int) int {
	var state int

	// 停售
	if status == TREASUREBOX_STATUS_STOP_SELL {
		state = TREASURE_STATE_STOP_SELL
		return state
	}

	if nowTime < beginTime {
		state = TREASURE_STATE_NOT_RUNNING
	} else if nowTime >= beginTime && nowTime < expireTime {
		state = TREASURE_STATE_RUNNING
	} else {
		state = TREASURE_STATE_ENDED
	}

	return state
}

// 获取订购的策略ID列表
func (this *Product) GetStrategyIdList(exps map[string]interface{}) ([]int64, error) {
	var ids []int64

	builder := this.Db.Select("p.RefID").From(this.TableName+" AS p").
		Join(TABLE_ORDERS+" AS o", "o.ProductID=p.ID")
	_, err := this.SelectWhere(builder, exps).LoadStructs(&ids)

	return ids, err
}

// 获取产品订购数
func (this *Product) GetSubscribeCount(exps map[string]interface{}) (int64, error) {
	var count int64

	builder := this.Db.Select("COUNT(0)").From(this.TableName+" AS p").
		Join(TABLE_ORDERS+" AS o", "o.ProductID=p.ID")
	err := this.SelectWhere(builder, exps).LoadStruct(&count)

	return count, err
}

// 获取已购买的产品列表
func (this *Product) GetAllBuyProduct() ([]*BuyProduct, error) {
	var productList []*BuyProduct

	exps := map[string]interface{}{
		"o.Status=?":         ORDER_STATUS_PAID,
		"o.IsDelete=?":       0,
		"o.ParentID=?":       0,
		"p.IsDelete=?":       0,
		"p.RefType not in ?": []int{REFTYPE_MEMBER_UPDATE},
	}
	builder := this.Db.Select("o.MemberID AS BuyerID, p.ID, p.ProductName, p.ExpireTime,p.RefID, p.RefType, IFNULL(pmb.ServiceEndTime,0) AS ServiceEndTime").
		From(this.TableName+" AS p").
		Join(TABLE_ORDERS+" AS o", "p.ID=o.ProductID").
		Join(TABLE_PRODUCT_MEMBER_BUY+" AS pmb", "o.ID=pmb.OrderID")
	_, err := this.SelectWhere(builder, exps).LoadStructs(&productList)

	// 锦囊、内参获取服务过期时间
	for i, v := range productList {
		if v.RefType == REFTYPE_REPORT || v.RefType == REFTYPE_TACTIC {
			productList[i].ExpireTime = v.ServiceEndTime
		}
	}
	return productList, err
}
