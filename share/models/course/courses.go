package course

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
/share/store/redis"
)

type Courses struct {
	Model       `db:"-" `
	ID          int64          // ID
	ApplyTo     dbr.NullString // 使用范围
	BeginTime   int64          // 启用时间
	CategoryID  int64          // 分类ID
	CreateTime  int64          // 创建时间
	Description dbr.NullString // 描述
	EndTime     int64          // 过期时间
	Extra1      dbr.NullString // 扩展字段1
	Extra2      dbr.NullString // 扩展字段2
	Extra3      dbr.NullString // 扩展字段3
	Extra4      dbr.NullString // 扩展字段4
	IsDelete    int            // 删除标记
	IsRecommend int            // 推荐标记
	MemberID    int64          // 会员ID
	Price       float32        // 价格
	RiskLevel   int            // 风险等级
	RiskTip     dbr.NullString // 风险提示
	Specialty   dbr.NullString // 特点
	Status      int            //
	Title       string         // 标题
	Type        int            // 宝箱类型 1:锦囊 4:研报
	UpdateTime  int64          // 更新时间

	//课堂改版添加字段 wdk 20170810 add
	RunStatus      int   // 运行状态
	ServiceEndTime int64 // 服务结束时间
	ServicePeriod  int   // 服务周期（天）
	StopSellTime   int64 // 停售时间
}

// 已购课程简要信息
type BuyCourseBrief struct {
	ID         int64
	Title      string
	BuyTime    int64
	ExpireTime int64
	RunStatus  int
	AdvisorID  int64
}

type BuyCourseBriefJson struct {
	ID         string             `json:"_id"`
	Title      string             `json:"title"`
	BuyTime    int64              `json:"buy_time"`
	ExpireTime int64              `json:"expire_time"`
	RunStatus  int                `json:"run_status"`
	BuyCount   int                `json:"buy_count"` // zxw 20170815 add 购买人数
	Advisor    common.AdvisorJson `json:"advisor"`
}

func NewCourses() *Courses {
	return &Courses{
		Model: Model{
			TableName: TABLE_COURSES,
			Db:        MyCat,
		},
	}
}

func NewCoursesTx(tx *dbr.Tx) *Courses {
	return &Courses{
		Model: Model{
			TableName: TABLE_COURSES,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

//
func (this *Courses) GetSingle(exps map[string]interface{}) error {

	builder := this.Db.Select("*").
		From(this.TableName)
	err := this.SelectWhere(builder, exps).
		LoadStruct(this)

	return err
}

//
func (this *Courses) GetStateByID(id int64) int {

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

// 获取已购课程列表
func (this *Courses) GetBuyListJson(mID int64) ([]BuyCourseBriefJson, error) {

	// 获取已购的课堂ID
	key := fmt.Sprintf(REDIS_MEMBER_ORDER_ACTIVITY, mID)
	values, err := redis.Smembers(key)
	if err != nil {
		logging.Debug("Redis Get Buy Course List | %v", err)
		return nil, err
	}
	ids := make([]int64, len(values))
	for i, v := range values {
		ids[i], _ = strconv.ParseInt(v, 10, 64)
	}
	if len(ids) == 0 {
		ids = append(ids, 0)
	}

	// 课堂简要信息
	exps := map[string]interface{}{
		"ID IN ?": ids,
	}
	var buyList []*BuyCourseBrief
	builder := this.Db.Select("ID, Title, RunStatus, MemberID AS AdvisorID").
		From(this.TableName).
		OrderBy("RunStatus ASC, CreateTime DESC")
	_, err = this.SelectWhere(builder, exps).Load(&buyList)
	if err != nil {
		return nil, err
	}

	data := make([]BuyCourseBriefJson, len(buyList))
	for i, v := range buyList {

		// 投顾信息
		advisor, err := common.NewMember().GetSingleAdvisor(v.AdvisorID)
		if err != nil {
			logging.Debug("Get Single Advisor | %v", err)
		}
		advisorJson, err := common.NewMember().GetAdvisorJson(advisor)
		if err != nil {
			logging.Debug("Get Single AdvisorJson | %v", err)
		}
		data[i].Advisor = *advisorJson

		// 下单时间、服务结束时间
		pmbExps := map[string]interface{}{
			"p.RefID=?":      v.ID,
			"p.RefType=?":    REFTYPE_COURSE,
			"pmb.MemberID=?": mID,
			"pmb.IsDelete=?": 0,
			"p.IsDelete=?":   0,
		}
		pmbData, err := common.NewProductMemberBuy().GetList(pmbExps, 1, "pmb.ServiceEndTime DESC")
		if err != nil {
			logging.Debug("Get Product Member Buy List | %v", err)
		}
		data[i].BuyTime = pmbData[0].CreateTime
		data[i].ExpireTime = pmbData[0].ServiceEndTime

		data[i].ID = IDEncrypt(v.ID)
		data[i].Title = v.Title
		data[i].RunStatus = v.RunStatus
		data[i].BuyCount = this.GetBuyCount(v.ID)
	}

	return data, nil
}

//课堂改版 获取课堂状态从新的字段RunStatus中得到 wdk 20170810 add
func (this *Courses) GetStateByIDNew(id int64) (int, error) {
	exps := map[string]interface{}{
		"ID=?": id,
	}
	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps).LoadStruct(this)
	if err != nil {
		return 0, err
	}

	return this.RunStatus, nil
}

// 购买数量 zxw 20170815 add
func (this *Courses) GetBuyCount(courseID int64) int {

	// 读缓存
	key := fmt.Sprintf(REDIS_COURSE_BUY_COUNT_CRIME, courseID)
	value, err := redis.Get(key)
	if err == nil {
		count, err := strconv.Atoi(value)
		if err == nil {
			return count
		}
	}

	return int(0)
}
