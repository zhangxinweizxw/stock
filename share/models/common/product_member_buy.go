package common

import (
    "fmt"
    "strconv"
    "time"

/share/models"

	"stock
/share/gocraft/dbr"
	"stock
/share/logging"
	"stock
/share/store/redis"
)

type ProductMemberBuy struct {
	Model            `db:"-" `
	ID               int64 // ID
	CreateTime       int64 // 创建时间
	IsDelete         int   // 删除标识
	MemberID         int64 // 用户ID
	OrderID          int64 // 订单ID
	ProductID        int64 // 产品ID
	ServiceBeginTime int64 // 服务开始时间
	ServiceEndTime   int64 // 服务结束时间
	RefType          int   // 产品类型
}

// 服务时间
type ServiceTime struct {
	ServiceBeginTime int64 // 服务开始时间
	ServiceEndTime   int64 // 服务结束时间
}

func NewProductMemberBuy() *ProductMemberBuy {
	return &ProductMemberBuy{
		Model: Model{
			TableName: TABLE_PRODUCT_MEMBER_BUY,
			Db:        MyCat,
		},
	}
}

func NewProductMemberBuyTx(tx *dbr.Tx) *ProductMemberBuy {
	return &ProductMemberBuy{
		Model: Model{
			TableName: TABLE_PRODUCT_MEMBER_BUY,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

func (this *ProductMemberBuy) GetSingle(exps map[string]interface{}, orderBy string) error {
	builder := this.Db.Select("*").From(this.TableName)
	if len(orderBy) > 0 {
		builder = builder.OrderBy(orderBy)
	}
	return this.SelectWhere(builder, exps).Limit(1).LoadStruct(this)
}

func (this *ProductMemberBuy) GetList(exps map[string]interface{}, limit uint64, orderBy string) ([]*ProductMemberBuy, error) {
	var data []*ProductMemberBuy

	builder := this.Db.Select("pmb.*").
		From(this.TableName+" AS pmb").
		Join(TABLE_PRODUCTS+" AS p", "pmb.ProductID = p.ID")
	if len(orderBy) > 0 {
		builder = builder.OrderBy(orderBy)
	}
	if limit > 0 {
		builder = builder.Limit(limit)
	}
	builder = this.SelectWhere(builder, exps)
	_, err := builder.LoadStructs(&data)

	return data, err
}

// 获取指定产品的最后服务结束时间
func (this *ProductMemberBuy) GetMaxServiceEndTime(exps map[string]interface{}) (int64, error) {
	var serviceEndTime dbr.NullInt64
	builder := this.Db.Select("MAX(pmb.ServiceEndTime)").From(this.TableName+" AS pmb").
		Join(TABLE_PRODUCTS+" AS p", "p.ID = pmb.ProductID").
		Join(TABLE_TREASURE_BOX+" AS t", "t.Type = p.RefType AND t.ID = p.RefID")
	err := this.SelectWhere(builder, exps).Limit(1).LoadStruct(&serviceEndTime)

	return serviceEndTime.Int64, err
}

// 判断用户能否购买该产品
func (this *ProductMemberBuy) IsPermitBuyProduct(mID int64, productID int) bool {

	now := time.Now().Unix()

	// 读缓存
	key := fmt.Sprintf(REDIS_BUY_TIME_PERMIT, mID, productID)
	value, err := redis.Get(key)
	if err == nil {
		valueInt, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			if now > valueInt {
				return true
			}
			return false
		}
	}
	redis.Del(key)

	// 读数据库
	var permitTime int64
	pmb := NewProductMemberBuy()
	exps := map[string]interface{}{
		"pmb.MemberID=?":  mID,
		"pmb.ProductID=?": productID,
		"pmb.IsDelete=?":  0,
	}
	data, err := pmb.GetList(exps, 2, "pmb.ServiceEndTime DESC")
	if err != nil {
		logging.Debug("Get Product Member Buy | %v", err)
	}
	if len(data) == 2 {
		permitTime = data[1].ServiceEndTime
	} else {
		permitTime = 0
	}

	// 写缓存
	redis.Set(key, []byte(fmt.Sprintf("%v", permitTime)))

	return now > permitTime
}

// 获取课堂的最后服务结束时间 课堂改版 wdk 20170810 add
func (this *ProductMemberBuy) GetCoursesMaxServiceEndTime(exps map[string]interface{}) (int64, error) {
	var serviceEndTime dbr.NullInt64
	builder := this.Db.Select("MAX(pmb.ServiceEndTime)").From(this.TableName+" AS pmb").
		Join(TABLE_PRODUCTS+" AS p", "p.ID = pmb.ProductID").
		Join(TABLE_COURSES+" AS t", "t.Type = p.RefType AND t.ID = p.RefID")
	err := this.SelectWhere(builder, exps).Limit(1).LoadStruct(&serviceEndTime)

	return serviceEndTime.Int64, err
}

// 获取服务时间列表
func (this *ProductMemberBuy) GetServiceTimeList(exps map[string]interface{}) ([]ServiceTime, error) {
	var timeList []ServiceTime

	builder := this.Db.Select("pmb.ServiceBeginTime,pmb.ServiceEndTime").From(this.TableName+" AS pmb").
		Join(TABLE_PRODUCTS+" AS p", "p.ID = pmb.ProductID").
		Join(TABLE_TREASURE_BOX+" AS t", "t.Type = p.RefType AND t.ID = p.RefID")
	_, err := this.SelectWhere(builder, exps).LoadStructs(&timeList)

	return timeList, err
}

// 判断用户和产品的关系 是否在购买运行中，是否购买的已结束 wdk 20170817 add
func (this *ProductMemberBuy) UserProductRelation(mID int64, productID int) (bool, bool) {

	now := time.Now().Unix()

	// 读数据库
	var permitTime int64
	var serviceEndtime0 int64
	pmb := NewProductMemberBuy()
	exps := map[string]interface{}{
		"pmb.MemberID=?":  mID,
		"pmb.ProductID=?": productID,
		"pmb.IsDelete=?":  0,
	}
	data, err := pmb.GetList(exps, 2, "pmb.ServiceEndTime DESC")
	if err != nil {
		logging.Debug("Get Product Member Buy | %v", err)
	}
	if len(data) == 2 {
		permitTime = data[1].ServiceEndTime
		serviceEndtime0 = data[0].ServiceEndTime
	} else if len(data) == 1 {
		serviceEndtime0 = data[0].ServiceEndTime
	} else {
		permitTime = 0
		serviceEndtime0 = 0
	}

	return now > permitTime, now > serviceEndtime0 //是否能买第二期，是否在运行中
}
