package common

import (
    "fmt"
    "strconv"
    "strings"

/share/models"

	"stock
/share/lib"
	"stock
/share/store/redis"
)

type Member struct {
	Model          `db:"-" `
	ID             int64
	MCode          int64  // 账户代码，来自 hn_member_code.Code
	AdvisorType    int    // 投顾类型：0.普通会员、1.内部员工、2.持牌投顾
	Avatar         string // 用户头像
	CreateTime     int64  // 创建时间
	FriendlyName   string // 全称（名字）
	GroupID        int64  // 会员组ID
	IsHelper       int    // 助手标签
	LastUpdateTime int64  // 最后更新时间
	MAN            string // 号码(4B+SK.AES)
	Mobile         string // 手机号
	Status         int    // 状态
	DeviceToken    string // jpush Token
}

type MemberJson struct {
	GUID               string `json:"_id"`
	AdvisorType        int    `json:"advisor_type"` // 投顾类型
	Avatar             string `json:"avatar"`       // 用户头像
	CreateTime         int64  `json:"create_time"`  // 创建时间
	FriendlyName       string `json:"name"`         // 全称（名字）
	FriendlyNamePinyin string `json:"name_pinyin"`  // 全称拼音（名字）
	MCode              int64  `json:"uid"`          // 会员账户代码
}

type Advisor struct {
	Model          `db:"-" `
	ID             int64  // ID
	MCode          int64  // 账户代码
	AdvisorType    int    // 投顾类型：0.普通会员、1.内部员工、2.持牌投顾
	Avatar         string // 用户头像
	CreateTime     int64  // 创建时间
	FriendlyName   string // 全称（名字）
	GroupID        int64  // 会员组ID
	Intro          string // 投顾简介
	IsHelper       int    // 助手标签
	LastUpdateTime int64  // 最后更新时间
	Level          int    // 投顾等级
	MAN            string // 号码(4B+SK.AES)
	Mobile         string // 手机号
	QCer           string // 证书编号
	Status         int    // 状态
	Tags           string // 投顾标签
	DeviceToken    string // jpush token
}

type MemberBrief struct {
	ID           int64
	GroupID      int64  // 会员组ID
	DeviceTicket string // 用户硬件识别ID
	DeviceOnline int    // 移动设备标记
}

type AdvisorJson struct {
	GUID               string `json:"_id"`
	AdvisorType        int    `json:"advisor_type"` // 投顾类型
	Avatar             string `json:"avatar"`       // 用户头像
	CreateTime         int64  `json:"create_time"`  // 创建时间
	Fans               int    `json:"fans"`         // 粉丝数量
	FriendlyName       string `json:"name"`         // 全称（名字）
	FriendlyNamePinyin string `json:"name_pinyin"`  // 全称拼音（名字）
	Intro              string `json:"intro"`        // 投顾简介
	Level              string `json:"level"`        // 投顾等级
	LevelName          string `json:"level_name"`   // 投顾等级名称
	MCode              int64  `json:"uid"`          // 会员账户代码
	QCer               string `json:"qcer"`         // 证书编号
	Tags               string `json:"tags"`         // 投顾标签
}

type MemberMobile struct {
	MAN    string
	MAQ    string
	Mobile string
}

// --------------------------------------------------------------------------------

func NewMember() *Member {
	return &Member{
		Model: Model{
			CacheKey:  REDIS_SIMPLE_MEMBERS,
			TableName: TABLE_MEMBERS,
			Db:        MyCat,
		},
	}
}

// 获取用户基本信息
func (this *Member) GetSingleBrief(id int64) (*MemberBrief, error) {
	if id <= 0 {
		return nil, fmt.Errorf("Id  invaild")
	}
	var m MemberBrief
	exps := map[string]interface{}{
		"m.ID=?": id,
	}
	builder := this.Db.Select("ID,GroupID,DeviceTicket,DeviceOnline").From(this.TableName)
	err := this.SelectWhere(builder, exps).LoadStruct(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// 获取单条数据
func (this *Member) GetSingle(id int64) error {
	advisor, err := this.GetSingleAdvisor(id)
	if err != nil {
		return err
	}

	this.ID = advisor.ID
	this.MCode = advisor.MCode
	this.AdvisorType = advisor.AdvisorType
	this.Avatar = advisor.Avatar
	this.CreateTime = advisor.CreateTime
	this.FriendlyName = advisor.FriendlyName
	this.GroupID = advisor.GroupID
	this.IsHelper = advisor.IsHelper
	this.LastUpdateTime = advisor.LastUpdateTime
	this.DeviceToken = advisor.DeviceToken
	this.MAN = advisor.MAN
	this.Mobile = advisor.Mobile
	this.Status = advisor.Status
	return nil
}

// 获取单条投顾数据
func (this *Member) GetSingleAdvisor(id int64) (*Advisor, error) {
	var advisor Advisor

	cacheKey := fmt.Sprintf(this.CacheKey, id)

	rec, err := redis.Hgetall(cacheKey)
	if err == nil && len(rec) > 0 {
		if err := MapToStruct(&advisor, rec); err != nil {
			redis.Del(cacheKey)
		} else {
			return &advisor, nil
		}
	}

	exps := map[string]interface{}{
		"m.ID=?": id,
	}
	builder := this.Db.Select("m.ID,m.MCode,m.AdvisorType,m.Avatar,m.CreateTime,m.FriendlyName,m.MAN,m.Mobile,m.GroupID,m.LastUpdateTime,IFNULL(a.Intro,\"\") AS Intro,IFNULL(a.QCer,\"\") AS QCer,IFNULL(a.Tags,\"\") AS Tags,IFNULL(a.Level,0) AS Level,m.Status,m.DeviceToken").
		From(this.TableName+" AS m").
		LeftJoin(TABLE_MEMBER_ADVISORS+" AS a", "m.ID=a.MemberID")
	if err := this.SelectWhere(builder, exps).LoadStruct(&advisor); err != nil {
		return nil, err
	}
	fmt.Print(advisor)
	if err := redis.Hmset(cacheKey, StructToMap(&advisor)); err != nil {
		return nil, err
	}

	return &advisor, nil
}

func (this *Member) GetAdvisorJson(a *Advisor) (*AdvisorJson, error) {
	var jsn AdvisorJson
	if a.ID < 1 {
		return &jsn, ErrUndefinedMemberID
	}

	fan, err := redis.Get(fmt.Sprintf(REDIS_MEMBERS_FANS, a.ID))
	if err != nil && !strings.EqualFold(err.Error(), "redigo: nil returned") {
		return nil, err
	}
	fans, err := strconv.Atoi(fan)
	if err != nil {
		return nil, err
	}

	levelGroup := NewGroup()
	if a.AdvisorType == ADVISOR_TYPE_ADVISOR {
		if err := levelGroup.GetSingle(int64(a.Level), GROUP_ASSORT_ADVISOR_LEVEL); err != nil {
			return &jsn, err
		}
	}

	jsn.GUID = IDEncrypt(a.ID)
	jsn.Avatar = GetAvtar(a.ID, a.LastUpdateTime, a.Avatar)
	jsn.CreateTime = a.CreateTime
	jsn.FriendlyName = a.FriendlyName
	jsn.FriendlyNamePinyin = lib.Pinyin(a.FriendlyName)
	jsn.Intro = a.Intro
	jsn.Level = IDEncrypt(int64(a.Level))
	jsn.LevelName = levelGroup.GroupName
	jsn.MCode = a.MCode
	jsn.QCer = a.QCer
	jsn.Tags = a.Tags
	jsn.Fans = fans

	return &jsn, nil
}

func (this *Member) GetSingleJson(m *Member) (*MemberJson, error) {
	var jsn MemberJson
	if m.ID < 1 {
		return &jsn, ErrUndefinedMemberID
	}

	jsn.GUID = IDEncrypt(m.ID)
	jsn.AdvisorType = m.AdvisorType
	jsn.Avatar = GetAvtar(m.ID, m.LastUpdateTime, m.Avatar)
	jsn.CreateTime = m.CreateTime
	jsn.FriendlyName = m.FriendlyName
	jsn.FriendlyNamePinyin = lib.Pinyin(m.FriendlyName)
	jsn.MCode = m.MCode

	return &jsn, nil
}

func (this *Member) GetMobileMAN(id int64) (*MemberMobile, error) {
	var data MemberMobile

	if id <= 0 {
		return nil, fmt.Errorf("Id  invaild")
	}
	exps := map[string]interface{}{
		"ID=?": id,
	}
	builder := this.Db.Select("Mobile, MAN, MAQ").From(this.TableName)
	err := this.SelectWhere(builder, exps).LoadStruct(&data)

	return &data, err
}

func (this *Member) ResetCache(id int64) error {
	this.DelCache(id)

	return this.GetSingle(id)
}
