package channel

import (
    "fmt"
    "strconv"
    "strings"
    "time"

/share/models"

	"stock
/share/gocraft/dbr"
	"stock
/share/logging"
	"stock
/share/store/redis"
)

type Channel struct {
	Model       `db:"-"`
	ID          int64
	ChannelName string         // 频道名称
	Color       string         // 颜色
	CreateTime  int64          // 创建时间
	Creator     int64          // 创建人
	Description dbr.NullString // 描述
	IsSystem    int            // 系统自动创建标记
	MemberIds   dbr.NullString // 群组成员串
	UpdateTime  int64          // 更新时间
	Visibility  int            // 可见范围
}

type ChannelJson struct {
	GUID          string                `json:"_id"`
	Color         string                `json:"color"`           // 颜色
	CreateTime    int64                 `json:"create_time"`     // 创建时间
	Creator       string                `json:"creator"`         // 创建人
	Description   dbr.NullString        `json:"desc"`            // 描述
	IsSystem      int                   `json:"is_system"`       // 系统自动创建标记
	IsTop         int                   `json:"is_top"`          // 置顶标记
	Joined        int                   `json:"joined"`          // 加入标记
	LastDoingTime int64                 `json:"last_doing_time"` // 最后操作时间
	LastMessage   string                `json:"last_message"`    // 最后消息内容
	Members       dbr.NullString        `json:"member_ids"`      // 群组成员串
	Name          string                `json:"name"`            // 频道名称
	NamePinyin    string                `json:"name_pinyin"`     // 频道名称（拼音）
	Preference    ChannelPreferenceJson `json:"preference"`      // 偏好设置节点
	Show          int                   `json:"show"`            // 显示标记
	Type          int                   `json:"type"`            // 类型
	Unread        int                   `json:"unread"`          // 未读消息数
	Visibility    int                   `json:"visibility"`      // 可见范围
}

type ChannelPreferenceJson struct {
	NotifyDesktop int `json:"notify_desktop"`
	NotifyMobile  int `json:"notify_mobile"`
}

// --------------------------------------------------------------------------------

func NewChannel() *Channel {
	return &Channel{
		Model: Model{
			CacheKey:  REDIS_CHANNELS,
			Db:        MyCat,
			TableName: TABLE_CHANNELS,
		},
	}
}

func NewChannelTx(tx *dbr.Tx) *Channel {
	return &Channel{
		Model: Model{
			CacheKey:  REDIS_CHANNELS,
			Db:        MyCat,
			TableName: TABLE_CHANNELS,
			Tx:        tx,
		},
	}
}

func (this *Channel) GetSingle(id int64) error {
	cacheKey := fmt.Sprintf(this.CacheKey, id)
	rec, err := redis.Hgetall(cacheKey)
	if err == nil && len(rec) > 1 {
		if err := MapToStruct(this, rec); err != nil {
			redis.Del(cacheKey)
		} else {
			return nil
		}
	}

	exps := map[string]interface{}{
		"ID=?": id,
	}
	builder := this.Db.Select("*").From(this.TableName)
	err = this.SelectWhere(builder, exps).
		Limit(1).
		LoadStruct(&this)
	if err != nil {
		return err
	}

	return redis.Hmset(cacheKey, StructToMap(this))
}

func (this *Channel) GetSingleByExps(exps map[string]interface{}) error {
	builder := this.Db.Select("*").From(this.TableName)
	return this.SelectWhere(builder, exps).Limit(1).LoadStruct(this)
}

func (this *Channel) GetListByMemberID(memberId int64) ([]Channel, error) {

	return nil, nil
}

// 获取未阅读数
func (this *Channel) GetUnread(memberId int64) int {
	if this.ID == 0 {
		return 0
	}

	num, _ := redis.Get(
		fmt.Sprintf(REDIS_MEMBERS_CHANNELS_UNREAD, memberId, this.ID))

	result, _ := strconv.Atoi(num)
	return result
}

// 设置未阅读数
func (this *Channel) SetUnread(refId int64, currMemberId int64) error {
	exps := map[string]interface{}{
		"RefID=?": refId,
	}
	cm := NewChannelMember()
	ids, err := cm.GetMemberIdsByExps(exps)
	if err == nil {
		for _, id := range ids {
			if currMemberId == id {
				continue
			}
			_, err := redis.Do("INCR", fmt.Sprintf(REDIS_MEMBERS_CHANNELS_UNREAD, id, refId))
			if err != nil {
				return err
			}
		}
	}
	return err
}

// 未阅读数清零
func (this *Channel) SetZeroUnread(id int64, memberId int64) error {
	return redis.Del(fmt.Sprintf(REDIS_MEMBERS_CHANNELS_UNREAD, memberId, id))
}

// 读取最后操作时间
func (this *Channel) GetLastDoingTime(id int64) (int64, error) {
	lastDongTime, err := redis.Get(
		fmt.Sprintf(REDIS_CHANNELS_LAST_DOING_TIME, id))
	result, _ := strconv.ParseInt(lastDongTime, 10, 64)
	return result, err
}

// 读取最后消息内容
func (this *Channel) GetLastMessage(id int64) (string, error) {
	result, err := redis.Get(
		fmt.Sprintf(REDIS_CHANNELS_LAST_MESSAGE, id))

	return result, err
}

// 刷新最后操作时间
func (this *Channel) RefreshLastDongTime(id int64) error {
	return redis.Set(
		fmt.Sprintf(REDIS_CHANNELS_LAST_DOING_TIME, id),
		[]byte(FormatInt(time.Now().Unix())))
}

// 刷新最后消息内容
func (this *Channel) RefreshLastMessage(id int64, data string) error {
	return redis.Set(
		fmt.Sprintf(REDIS_CHANNELS_LAST_MESSAGE, id),
		[]byte(data))
}

// 更新频道成员
func (this *Channel) UpdateMemberIds(id int64) error {
	var memberIds []string

	exps := map[string]interface{}{
		"RefID=?": id,
	}
	cm := NewChannelMember()
	ids, err := cm.GetMemberIdsByExps(exps)
	if err == nil {
		for _, id := range ids {
			memberIds = append(memberIds, IDEncrypt(id))
		}
	}

	idString := strings.Join(memberIds, ",")

	rec, err := redis.Hgetall(fmt.Sprintf(this.CacheKey, id))
	if err == nil && len(rec) > 1 {
		redis.Do("HMSET", fmt.Sprintf(this.CacheKey, id), "MemberIds", []byte(idString))
	}

	params := map[string]interface{}{
		"MemberIds": idString,
	}
	exps = map[string]interface{}{
		"ID=?": id,
	}

	return this.Update(params, exps)
}

// 获取频道的MemberId,仅针对百宝箱业务
func (this *Channel) GetMemberIDByChanneID(chId int64) (int64, int) {
	if chId < 1 {
		return 0, 0
	}
	exps := map[string]interface{}{
		"ch.ID=?": chId,
	}
	var tempData struct {
		MemberID int64
		RefType  int
	}
	builder := this.Db.Select("t.MemberID,ch.RefType").From(this.TableName+" AS ch").Join(TABLE_TREASURE_BOX+" AS t", "ch.RefID=t.ID")
	this.SelectWhere(builder, exps).Limit(1).LoadStruct(&tempData)
	return tempData.MemberID, tempData.RefType
}

func (this *Channel) GetLastChannelID() int64 {
	var (
		chanellMaxID int64
		sessionID    int64
	)
	builder := this.Db.Select(`MAX("ID")`).From(this.TableName)
	this.SelectWhere(builder, nil).LoadValue(chanellMaxID)
	builder = this.Db.Select(`MAX("ID")`).From(TABLE_SESSIONS)
	this.SelectWhere(builder, nil).LoadValue(sessionID)
	if chanellMaxID > sessionID {
		return chanellMaxID + 1
	} else {
		return sessionID + 1
	}
}

func (this *Channel) GetChannelIncrID() int64 {
	id, err := redis.Do("INCR", REDIS_ADVISOR_CHANNEL_INC_ID)
	if err != nil || id.(int64) < 1 {
		cid := this.GetLastChannelID()
		redis.Set(REDIS_ADVISOR_CHANNEL_INC_ID, []byte(strconv.Itoa(int(cid))))
		return cid
	}

	return id.(int64)
}

// 初始化用户订阅频道
func (this *Channel) InitNotifyChannelCache(memberID int64) []int64 {

	// 读数据库
	var ids []int64
	exps := map[string]interface{}{
		"n.MemberID=?": memberID,
	}
	builder := this.Db.Select("c.ID").From(TABLE_NOTIFY_SUBSCRIBE+" AS n").
		Join(TABLE_CHANNELS+" AS c", fmt.Sprintf("n.RefType=c.RefType AND n.RefID=c.RefID "))
	_, err := this.SelectWhere(builder, exps).
		LoadStructs(&ids)
	if err != nil {
		logging.Debug("Get Subscribe Channel | %v", err)
		return ids
	}

	// 写缓存
	key := fmt.Sprintf(REDIS_NOTIFY_SUBSCRIBE_CHANNEL, memberID)
	redis.Del(key)
	var value string
	for _, v := range ids {
		value = fmt.Sprintf("%v", v)
		redis.Sadd(key, value)
	}
	return ids
}
