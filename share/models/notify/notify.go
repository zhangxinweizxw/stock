package notify

import (
    "fmt"
    "strconv"
    "strings"

/share/models"

	"encoding/json"

	"stock
/share/gocraft/dbr"
	"stock
/share/logging"
	"stock
/share/models/common"
	"stock
/share/models/message"
)

type Notify struct {
	Model      `db:"-" `
	ID         int64  // ID
	CreateTime int64  // 创建时间
	Content    string // 通知内容
	MemberID   int64  // 创建者
	RefID      int64  // 类型ID
	RefType    int    // 类型 (1.锦囊包、2.视频聊天、3.图文直播、4.研报、5.会员组)
}

// 业务模板
var TypeMode = map[int]string{
	REFTYPE_TACTIC:   "model_treasure",
	REFTYPE_REPORT:   "model_treasure",
	REFTYPE_TEXTLIVE: "model_weblive",
	REFTYPE_WEBLIVE:  "model_weblive",
}

// 推送提醒数据包
type NotifyJson struct {
	ID         string            `json:"_id"`
	Body       NotifyBody        `json:"body"`
	CreateTime float64           `json:"create_time"`
	RefID      string            `json:"ref_id"`
	RefType    int               `json:"ref_type"`
	Template   string            `json:"template"`
	Type       int               `json:"type"` // 提醒类型 2.文本、3.附件、4.内联
	From       common.MemberJson `json:"from"`
}

type NotifyBody struct {
	Attachment interface{} `json:"attachment"`
	Content    string      `json:"content"`
	Inline     string      `json:"inline"`
}

func NewNotify() *Notify {
	return &Notify{
		Model: Model{
			TableName: TABLE_NOTIFY,
			Db:        MyCat,
		},
	}
}

//
func (this *Notify) GetSingle(exps map[string]interface{}) error {
	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps).
		LoadStruct(this)

	return err
}

// 获取推送提醒数据包
func (this *Notify) GetNotifyJson(exps map[string]interface{}) (*NotifyJson, error) {
	var jsn *NotifyJson

	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps).
		LoadStruct(this)
	if err != nil {
		return jsn, err
	}
	return this.GetJson(this)
}

// 获取推送提醒数据包列表
func (this *Notify) GetNotifyListJson(mID int64, limit int, orderBy string, whereClause string) ([]*NotifyJson, error) {
	var list []*Notify

	// 获取关注的老师ID列表
	ids, _ := common.NewMemberFollow().GetFollowIdsByMemberId(mID)
	idsString := make([]string, len(ids))
	for i, v := range ids {
		idsString[i] = strconv.FormatInt(v, 10)
	}
	commaString := strings.Join(idsString, ",")
	if len(commaString) == 0 {
		commaString = "0"
	}

	sql := fmt.Sprintf(`SELECT a.* FROM (
							SELECT n.* FROM hn_notify AS n WHERE n.RefType = %v AND n.RefID = %v
							UNION
							SELECT n.* FROM hn_notify AS n JOIN hn_notify_subscribe AS s ON n.RefID = s.RefID AND n.RefType = s.RefType WHERE s.MemberID = %v
							UNION 
							SELECT n.* FROM hn_notify AS n WHERE n.RefType = %v AND n.MemberID IN (%v)
							) AS a %v ORDER BY %v LIMIT %v `, REFTYPE_PRIVATE_NOTIFY, mID, mID, REFTYPE_MEMBER_FOLLOW, commaString, whereClause, orderBy, limit)
	_, err := this.Db.SelectBySql(sql).
		LoadStructs(&list)
	if err != nil {
		return nil, err
	}

	listJson := make([]*NotifyJson, len(list))
	for i, v := range list {
		listJson[i], err = this.GetJson(v)
		if err != nil {
			return nil, err
		}
	}
	return listJson, err
}

// 根据目标获取推送提醒数据包列表
func (this *Notify) GetNotifyListJsonByTarget(limit int, exps map[string]interface{}, orderBy string) ([]*NotifyJson, error) {
	var list []*Notify
	var ( //wdk 20170802 add 临时改动。为了配合移动端线上环境咨询推送里面早报晚报周报月报没有绘制图片的bug
		nm        NewsModel
		isReplace bool
	)

	builder := this.Db.Select("*").From(this.TableName)
	_, err := this.SelectWhere(builder, exps).
		Limit(uint64(limit)).
		OrderBy(orderBy).
		LoadStructs(&list)
	if err != nil {
		return nil, err
	}

	listJson := make([]*NotifyJson, len(list))
	for i, v := range list {

		//wdk 20170802 add -begin 临时改动。为了配合移动端线上环境咨询推送里面早报晚报周报月报没有绘制图片的bug
		if err := json.Unmarshal([]byte(v.Content), &nm); err != nil {
			logging.Debug("GetNotifyListJsonByTarget Json Unmarshal to NewsModel | %v", err)
		} else {
			if nm.CategoryName == "晨报" {
				nm.Thumbnail = "http://attach.0606.com.cn/magazine/2ca4d44c-61e4-4d8a-9f09-3fe19ece9fee.png"
				isReplace = true
			} else if nm.CategoryName == "晚报" {
				nm.Thumbnail = "http://attach.0606.com.cn/magazine/e92d7f3a-76e5-4046-b020-30c68deb1dc9.png"
				isReplace = true
			} else if nm.CategoryName == "周报" {
				nm.Thumbnail = "http://attach.0606.com.cn/magazine/16238f8d-25fa-4f7e-bfd1-9921d419a47d.png"
				isReplace = true
			} else if nm.CategoryName == "月报" {
				nm.Thumbnail = "http://attach.0606.com.cn/magazine/c8171125-c8fe-4f69-bd8c-6fda5854aa25.png"
				isReplace = true
			}
		}
		if isReplace {
			content, err := json.Marshal(nm)
			if err != nil {
				logging.Debug("Json Marshal | %v", err)
			}
			v.Content = string(content)
		}
		//wdk 20170802 add -end

		listJson[i], err = this.GetJson(v)
		if err != nil {
			return nil, err
		}
	}
	return listJson, err
}

// 获取推送消息未读数
func (this *Notify) GetUnreadCount(exps map[string]interface{}) (int64, error) {
	var count int64

	builder := this.Db.Select("COUNT(n.ID)").From(this.TableName+" AS n").
		Join(TABLE_NOTIFY_SUBSCRIBE+" AS s", "n.RefType=s.RefType AND n.RefID=s.RefID")
	err := this.SelectWhere(builder, exps).
		LoadStruct(&count)
	return count, err
}

func (this *Notify) GetJson(notify *Notify) (*NotifyJson, error) {
	var jsn NotifyJson
	mjsn, err := message.GetCahceMember(notify.MemberID)
	if err != nil && err != dbr.ErrNotFound {
		logging.Debug("Get CacheMember | %v", err)
		return &jsn, err
	}
	jsn.ID = IDEncrypt(notify.ID)
	jsn.Body.Inline = notify.Content
	jsn.CreateTime, _ = strconv.ParseFloat(fmt.Sprintf("%.6f", float64(notify.CreateTime)/1000000000), 64)
	jsn.RefType = notify.RefType
	jsn.RefID = IDEncrypt(notify.RefID)
	jsn.Template = TypeMode[notify.RefType]
	jsn.Type = 4
	if mjsn != nil {
		jsn.From = *mjsn
	}
	return &jsn, nil
}
