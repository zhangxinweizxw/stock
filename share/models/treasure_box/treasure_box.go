package treasure_box

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
/share/models/common"
	"stock
/share/models/notify"
	"stock
/share/store/redis"
)

type TreasureBox struct {
	Model          `db:"-"`
	ID             int64          // ID
	ApplyTo        dbr.NullString // 使用范围
	BeginTime      int64          // 启用时间
	CategoryID     int64          // 分类ID
	CreateTime     int64          // 创建时间
	CoverUrl       dbr.NullString // PC端封面url
	Description    dbr.NullString // 描述
	DeviceCoverUrl dbr.NullString // 移动端封面url
	EndTime        int64          // 过期时间
	IsDelete       int            // 删除标记
	IsRecommend    int            // 推荐标记
	MemberID       int64          // 会员ID
	Price          float32        // 价格
	RiskLevel      int            // 风险等级
	RiskTip        dbr.NullString // 风险提示
	RunStatus      int            // 运行状态
	ServiceEndTime int64          // 服务结束时间
	ServicePeriod  int            // 服务周期（天）
	Status         int            // 启用状态
	StopSellTime   int64          // 停售时间
	Specialty      dbr.NullString // 特点
	Title          string         // 标题
	Type           int            // 宝箱类型 1:锦囊 4:研报
	UpdateTime     int64          // 更新时间
}

type TreasureBoxJson struct {
	ID                string             `json:"_id"`                 // GUID
	Advisor           common.AdvisorJson `json:"advisor"`             // 投顾
	ApplyTo           string             `json:"apply_to"`            // 适用范围
	BeginTime         int64              `json:"begin_time"`          // 启用时间
	CategoryID        string             `json:"category_id"`         // 分类ID
	CoverUrl          string             `json:"cover_url"`           // PC端封面url
	CreateTime        int64              `json:"create_time"`         // 创建时间
	Description       string             `json:"description"`         // 描述
	DeviceCoverUrl    string             `json:"device_cover_url"`    // 移动端封面url
	EndTime           int64              `json:"end_time"`            // 过期时间
	IsNotify          int                `json:"is_notify"`           // 通知标志
	IsPay             int                `json:"is_pay"`              // 支付标志
	PermitContinuePay int                `json:"permit_continue_pay"` // 是否继续购买
	IsSubscribe       int                `json:"is_subscribe"`        // 是否订阅
	OrderTotal        int                `json:"order_total"`         // 订购数
	RiskTip           string             `json:"risk_tip"`            // 风险提示
	RiskLevel         int                `json:"risk_level"`          // 风险等级
	Price             float32            `json:"price"`               // 价格
	ServiceEndTime    int64              `json:"service_endtime"`     // 产品服务结束时间
	ServicePeriod     int                `json:"service_period"`      // 服务周期
	MyServiceEndTime  int64              `json:"my_service_endtime"`  // 我购买的服务结束时间
	State             int                `json:"state"`               // 运行状态
	StateV2           int                `json:"state_v2"`            // 运行状态(版本2)
	Subscribes        int                `json:"subscribes"`          // 订阅数量
	StopSellTime      int64              `json:"stopsell_time"`       // 停售时间
	Specialty         string             `json:"specialty"`           // 特点
	Title             string             `json:"title"`               // 标题
	Type              int                `json:"type"`                // 宝箱类型
	UpdateTime        int64              `json:"update_time"`         // 更新时间
	DiscountPrice     float32            `json:"discount_price"`      // 打折后价格 add by yh 20170809
}

type BoxJson struct {
	TotalBoxes int                `json:"total"` // 宝箱总数
	Rows       []*TreasureBoxJson `json:"rows"`  // 宝箱列表
}

func NewTreasureBox() *TreasureBox {
	return &TreasureBox{
		Model: Model{
			CacheKey:  REDIS_ADVISOR_BOXES,
			TableName: TABLE_TREASURE_BOX,
			Db:        MyCat,
		},
	}
}

func NewTreasureBoxTx(tx *dbr.Tx) *TreasureBox {
	return &TreasureBox{
		Model: Model{
			CacheKey:  REDIS_ADVISOR_BOXES,
			TableName: TABLE_TREASURE_BOX,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

func (this *TreasureBox) GetListByExps(exps map[string]interface{}, limit int, page int, orderBy string) ([]*TreasureBox, error) {
	var data []*TreasureBox

	builder := this.Db.Select("t.*").
		From(this.TableName+" AS t").Join(TABLE_MEMBERS+" AS m", "t.MemberID=m.ID")
	err := this.SelectWhere(builder, exps).OrderBy(orderBy).Offset(uint64((page - 1) * limit)).Limit(uint64(limit)).LoadStruct(&data)

	return data, err
}

func (this *TreasureBox) GetListJsonByExps(exps map[string]interface{}, limit int, page int, orderBy string, memberId int64) ([]*TreasureBoxJson, error) {
	data, err := this.GetListByExps(exps, limit, page, orderBy)
	if err != nil {
		if err != dbr.ErrNotFound {
			return nil, err
		}

		return []*TreasureBoxJson{}, nil
	}

	jsns := make([]*TreasureBoxJson, len(data))
	for i, v := range data {
		jsns[i], err = this.GetJson(v, memberId, false)
		if err != nil {
			return nil, err
		}
	}

	return jsns, nil
}

func (this *TreasureBox) GetBoxJsonByExps(exps map[string]interface{}, limit int, page int, orderBy string, advisorId int64, memberId int64) (*BoxJson, error) {
	total, err := this.GetCountByAdvisorId(advisorId)
	if err != nil {
		return nil, err
	}

	jsns, err := this.GetListJsonByExps(exps, limit, page, orderBy, memberId)
	if err != nil {
		return nil, err
	}

	return &BoxJson{
		TotalBoxes: total,
		Rows:       jsns,
	}, nil
}

func (this *TreasureBox) GetPageJsonListByExps(exps map[string]interface{}, page int, limit int, order string, memberId int64) (*BoxJson, error) {
	total, err := this.GetCountByExps(exps)
	if err != nil {
		return nil, err
	}

	var data []*TreasureBox
	if total > 0 {
		data, err = this.GetPageListByExps(exps, page, limit, order)
		if err != nil {
			return nil, err
		}
	}

	jsns := make([]*TreasureBoxJson, len(data))
	for i, v := range data {
		jsn, err := this.GetJson(v, memberId, true)
		if err != nil {
			return nil, err
		}

		jsns[i] = jsn
	}

	return &BoxJson{
		TotalBoxes: total,
		Rows:       jsns,
	}, nil
}

func (this *TreasureBox) GetPageListByExps(exps map[string]interface{}, page int, limit int, order string) ([]*TreasureBox, error) {
	var data []*TreasureBox

	builder := this.Db.Select("t.*").
		From(this.TableName+" AS t").
		Join(TABLE_MEMBERS+" AS m", fmt.Sprintf("t.`MemberID`=m.`ID` AND m.Status=%v", MEMBER_STATUS_NORMAL))
	builder = this.SelectWhere(builder, exps)

	if len(order) > 0 {
		builder = builder.OrderBy(order)
	}

	_, err := builder.
		Offset(uint64((page - 1) * limit)).
		Limit(uint64(limit)).
		LoadStructs(&data)

	return data, err
}

func (this *TreasureBox) GetJson(p *TreasureBox, memberId int64, withAdvisor bool) (*TreasureBoxJson, error) {
	product := common.NewProduct()
	sub, _ := redis.Get(fmt.Sprintf(REDIS_TREASURE_BOX_SUBSCRIBES, p.ID))
	subInt, _ := strconv.Atoi(sub)

	// isPaid
	isPaid := 0
	paid, err := common.NewProduct().IsPaid(p.ID, p.Type, memberId)
	if err != nil {
		return nil, err
	}
	if paid {
		isPaid = 1
	}

	// isNotify
	var isNotify int
	if notify.NewNotifySubscribe().IsSubscribe(p.Type, p.ID, memberId) {
		isNotify = 1
	}

	// url
	var deviceCoverUrl string
	var coverUrl string
	if len(p.CoverUrl.String) > 0 {
		coverUrl = AFFIX_URL + p.CoverUrl.String
	}
	if len(p.DeviceCoverUrl.String) > 0 {
		deviceCoverUrl = AFFIX_URL + p.DeviceCoverUrl.String
	}

	jsn := &TreasureBoxJson{
		ID:             IDEncrypt(int64(p.ID)),
		BeginTime:      p.BeginTime,
		CategoryID:     IDEncrypt(p.CategoryID),
		CreateTime:     p.CreateTime,
		DeviceCoverUrl: deviceCoverUrl,
		CoverUrl:       coverUrl,
		Specialty:      p.Specialty.String,
		Description:    p.Description.String,
		EndTime:        p.EndTime,
		IsNotify:       isNotify,
		IsPay:          isPaid,
		Price:          p.Price,
		RiskLevel:      p.RiskLevel,
		ServicePeriod:  p.ServicePeriod,
		State:          product.GetState(p.BeginTime, p.EndTime, time.Now().Unix(), 0),
		StateV2:        p.RunStatus,
		Subscribes:     subInt,
		StopSellTime:   p.StopSellTime,
		Title:          p.Title,
		Type:           p.Type,
		UpdateTime:     p.UpdateTime,
		DiscountPrice:  p.Price, // add by yh 20170809
	}

	if withAdvisor {
		member := common.NewMember()
		advisor, err := member.GetSingleAdvisor(p.MemberID)
		if err != nil {
			return nil, err
		}

		advisorJson, err := member.GetAdvisorJson(advisor)
		if err != nil {
			return nil, err
		}

		jsn.Advisor = *advisorJson
	}

	return jsn, nil
}

// 获取多条数据
func (this *TreasureBox) GetListJson(exps map[string]interface{}, orderBy string, limit uint64, memberID int64) ([]*TreasureBoxJson, error) {
	var treasureList []*TreasureBox

	builder := this.Db.Select("t.*").
		From(this.TableName+" AS t").Join(TABLE_MEMBERS+" AS m", fmt.Sprintf("t.MemberID=m.ID AND m.Status=%v", MEMBER_STATUS_NORMAL))
	_, err := this.SelectWhere(builder, exps).
		OrderBy(orderBy).
		LoadStructs(&treasureList)
	if err != nil {
		return nil, err
	}

	dataList := make([]*TreasureBoxJson, len(treasureList))
	for i, v := range treasureList {
		dataList[i], err = this.getSingleJson(v, memberID)
		if err != nil {
			return dataList, err
		}
	}
	return dataList, nil
}

func (this *TreasureBox) GetSingleById(id int64, typ int) error {
	exps := map[string]interface{}{
		"ID=?":   id,
		"Type=?": typ,
	}

	return this.GetSingleByExps(exps)
}

func (this *TreasureBox) GetSingleByExps(exps map[string]interface{}) error {
	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps).LoadStruct(this)

	return err
}

// 获取单条数据
func (this *TreasureBox) GetSingleJson(memberID int64, exps map[string]interface{}, conditions ...dbr.Condition) (*TreasureBoxJson, error) {
	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps, conditions...).
		LoadStruct(this)
	if err != nil {
		return nil, err
	}

	return this.getSingleJson(this, memberID)
}

// 更新数据
func (this *TreasureBox) UpdateData(params map[string]interface{}, exps map[string]interface{}) error {
	builder := this.Db.Update(this.TableName)
	this.UpdateParams(builder, params)
	r, err := this.UpdateWhere(builder, exps).Exec()
	if err != nil {
		return err
	}
	n, _ := r.RowsAffected()
	if n == 0 {
		return dbr.ErrNotFound
	}
	return nil
}

func (this *TreasureBox) GetCountByAdvisorId(advisorId int64) (int, error) {
	if advisorId == 0 {
		return 0, nil
	}

	cacheKey := fmt.Sprintf(this.CacheKey, advisorId)
	opinion, err := redis.Get(cacheKey)
	if err == nil && len(opinion) > 0 {
		boxes, err := strconv.Atoi(opinion)
		if err == nil {
			return boxes, nil
		}

		redis.Del(cacheKey)
	}

	exps := map[string]interface{}{
		"MemberID=?": advisorId,
		"IsDelete=?": 0,
		"Status=?":   TACTIC_STATUS_ENABLED,
	}
	boxes, err := this.GetCount(exps)
	if err != nil {
		return 0, err
	}

	redis.Set(cacheKey, []byte(fmt.Sprintf("%v", boxes)))

	return boxes, nil
}

func (this *TreasureBox) GetCountByExps(exps map[string]interface{}) (int, error) {
	var count int

	builder := this.Db.Select("COUNT(0)").
		From(this.TableName+" AS t").
		Join(TABLE_MEMBERS+" AS m", "t.MemberID=m.ID")
	_, err := this.SelectWhere(builder, exps).
		Limit(1).
		LoadStructs(&count)

	return count, err
}

// --------------------------------------------------------------------------------

// 获取单条数据JSON
func (this *TreasureBox) getSingleJson(t *TreasureBox, memberId int64) (*TreasureBoxJson, error) {
	var jsn TreasureBoxJson

	// Advisor
	advisor, err := common.NewMember().GetSingleAdvisor(t.MemberID)
	if err != nil {
		logging.Debug("Get Single Advisor | %v", err)
	}
	advisorJson, err := common.NewMember().GetAdvisorJson(advisor)
	if err != nil {
		logging.Debug("Get Single Advisor Json| %v", err)
	}

	// url
	var deviceCoverUrl string
	var coverUrl string
	if len(t.DeviceCoverUrl.String) > 0 {
		deviceCoverUrl = AFFIX_URL + t.DeviceCoverUrl.String
	}
	if len(t.CoverUrl.String) > 0 {
		coverUrl = AFFIX_URL + t.CoverUrl.String
	}

	product := common.NewProduct()
	nowTime := time.Now().Unix()

	jsn.ID = IDEncrypt(t.ID)
	jsn.Advisor = *advisorJson
	jsn.ApplyTo = t.ApplyTo.String
	jsn.BeginTime = t.BeginTime
	jsn.CategoryID = IDEncrypt(t.CategoryID)
	jsn.CreateTime = t.CreateTime
	jsn.DeviceCoverUrl = deviceCoverUrl
	jsn.CoverUrl = coverUrl
	jsn.Description = t.Description.String
	jsn.EndTime = t.EndTime

	paid, err := product.IsPaid(t.ID, t.Type, memberId)
	if err != nil {
		return &jsn, err
	}

	if paid {
		jsn.IsPay = 1
	} else {
		jsn.IsPay = 0
	}
	if notify.NewNotifySubscribe().IsSubscribe(t.Type, t.ID, memberId) {
		jsn.IsNotify = 1
	} else {
		jsn.IsNotify = 0
	}
	jsn.IsSubscribe = 0
	jsn.OrderTotal = 0
	jsn.Price = t.Price
	jsn.RiskTip = t.RiskTip.String
	jsn.RiskLevel = t.RiskLevel
	jsn.Specialty = t.Specialty.String
	jsn.ServicePeriod = t.ServicePeriod
	jsn.State = product.GetState(t.BeginTime, t.EndTime, nowTime, 0)
	jsn.StateV2 = t.RunStatus
	jsn.StopSellTime = t.StopSellTime
	jsn.Title = t.Title
	jsn.UpdateTime = t.UpdateTime
	jsn.DiscountPrice = t.Price // add by yh 20170809

	return &jsn, nil
}
