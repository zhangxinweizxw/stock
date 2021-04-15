package notify

import (
    "fmt"
    "strconv"
    "time"

/share/models"

	"stock
/share/logging"
	"stock
/share/models/channel"
	"stock
/share/store/redis"
)

type NotifySubscribe struct {
	Model      `db:"-"`
	ID         int64 // ID
	CreateTime int64 // 创建时间
	MemberID   int64 // 会员ID
	RefID      int64 // 类型ID
	RefType    int   // 类型
	ChannelID  int64 // 频道ID
}

// 订阅的ID
type SubscribeID struct {
	RefID   int64 // 类型ID
	RefType int   // 类型
}

// 订阅列表
type SubscribeListJson struct {
	ChannelID string `json:"channel"`
	Type      int    `json:"type"`
}

// 订阅的ID Json
type SubscribeIDJson struct {
	RefID   string `json:"ref_id" binding:"required"`   // 类型ID
	RefType int    `json:"ref_type" binding:"required"` // 类型
}

// 订阅返回JSON
type SubscribeResponseJson struct {
	ChannelID string `json:"channel"`
}

// 订阅通知未读数
type UnReadNumberJson struct {
	Notify int64 `json:"notify"` // 订阅消息未读数
	Cms    int64 `json:"cms"`    // 资讯未读数
	System int64 `json:"system"` // 系统消息未读数
}

// 用户订阅信息
type MemberSubscribeInfo struct {
	SubscribeID int64 // 订阅GUID
	MemberID    int64 // 用户GUID
	ChannelID   int64 // 频道GUID
}

func NewNotifySubscribe() *NotifySubscribe {
	return &NotifySubscribe{
		Model: Model{
			TableName: TABLE_NOTIFY_SUBSCRIBE,
			Db:        MyCat,
		},
	}
}

// 订阅
func (this *NotifySubscribe) Subscribe(refType int, refID int64, memberID int64) error {

	// 数据库
	params := map[string]interface{}{
		"RefType":    refType,
		"RefID":      refID,
		"MemberID":   memberID,
		"CreateTime": time.Now().Unix(),
	}
	_, err := this.Insert(params)
	if err != nil {
		return err
	}

	// 产品缓存
	productKey := fmt.Sprintf(REDIS_NOTIFY_SUBSCRIBE_PRODUCT, memberID)
	value := fmt.Sprintf("%v#%v", refType, refID)
	redis.Sadd(productKey, []byte(value))

	// 频道缓存
	exps := map[string]interface{}{
		"RefType=?": refType,
		"RefID=?":   refID,
	}
	ids, _ := channel.NewChannel().GetIds(exps)
	channelKey := fmt.Sprintf(REDIS_NOTIFY_SUBSCRIBE_CHANNEL, memberID)
	value = fmt.Sprintf("%v", ids[0])
	redis.Sadd(channelKey, []byte(value))

	return nil
}

// 取消订阅
func (this *NotifySubscribe) UnSubscribe(refType int, refID int64, memberID int64) error {

	// 数据库
	exps := map[string]interface{}{
		"RefType=?":  refType,
		"RefID=?":    refID,
		"MemberID=?": memberID,
	}
	if err := this.Delete(exps); err != nil {
		return err
	}

	// 产品缓存
	key := fmt.Sprintf(REDIS_NOTIFY_SUBSCRIBE_PRODUCT, memberID)
	value := fmt.Sprintf("%v#%v", refType, refID)
	redis.Srem(key, []byte(value))

	// 频道缓存
	exps = map[string]interface{}{
		"RefType=?": refType,
		"RefID=?":   refID,
	}
	ids, _ := channel.NewChannel().GetIds(exps)
	channelKey := fmt.Sprintf(REDIS_NOTIFY_SUBSCRIBE_CHANNEL, memberID)
	value = fmt.Sprintf("%v", ids[0])
	redis.Srem(channelKey, []byte(value))

	return nil
}

// 判断用户是否订阅了该通知
func (this *NotifySubscribe) IsSubscribe(refType int, refID int64, memberID int64) bool {
	if memberID == 0 {
		return false
	}

	// key是否存在
	key := fmt.Sprintf(REDIS_NOTIFY_SUBSCRIBE_PRODUCT, memberID)
	exist, err := redis.Exists(key)
	if err != nil || !exist {
		this.initNotifyProductCache(memberID)
	}

	// 读缓存
	value := fmt.Sprintf("%v#%v", refType, refID)
	exist, _ = redis.Sismember(key, value)
	return exist
}

// 获取用户订阅的产品
func (this *NotifySubscribe) GetSubscribeProduct(memberID int64) []string {

	// key是否存在
	key := fmt.Sprintf(REDIS_NOTIFY_SUBSCRIBE_PRODUCT, memberID)
	exist, err := redis.Exists(key)
	if err != nil || !exist {
		this.initNotifyProductCache(memberID)
	}

	// 读缓存
	valueList, _ := redis.Smembers(key)
	return valueList
}

// 获取用户订阅的频道
func (this *NotifySubscribe) GetSubscribeChannel(memberID int64) []int64 {

	// key是否存在
	key := fmt.Sprintf(REDIS_NOTIFY_SUBSCRIBE_CHANNEL, memberID)
	exist, err := redis.Exists(key)
	if err != nil || !exist {
		channel.NewChannel().InitNotifyChannelCache(memberID)
	}

	// 读缓存
	valueList, _ := redis.Smembers(key)
	ids := make([]int64, len(valueList))
	for i, v := range valueList {
		id, _ := strconv.ParseInt(v, 10, 64)
		ids[i] = id
	}
	return ids
}

func (this *NotifySubscribe) GetlistByExps(exps map[string]interface{}) ([]*NotifySubscribe, error) {

	data := []*NotifySubscribe{}
	if len(exps) == 0 {
		return data, fmt.Errorf("Map is nil")
	}
	builder := this.Db.Select("s.*,c.ID as ChannelID").From(this.TableName+" AS s").Join(TABLE_CHANNELS+" AS c", "s.RefType = c.RefType AND s.RefID = c.RefID")
	_, err := this.SelectWhere(builder, exps).LoadStructs(&data)
	return data, err
}

// 获取服务到期的订阅信息列表
func (this *NotifySubscribe) GetServiceExpiredSubscribeList() ([]*MemberSubscribeInfo, error) {
	var data []*MemberSubscribeInfo

	exps := map[string]interface{}{
		"pmb.ServiceEndTime<=?": time.Now().Unix(),
	}

	builder := this.Db.Select("MAX(s.ID) AS SubscribeID, s.MemberID, MAX(c.ID) AS ChannelID").
		From(TABLE_PRODUCT_MEMBER_BUY+" AS pmb").
		Join(TABLE_PRODUCTS+" AS p", "pmb.ProductID=p.ID").
		Join(TABLE_NOTIFY_SUBSCRIBE+" AS s", "p.RefID=s.RefID AND p.RefType=s.RefType AND s.MemberID=pmb.MemberID").
		Join(TABLE_CHANNELS+" AS c", "c.RefID=p.RefID AND c.RefType=p.RefType")
	_, err := this.SelectWhere(builder, exps).
		GroupBy("pmb.MemberID,pmb.ProductID").
		LoadStructs(&data)

	return data, err
}

// 初始化订阅产品缓存
func (this *NotifySubscribe) initNotifyProductCache(memberID int64) {

	// 读数据库
	var ids []SubscribeID
	exps := map[string]interface{}{
		"MemberID=?": memberID,
	}
	builder := this.Db.Select("RefType, RefID").From(this.TableName)
	_, err := this.SelectWhere(builder, exps).
		LoadStructs(&ids)
	if err != nil {
		logging.Debug("Get Subscribe Product| %v", err)
		return
	}

	// 写缓存
	key := fmt.Sprintf(REDIS_NOTIFY_SUBSCRIBE_PRODUCT, memberID)
	var value string
	for _, v := range ids {
		value = fmt.Sprintf("%v#%v", v.RefType, v.RefID)
		redis.Sadd(key, value)
	}
}
