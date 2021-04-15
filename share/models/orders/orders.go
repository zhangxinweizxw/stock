package orders

import (
    "fmt"

/share/models"

	redigo "stock
/share/garyburd/redigo/redis"
	"stock
/share/gocraft/dbr"
	"stock
/share/logging"
	"stock
/share/store/redis"
)

type Order struct {
	Model              `db:"-"`
	ID                 int64
	TreasureAdvisorID  dbr.NullInt64  // 百宝箱投顾GUID
	CourseAdvisorID    dbr.NullInt64  // 课堂投顾GUID
	AssembleAdvisorID  dbr.NullInt64  // 复合产品投顾GUID
	BatchNo            string         // 订单号
	CancelTime         int64          // 取消时间
	CreateTime         int64          // 创建时间
	IsDelete           int            // 删除标记
	MemberID           int64          // 会员ID
	PayMethod          int            // 支付方式 0.线下支付、1.微信支付、2.支付宝支付
	PayTime            int64          // 支付时间
	ProductID          int64          // 产品ID
	ProductName        string         // 产品名称
	ProductPrice       float32        // 产品价格
	ProductBeginTime   dbr.NullInt64  // 启用时间
	ProductDescription dbr.NullString // 产品简介
	ProductEndTime     dbr.NullInt64  // 过期时间
	ProductRefID       dbr.NullInt64  // 产品关联ID
	ProductRefType     dbr.NullInt64  // 产品关联类型
	Price              float32        // 订单价格
	Status             int            // 订单状态 1.待支付、2.已支付、3.已取消、4.提交退款、5.退款中、6.已退款
	Transcation        string         // 支付交易流水号
	ProtocolNumber     string         // 订单协议编号
	TreasureState      dbr.NullInt64  // 百宝箱状态
	CourseState        dbr.NullInt64  // 课堂状态
	AssembleState      dbr.NullInt64  // 复合产品状态
	Extra1             dbr.NullString // 通用参数（策略Guid）
	UpdateTime         int64          // 更新时间
	LinkUrl            string         // 跳转路径
}

func NewOrder() *Order {
	return &Order{
		Model: Model{
			TableName: TABLE_ORDERS,
			Db:        MyCat,
		},
	}
}

func NewOrderTx(tx *dbr.Tx) *Order {
	return &Order{
		Model: Model{
			TableName: TABLE_ORDERS,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

// 获取用户购买产品的订单总金额
func (this *Order) GetOrderSumPrice(exps map[string]interface{}) (float64, error) {
	builder := this.Db.Select("SUM(o.`Price`").
		From(this.TableName+" AS o").
		Join(TABLE_PRODUCTS+" AS p", "o.`ProductID`=p.`ID`")

	var sum float64
	err := this.SelectWhere(builder, exps).Limit(1).LoadStruct(&sum)
	if err != nil {
		if err == dbr.ErrNotFound {
			return 0, nil
		}
		return -1, err
	}
	return sum, nil
}

// 获取用户购买产品的订单数量
func (this *Order) GetOrderCount(exps map[string]interface{}) (int, error) {
	builder := this.Db.Select("COUNT(1)").
		From(this.TableName+" AS o").
		Join(TABLE_PRODUCTS+" AS p", "o.`ProductID`=p.`ID`")

	var count int
	err := this.SelectWhere(builder, exps).Limit(1).LoadStruct(&count)
	if err != nil {
		if err == dbr.ErrNotFound {
			return 0, nil
		}
		return -1, err
	}
	return count, nil
}

// 判断 memberId 是否老用户
func IsOldUser(memberId int64) (bool, error) {
	// 和产品沟通需求后不能用 redis 来优化判断用户是否老用户，原因是：
	// 用户购买一个产品后退款了，再次购买产品相当于新购买产品，不享受折扣
	//b, err := IsOldUserByRedis(memberId)
	//if err == nil {
	//	return b, nil
	//}
	return IsOldUserByMysql(memberId)
}

// 通过 mysql 订单记录判断 memberId 是否老用户
// 老用户定义：至少购买过一次 锦囊、内参、课堂、组合产品的客户(量化策略除外)
func IsOldUserByMysql(memberId int64) (bool, error) {
	order := NewOrder()
	count, err := order.GetOrderCount(map[string]interface{}{
		"o.MemberID=?":   memberId,
		"o.Status=?":     ORDER_STATUS_PAID,
		"p.RefType IN ?": []int{REFTYPE_TACTIC, REFTYPE_REPORT, REFTYPE_COURSE, REFTYPE_ASSEMBLE},
	})
	if err != nil && err != dbr.ErrNotFound {
		logging.Error("GetOrderCount | memberId %v | %v", memberId, err)
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func memberOrderScard(keyFmt string, memberId int64) (int, error) {
	key := fmt.Sprintf(keyFmt, memberId)
	n, err := redis.Scard(key)
	if err != nil && err != redigo.ErrNil {
		logging.Error("Scard %v | %v", key, err)
		return n, err
	}
	return n, nil
}

// 通过 redis 订单记录判断 memberId 是否老用户
// 老用户定义：至少购买过一次 锦囊、内参、课堂、组合产品的客户(量化策略除外)
func IsOldUserByRedis(memberId int64) (bool, error) {
	// 锦囊、内参
	n, err := memberOrderScard(REDIS_MEMBER_ORDER_TREASURE_BOX, memberId)
	if err != nil {
		return false, err
	}
	if n > 0 {
		return true, nil
	}

	// 活动课堂
	n, err = memberOrderScard(REDIS_MEMBER_ORDER_ACTIVITY, memberId)
	if err != nil {
		return false, err
	}
	if n > 0 {
		return true, nil
	}

	// 组合产品(牛人计划)
	n, err = memberOrderScard(REDIS_MEMBER_ORDER_ASSEMBLE, memberId)
	if err != nil {
		return false, err
	}
	if n > 0 {
		return true, nil
	}

	return false, nil
}
