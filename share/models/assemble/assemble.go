package assemble

import (
    "time"

/share/models"

	"stock
/share/gocraft/dbr"
	"stock
/share/logging"
	"stock
/share/models/common"
)

type Assemble struct {
	Model          `db:"-"`
	ID             int64   // GUID
	BeginTime      int64   // 开始时间
	CreateTime     int64   // 创建时间
	CourseID       int64   // 课程ID
	DeviceImageUrl string  // 移动端图片URL
	DeviceLinkUrl  string  // 移动端链接URL
	EndTime        int64   // 结束时间
	IsDelete       int     // 删除标识
	ImageUrl       string  // PC端图片URL
	LinkUrl        string  // PC端链接URL
	MemberID       int64   // 用户ID
	Order          int     // 排序
	Price          float32 // 价格
	RiskLevel      int     // 风险等级
	ReportID       int64   // 内参ID
	Status         int     // 状态（0.禁用、1.启用、2.停购）
	TacticID       int64   // 锦囊ID
	UpdateTime     int64   // 更新时间
}

type AssembleJson struct {
	ID             string             `json:"_id"`
	BeginTime      int64              `json:"begin_time"`
	CreateTime     int64              `json:"create_time"`
	CourseID       string             `json:"course_id"`
	DeviceImageUrl string             `json:"device_image_url"`
	DeviceLinkUrl  string             `json:"device_link_url"`
	EndTime        int64              `json:"end_time"`
	ImageUrl       string             `json:"image_url"`
	LinkUrl        string             `json:"link_url"`
	MemberID       string             `json:"member_id"`
	Price          float32            `json:"price"`
	RiskLevel      int                `json:"risk_level"`
	ReportID       string             `json:"report_id"`
	TacticID       string             `json:"tactic_id"`
	Advisor        common.AdvisorJson `json:"advisor"`
	DiscountPrice  float32            `json:"discount_price"` // add by yh 20170807(打折后价格)
}

type AssembleIsPaidJson struct { //wdk 20170802 add
	ID     string `json:"ref_id"`
	IsPaid bool   `json:"is_paid"`
}

func NewAssemble() *Assemble {
	return &Assemble{
		Model: Model{
			TableName: TABLE_ASSEMBLE,
			Db:        MyCat,
		},
	}
}

func NewAssembleTx(tx *dbr.Tx) *Assemble {
	return &Assemble{
		Model: Model{
			TableName: TABLE_ASSEMBLE,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

func (this *Assemble) GetSingle(exps map[string]interface{}) error {
	builder := this.Db.Select("*").
		From(this.TableName)
	err := this.SelectWhere(builder, exps).
		LoadStruct(this)
	return err
}

// 获取复合产品列表
func (this *Assemble) GetListJson(exps map[string]interface{}, orderBy string) ([]*AssembleJson, error) {
	var assembleList []*Assemble

	builder := this.Db.Select("*").
		From(this.TableName)
	_, err := this.SelectWhere(builder, exps).
		OrderBy(orderBy).
		LoadStructs(&assembleList)
	if err != nil {
		return nil, err
	}

	data := make([]*AssembleJson, len(assembleList))
	for i, v := range assembleList {
		temp, err := this.getJson(v)
		if err != nil {
			logging.Debug("Get Assemble Single Json | %v", err)
			continue
		}
		data[i] = temp
	}

	return data, nil
}

// 获取一个复合产品
func (this *Assemble) GetSingleJson(exps map[string]interface{}) (*AssembleJson, error) {
	var assemble Assemble

	builder := this.Db.Select("*").
		From(this.TableName)
	_, err := this.SelectWhere(builder, exps).
		Limit(1).
		LoadStructs(&assemble)
	if err != nil {
		return nil, err
	}

	data, err := this.getJson(&assemble)
	if err != nil {
		logging.Debug("Get Assemble Single Json | %v", err)
		return nil, err
	}

	return data, nil
}

//
func (this *Assemble) GetStateByID(id int64) int {

	exps := map[string]interface{}{
		"ID=?": id,
	}
	builder := this.Db.Select("*").
		From(this.TableName)
	this.SelectWhere(builder, exps).
		LoadStruct(this)

	var state int
	now := time.Now().Unix()
	if this.Status == TREASUREBOX_STATUS_STOP_SELL {
		state = TREASURE_RUN_STATE_STOP_SELL
	} else if now < this.BeginTime {
		state = TREASURE_RUN_STATE_BEFORE_SELLING
	} else if now >= this.BeginTime && now <= this.EndTime {
		state = TREASURE_RUN_STATE_RUNNING
	} else {
		state = TREASURE_RUN_STATE_ENDED
	}

	return state
}

func (this *Assemble) getJson(a *Assemble) (*AssembleJson, error) {
	var data AssembleJson

	// Advisor
	advisor, err := common.NewMember().GetSingleAdvisor(a.MemberID)
	if err != nil {
		return nil, err
	}
	advisorJson, err := common.NewMember().GetAdvisorJson(advisor)
	if err != nil {
		return nil, err
	}

	data.ID = IDEncrypt(a.ID)
	data.BeginTime = a.BeginTime
	data.CreateTime = a.CreateTime
	data.CourseID = IDEncrypt(a.CourseID)
	data.DeviceImageUrl = a.DeviceImageUrl
	data.DeviceLinkUrl = a.DeviceLinkUrl
	data.EndTime = a.EndTime
	data.ImageUrl = a.ImageUrl
	data.LinkUrl = a.LinkUrl
	data.MemberID = IDEncrypt(a.MemberID)
	data.Price = a.Price
	data.RiskLevel = a.RiskLevel
	data.ReportID = IDEncrypt(a.ReportID)
	data.TacticID = IDEncrypt(a.TacticID)
	data.Advisor = *advisorJson
	data.DiscountPrice = data.Price

	return &data, nil
}
