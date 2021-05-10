package message

import (
	"fmt"

	. "stock/share/models"

	"encoding/json"
	"strconv"

	"sync"

	"stock/share/gocraft/dbr"
	"stock/share/middleware"
	"stock/share/models/channel"
	"stock/share/models/common"
	"stock/share/models/upload"
	"stock/share/store/redis"
)

type Message struct {
	Model      `db:"-"`
	ID         int64
	AffixID    int64          // 附件ID
	Content    dbr.NullString // 消息文本内容
	CreateTime int64          // 创建时间
	FromUser   int64          // 发送人
	Inline     dbr.NullString // 消息内联信息
	IsFavorite int            // 收藏标记
	IsUndo     int            // 撤回标记
	IsUnread   int            // 未读标记
	ToTarget   int64          // 消息接收目标
	ToType     int            // 消息类型 MESSAGE_TARGET_TYPE_%v
	Type       int            // 消息主体类型 MESSAGE_BODY_TYPE_%v
	State      string         // 用户自定义数据
}

// 消息目标（群组或者会话）
type MessageTargetJson struct {
	RefID         string        `json:"ref_id"`
	RefType       int           `json:"ref_type"`
	Messages      []MessageJson `json:"messages",array`
	LastDoingTime int64         `json:"last_doing_time"`
}

type MessageJson struct {
	GUID       string             `json:"_id"`
	Body       MesageBodyJson     `json:"body"`        // 消息主体
	CreateTime float64            `json:"create_time"` // 创建时间
	From       *common.MemberJson `json:"from"`        // 发送人节点
	Type       int                `json:"type"`        // 消息主体类型
	State      string             `json:"state"`       // 用户自定义数据
}

type MesageBodyJson struct {
	Attachment *upload.UploadAttachmentJson `json:"attachment"` // 附件
	Content    string                       `json:"content"`    // 正文
	Inline     *InlineJson                  `json:"inline"`     // 内联
}

type InlineJson struct {
	Fields   []InlineFieldJson `json:"fields"`   // 字段
	Link     string            `json:"link"`     // 链接
	Template string            `json:"template"` // 模版
	Title    string            `json:"title"`    // 标题
	Comment  string            `json:"comment"`  // 内容
}

type InlineFieldJson struct {
	Item  string      `json:"item"`  // 子项标题
	Value interface{} `json:"value"` // 子项内容
	Link  string      `json:"link"`  // 子项链接
}

type TopicGatherJson struct {
	GUID    string      `json:"_id"`
	Message MessageJson `json:"message"`
}

var memberJson map[int64]*common.MemberJson = make(map[int64]*common.MemberJson)
var l sync.RWMutex

// --------------------------------------------------------------------------------

func NewMessage() *Message {
	return &Message{
		Model: Model{
			TableName: TABLE_MESSAGES,
			Db:        MyCat,
		},
	}
}

func NewMessageTx(tx *dbr.Tx) *Message {
	return &Message{
		Model: Model{
			TableName: TABLE_MESSAGES,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

func (this *Message) GetSingle(id int64, memberId int64) error {
	exps := map[string]interface{}{
		"m.ID=?": id,
	}
	builder := this.Db.Select("m.*, r.ID IS NULL AS IsUnread").
		From(this.TableName+" AS m").
		LeftJoin(TABLE_MESSAGE_READ+" AS r", fmt.Sprintf("r.RefID=m.ID AND r.MemberID=%v", memberId))

	err := this.SelectWhere(builder, exps).Limit(1).
		LoadStruct(&this)

	return err
}

func (this *Message) GetListByExps(exps map[string]interface{}, sort string, limit uint64, memberId int64, conditions ...dbr.Condition) ([]*Message, error) {
	var data []*Message

	builder := this.Db.Select(`m.*`).
		From(this.TableName+" AS m").
		LeftJoin(TABLE_MESSAGE_READ+" AS r", fmt.Sprintf("r.RefID=m.ID AND r.MemberID=%v", memberId))

	_, err := this.SelectWhere(builder, exps, conditions...).
		//		OrderBy("CreateTime DESC").
		OrderBy("ID " + sort).
		Limit(limit + 1).
		LoadStructs(&data)

	return data, err
}

func (this *Message) GetLiveListExps(exps map[string]interface{}, limit uint64) ([]*Message, error) {
	var data []*Message

	builder := this.Db.Select(`*`).From(this.TableName + " AS m")

	_, err := this.SelectWhere(builder, exps).
		//		OrderBy("CreateTime DESC").
		OrderBy("ID DESC").
		Limit(limit + 1).
		LoadStructs(&data)

	return data, err

}

func (this *Message) GetFavoritesListByExps(exps map[string]interface{}, limit uint64) ([]*Message, error) {
	var data []*Message
	builder := this.Db.Select(`m.*, 1 AS IsFavorite`).
		From(this.TableName+" AS m").
		Join(TABLE_MESSAGE_FAVORITES+" AS mf", "m.ID=mf.RefID")

	_, err := this.SelectWhere(builder, exps).
		//		OrderBy("CreateTime DESC").
		OrderBy("ID DESC").
		Limit(limit + 1).
		LoadStructs(&data)

	return data, err
}

// 清除所有消息未读数
func (this *Message) ClearUnreadAll(memberId int64) error {
	res, err := redis.Keys(fmt.Sprintf(REIDS_MEMBERS_UNREAD, memberId))
	if err != nil {
		return err
	}

	for _, v := range res {
		if err := redis.Del(v); err != nil {
			return err
		}
	}
	return nil
}

// 设置消息未读数
func (this *Message) SetUnread(toTarget int64, toType int, client *middleware.Client) error {
	switch toType {
	case MESSAGE_TARGET_TYPE_CHANNEL:
		ch := channel.NewChannel()
		return ch.SetUnread(toTarget, client.Member.ID)

	case MESSAGE_TARGET_TYPE_SESSION:
		sess := NewSession()
		return sess.SetUnread(toTarget, client.Member.ID)

	default:
		return dbr.ErrNotSupported
	}
	return nil
}

// 获取最后操作时间
func (this *Message) GetLastDoingTime(refId int64, refType int, client *middleware.Client) (int64, error) {
	switch refType {
	case MESSAGE_TARGET_TYPE_CHANNEL:
		ch := channel.NewChannel()
		num, err := ch.GetLastDoingTime(refId)
		return num, err

	case MESSAGE_TARGET_TYPE_SESSION:
		sess := NewSession()
		num, err := sess.GetLastDoingTime(refId)
		return num, err

	case MESSAGE_TARGET_TYPE_LIVE:
		live := NewMessageLive()
		num, err := live.GetLastDoingTime(refId)
		return num, err

	default:
		return 0, dbr.ErrNotSupported
	}

	return 0, nil
}

// 获取最后消息内容
func (this *Message) GetLastMessage(refId int64, refType int, client *middleware.Client) (string, error) {
	return "", nil
}

// 刷新最后操作时间
func (this *Message) RefreshLastDoingTime(refId int64, refType int, client *middleware.Client) error {
	switch refType {
	case MESSAGE_TARGET_TYPE_CHANNEL:
		ch := channel.NewChannel()
		return ch.RefreshLastDongTime(refId)

	case MESSAGE_TARGET_TYPE_SESSION:
		sess := NewSession()
		return sess.RefreshLastDongTime(refId)

	case MESSAGE_TARGET_TYPE_LIVE:
		live := NewMessageLive()
		return live.RefreshLastDoingTime(refId)

	default:
		return dbr.ErrNotSupported
	}
}

// 刷新最后消息内容
func (this *Message) RefreshLastMessage(refId int64, refType int, data string, client *middleware.Client) error {
	switch refType {
	case MESSAGE_TARGET_TYPE_CHANNEL:
		ch := channel.NewChannel()
		return ch.RefreshLastMessage(refId, data)

	case MESSAGE_TARGET_TYPE_SESSION:
		sess := NewSession()
		return sess.RefreshLastMessage(refId, data)

	default:
		return dbr.ErrNotSupported
	}

	return nil
}

// 获取直播聚集
func (this *Message) GetTopicGather(exps map[string]interface{}, limit uint64) ([]*Message, error) {
	return []*Message{}, nil
}

func (this *Message) GetLastMessageID() int64 {
	var id int64

	builder := this.Db.Select("ID").From(this.TableName).OrderBy("ID DESC")
	this.SelectWhere(builder, nil).Limit(1).LoadValue(&id)

	return id
}

func (this *Message) GetMessageIncrID() int64 {
	id, err := redis.Do("INCR", REDIS_ADVISOR_MESSAGE_INC_ID)
	if err != nil || id.(int64) < 1 {
		return this.GetLastMessageID()
	}

	return id.(int64)
}

func (this *Message) GetSingleJson(id int64, client *middleware.Client) (*MessageTargetJson, error) {

	if err := this.GetSingle(id, client.Member.ID); err != nil {
		return nil, err
	}

	messages := []*Message{this}
	jsn, err := getMessageTargetJson(messages)
	if err != nil {
		return nil, err
	}

	jsn.RefID = IDEncrypt(this.ToTarget)
	jsn.RefType = this.ToType

	return jsn, nil
}

func (this *Message) GetMultiJson(exps map[string]interface{}, limit int, sort string, refId string, toType int, client *middleware.Client, conditions ...dbr.Condition) (*MessageTargetJson, error) {
	var err error
	var messages []*Message

	messages, err = this.GetListByExps(exps, sort, uint64(limit), client.Member.ID, conditions...)
	if err != nil {
		return nil, err
	}

	jsn, err := getMessageTargetJson(messages)
	if err != nil {
		return nil, err
	}
	jsn.RefID = refId
	jsn.RefType = toType
	return jsn, nil
}

func getMessageTargetJson(messages []*Message) (*MessageTargetJson, error) {
	jsns := MessageTargetJson{
		Messages: make([]MessageJson, len(messages)),
	}

	if len(messages) == 0 {
		return &jsns, nil
	}

	for i, m := range messages {
		var jsn MessageJson
		jsn.GUID = IDEncrypt(m.ID)
		jsn.CreateTime, _ = strconv.ParseFloat(fmt.Sprintf("%.6f", float64(m.CreateTime)/1000000000), 64)
		if memberJsn, err := GetCahceMember(m.FromUser); err != nil {
			return &jsns, err
		} else {
			jsn.From = memberJsn
		}
		jsn.Type = m.Type
		jsn.State = m.State
		// 非撤回消息
		if m.IsUndo == 0 {

			// 正文内容
			jsn.Body = MesageBodyJson{Content: m.Content.String}

			// 上传附件信息
			if m.AffixID != 0 {
				ua := upload.NewUploadAffixs()
				if err := ua.GetSingle(m.AffixID); err == nil {
					j, _ := ua.GetSingleJson(ua)
					jsn.Body.Attachment = &j

					// 关联数据
					exps := map[string]interface{}{
						"RefID=?":   m.ToTarget,
						"RefType=?": m.ToType,
						"DataID=?":  m.ID,
					}
					uar := upload.NewUploadRelevance()
					if relevanceId, err := uar.GetRelevanceIdByExps(exps); err == nil {
						jsn.Body.Attachment.Relevance = IDEncrypt(relevanceId)
					}
				}
			}

			// 内联信息
			if len(m.Inline.String) > 0 {
				err := json.Unmarshal([]byte(m.Inline.String), &jsn.Body.Inline)
				if err != nil {
					return nil, err
				}
			}
		} else {
			jsn.Body = MesageBodyJson{Content: "消息已撤回"}
			jsn.Type = MESSAGE_BODY_TYPE_UNDO
		}

		jsns.Messages[i] = jsn

	}
	if len(jsns.Messages) == 0 {
		jsns.Messages = []MessageJson{}
	}
	// 重置内存数据
	memberJson = make(map[int64]*common.MemberJson)
	return &jsns, nil
}

func GetCahceMember(memberID int64) (*common.MemberJson, error) {
	l.Lock()
	defer l.Unlock()
	if mjsn, ok := memberJson[memberID]; ok {
		return mjsn, nil
	}
	member := common.NewMember()
	if err := member.GetSingle(memberID); err == nil {
		m, err := member.GetSingleJson(member)
		memberJson[memberID] = m
		return m, err

	} else {
		return nil, err
	}
}
