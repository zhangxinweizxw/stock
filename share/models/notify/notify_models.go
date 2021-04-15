package notify

// 百宝箱推送模板
type TreasureModel struct {
	Title string `json:"title"` // 标题
	Intro string `json:"intro"` // 简介
}

// 视频推送模板
type WebliveModel struct {
	Assort    int    `json:"assort"`    // 直播类型
	Content   string `json:"content"`   // 提醒内容
	Topic     string `json:"topic"`     // 直播主题
	Thumbnail string `json:"thumbnail"` // avatar(预留)
}

// 系统推送模板
type SystemModel struct {
	Assort  int    `json:"assort"`  // 分类 (预留)
	Affix   string `json:"affix"`   // 图片附件
	Title   string `json:"title"`   // 标题
	Content string `json:"content"` // 内容
}

// 资讯推送模板
type NewsModel struct {
	CategoryID   string `json:"category_id"`   // 分类ID
	CategoryName string `json:"category_name"` // 分类名称
	Title        string `json:"title"`         // 标题
	Intro        string `json:"intro"`         // 简介
	Thumbnail    string `json:"thumbnail"`     // avatar
	Link         string `json:"link"`          // 链接地址
}

// 策略模板
type StrategyModel struct {
	StrategyID string `json:"strategy_id"` // 策略ID
	Assort     int    `json:"assort"`      // 类型
	Title      string `json:"title"`       // 标题
	Intro      string `json:"intro"`       // 简介
}

// 产品成立模板
type ProductCreateModel struct {
	Title   string `json:"title"`    // 标题
	Intro   string `json:"intro"`    // 简介
	RefID   string `json:"ref_id"`   // 参照ID
	RefType int    `json:"ref_type"` // 参照类型
}

// 私信提醒模板
type PrivateModel struct {
	Assort int           `json:"assort"` // 类型 (1订阅成功 2服务到期前五天(有其他产品) 3服务到期前五天 4产品到期(有其他产品) 5产品到期)
	Title  string        `json:"title"`  // 标题
	Intro  string        `json:"intro"`  // 简介
	Item   []PrivateItem `json:"item"`
}

type PrivateItem struct {
	RefID   string `json:"ref_id"`   // 参照ID
	RefType int    `json:"ref_type"` // 参照类型
}
