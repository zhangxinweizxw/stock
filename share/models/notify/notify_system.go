package notify

import (
    "fmt"
    "strconv"

/share/models"

	"stock
/share/gocraft/dbr"
	"stock
/share/logging"
	"stock
/share/models/message"
)

type NotifySystem struct {
	Model      `db:"-" `
	ID         int64  // ID
	CreateTime int64  // 创建时间
	Content    string // 通知内容
	MemberID   int64  // 创建者
	RefID      int64  // 类型ID
	RefType    int    // 类型 (1.锦囊包、2.视频聊天、3.图文直播、4.研报、5.会员组)
}

func NewNotifySystem() *NotifySystem {
	return &NotifySystem{
		Model: Model{
			TableName: TABLE_NOTIFY_SYSTEM,
			Db:        MyCat,
		},
	}
}

//
func (this *NotifySystem) GetSingle(exps map[string]interface{}) error {

	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps).
		LoadStruct(this)
	return err
}

// 获取推送提醒数据包
func (this *NotifySystem) GetSingleJson(exps map[string]interface{}) (*NotifyJson, error) {
	var jsn *NotifyJson

	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps).
		LoadStruct(this)
	if err != nil {
		return jsn, err
	}
	return this.getJson(this)
}

// 获取推送提醒数据包列表
func (this *NotifySystem) GetListJson(limit uint64, exps map[string]interface{}, conditions ...dbr.Condition) ([]*NotifyJson, error) {
	var list []*NotifySystem

	builder := this.Db.Select("*").From(this.TableName)
	_, err := this.SelectWhere(builder, exps, conditions...).
		Limit(limit).
		OrderBy("CreateTime DESC").
		LoadStructs(&list)
	if err != nil {
		return nil, err
	}

	listJson := make([]*NotifyJson, len(list))
	for i, v := range list {
		listJson[i], err = this.getJson(v)
		if err != nil {
			return nil, err
		}
	}
	return listJson, err
}

func (this *NotifySystem) getJson(notify *NotifySystem) (*NotifyJson, error) {
	var jsn NotifyJson
	mjsn, err := message.GetCahceMember(notify.MemberID)
	if err != nil && err != dbr.ErrNotFound {
		logging.Debug("Get CacheMember | %v", err)
		return &jsn, err
	}
	jsn.ID = IDEncrypt(notify.ID)
	jsn.Body.Inline = notify.Content
	jsn.CreateTime, _ = strconv.ParseFloat(fmt.Sprintf("%.6f", float64(notify.CreateTime)/1000000000), 64)
	jsn.RefType = REFTYPE_SYSTEM
	jsn.RefID = IDEncrypt(notify.RefID)
	jsn.Template = TypeMode[notify.RefType]
	jsn.Type = 4
	if mjsn != nil {
		jsn.From = *mjsn
	}
	return &jsn, nil
}
