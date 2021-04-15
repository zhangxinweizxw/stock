package message

import (
    "fmt"
    "strconv"
    "time"

/share/models"

	"stock
/share/gocraft/dbr"
	"stock
/share/store/redis"
)

type Session struct {
	Model        `db:"-"`
	ID           int64
	SID          int64          // 会话ID
	CreateTime   int64          // 创建时间
	Creator      int64          // 创建人
	IsDeleted    int            // 删除标记
	IsHelper     int            // 助手标记
	IsTop        int            // 置顶标记
	MemberAvatar dbr.NullString // 成员头像
	MemberDesc   dbr.NullString // 成员描述
	MemberID     int64          // 团队成员ID
	MemberName   dbr.NullString // 成员昵称
	RefChannel   int64          // 关联的订阅频道（成员ID）
	RefID        int64          // 关联的成员ID
	UpdateTime   int64          // 更新时间
}

type SessionJson struct {
	GUID          string `json:"_id"`
	Avatar        string `json:"avatar"`          // 头像
	CreateTime    int64  `json:"create_time"`     // 创建时间
	Creator       string `json:"creator"`         // 创建人
	Description   string `json:"desc"`            // 描述
	IsHelper      int    `json:"is_helper"`       // 助手标记
	IsTop         int    `json:"is_top"`          // 置顶标记
	LastDoingTime int64  `json:"last_doing_time"` // 最后操作时间
	LastMessage   string `json:"last_message"`    // 最后消息内容
	Name          string `json:"name"`            // 成员名称
	NamePinyin    string `json:"name_pinyin"`     // 成员名称（拼音）
	RefID         string `json:"ref_id"`          // 关联成员ID
	Type          int    `json:"type"`            // 类型
	Unread        int    `json:"unread"`          // 未读消息数
}

// --------------------------------------------------------------------------------

func NewSession() *Session {
	return &Session{
		Model: Model{
			TableName: TABLE_SESSIONS,
			Db:        MyCat,
		},
	}
}

func NewSessionTx(tx *dbr.Tx) *Session {
	return &Session{
		Model: Model{
			TableName: TABLE_SESSIONS,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

func (this *Session) GetSingleByExps(exps map[string]interface{}) error {
	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps).
		Limit(1).
		LoadStruct(&this)

	return err
}

func (this *Session) GetMemberIdsBySID(exps map[string]interface{}) ([]int64, error) {
	var ids []int64
	builder := this.Db.Select("MemberID").From(this.TableName)
	err := this.SelectWhere(builder, exps).LoadValue(&ids)
	return ids, err
}

func (this *Session) GetListByExps(exps map[string]interface{}) ([]Session, error) {
	var data []Session
	builder := this.Db.Select(
		"s.*, IFNULL(m.IsHelper,0) AS IsHelper, m.FriendlyName AS MemberName, m.Description AS MemberDesc, m.Avatar AS MemberAvatar").
		From(this.TableName+" AS s").
		LeftJoin(TABLE_MEMBERS+" AS m", "s.RefID=m.ID")
	_, err := this.SelectWhere(builder, exps).LoadStructs(&data)

	return data, err
}

func (this *Session) GetListLimitByExps(exps map[string]interface{}) ([]Session, error) {
	var data []Session
	builder := this.Db.Select(
		"s.*, IFNULL(m.IsHelper,0) AS IsHelper, m.FriendlyName AS MemberName, m.Description AS MemberDesc, m.Avatar AS MemberAvatar").
		From(this.TableName+" AS s").
		LeftJoin(TABLE_MEMBERS+" AS m", "s.RefID=m.ID").OrderBy("ID DESC")
	_, err := this.SelectWhere(builder, exps).LoadStructs(&data)

	return data, err
}

// 获取会话序列号
func (this *Session) GetSequenceNum() (int64, error) {
	rec, _ := redis.Do("INCR", REDIS_MAJOR_SEQ_SESSION)
	seqnum := rec.(int64)
	if seqnum == 1 {
		var data int64
		err := this.Db.Select("IFNULL(MAX(SID),0) AS MaxID").From(this.TableName).LoadValue(&data)
		if err != nil {
			return data, err
		}

		err = redis.Set(REDIS_MAJOR_SEQ_SESSION, []byte(fmt.Sprintf("%v", data+1)))
		return data + 1, err
	}
	return seqnum, nil
}

// 获取未阅读数
func (this *Session) GetUnread(memberId int64) int {
	if memberId < 1 {
		return 0
	}
	num, _ := redis.Get(
		fmt.Sprintf(REDIS_MEMBERS_SESSIONS_UNREAD, memberId, this.SID))

	result, _ := strconv.Atoi(num)
	return result
}

// 设置未阅读数
func (this *Session) SetUnread(id int64, currMemberId int64) error {
	exps := map[string]interface{}{
		"SID=?": id,
	}
	ids, err := this.GetMemberIdsBySID(exps)
	if err == nil {
		for _, memberId := range ids {
			if currMemberId == memberId {
				continue
			}
			_, err = redis.Do("INCR", fmt.Sprintf(REDIS_MEMBERS_SESSIONS_UNREAD, memberId, id))
			if err != nil {
				return err
			}
		}
	}
	return err
}

// 未阅读数清零
func (this *Session) SetZeroUnread(id int64, memberId int64) error {
	return redis.Del(fmt.Sprintf(REDIS_MEMBERS_SESSIONS_UNREAD, memberId, id))
}

// 读取最后操作时间
func (this *Session) GetLastDoingTime(id int64) (int64, error) {
	lastDoingTime, err := redis.Get(
		fmt.Sprintf(REDIS_SESSIONS_LAST_DOING_TIME, id))

	result, _ := strconv.ParseInt(lastDoingTime, 10, 64)
	return result, err
}

// 读取最后消息内容
func (this *Session) GetLastMessage(id int64) (string, error) {
	result, err := redis.Get(
		fmt.Sprintf(REDIS_SESSIONS_LAST_MESSAGE, id))

	return result, err
}

// 刷新最后操作时间
func (this *Session) RefreshLastDongTime(id int64) error {
	return redis.Set(
		fmt.Sprintf(REDIS_SESSIONS_LAST_DOING_TIME, id),
		[]byte(FormatInt(time.Now().Unix())))
}

// 刷新最后消息内容
func (this *Session) RefreshLastMessage(id int64, data string) error {
	return redis.Set(
		fmt.Sprintf(REDIS_SESSIONS_LAST_MESSAGE, id),
		[]byte(data))
}

// 更新删除标记
// 当A成员向B成员发送消息时，需要将此会话删除标记清零
func (this *Session) UpdateFieldIsDeleted(id int64) error {
	params := map[string]interface{}{
		"IsDeleted":  0,
		"UpdateTime": time.Now().Unix(),
	}
	exps := map[string]interface{}{
		"SID=?":       id,
		"IsDeleted=?": 1,
	}

	return this.Update(params, exps)
}
