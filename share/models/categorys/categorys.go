package categorys

/share/models"

	"stock
/share/gocraft/dbr"
)

type Categorys struct {
	Model        `db:"-"`
	ID           int64
	Assort       int    // 分类
	CategoryName string // 标题
	Color        string // 分类颜色
	ParentID     int64  // 父级ID
	IsProtect    int    // 保护标记
	Order        int    // 排序
	Tag          int    // 自定义标签(直播时： 1、视频直播、2、图文直播)
}

// 分类简短Json
type CategoryBriefJson struct {
	Assort     int    `json:"assort"`
	CategoryID string `json:"category_id"`
	Color      string `json:"color"`
	Name       string `json:"name"`
}

// 分类详情Json
type CategoryDetailJson struct {
	ID           string `json:"_id"`
	Assort       int    `json:"assort"`
	CategoryName string `json:"category_name"`
	Color        string `json:"color"`
	ParentId     string `json:"parent_id"`
	Order        int    `json:"order"`
	Tag          int    `json:"tag"`
}

type CategorysJson struct {
	ID           string `json:"_id"`
	Assort       int    `json:"assort"`
	CategoryName string `json:"category_name"`
	Color        string `json:"color"`
	Order        int    `json:"order"`
}

func NewCategorys() *Categorys {
	return &Categorys{
		Model: Model{
			TableName: TABLE_CATEGORYS,
			Db:        MyCat,
		},
	}
}

//
func (this *Categorys) GetListByAssort(assort int) ([]*Categorys, error) {
	var lst []*Categorys
	exps := map[string]interface{}{
		"Assort=?": assort,
	}

	builder := this.Db.Select("*").From(this.TableName + " AS c")
	_, err := this.SelectWhere(builder, exps).
		OrderBy("c.Order Asc").
		LoadStructs(&lst)

	return lst, err
}

// 获取分类列表
func (this *Categorys) GetCategoryList(exps map[string]interface{}) ([]Categorys, error) {
	var cat []Categorys

	builder := this.Db.Select("*").From(this.TableName)
	_, err := this.SelectWhere(builder, exps).LoadStructs(&cat)
	return cat, err
}

// 获取简短分类列表Json
func (this *Categorys) GetCategoryBriefListJson(lst []Categorys) []CategoryBriefJson {
	jsn := make([]CategoryBriefJson, len(lst))
	for i, v := range lst {
		jsn[i].Assort = v.Assort
		jsn[i].CategoryID = IDEncrypt(v.ID)
		jsn[i].Color = v.Color
		jsn[i].Name = v.CategoryName
	}
	return jsn
}

// 获取详细分类列表Json
func (this *Categorys) GetCategoryDetailListJson(lst []Categorys) []CategoryDetailJson {
	jsn := make([]CategoryDetailJson, len(lst))
	for i, v := range lst {
		jsn[i].ID = IDEncrypt(v.ID)
		jsn[i].Assort = v.Assort
		jsn[i].CategoryName = v.CategoryName
		jsn[i].Color = v.Color
		jsn[i].ParentId = IDEncrypt(v.ParentID)
		jsn[i].Order = v.Order
		jsn[i].Tag = v.Tag
	}
	return jsn
}

// 获取单条数据
func (this *Categorys) GetSingle(exps map[string]interface{}) error {
	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps).LoadStruct(this)

	return err
}

func (this *Categorys) GetListByExps(exps map[string]interface{}, limit int) ([]*Categorys, error) {
	var data []*Categorys

	builder := this.Db.Select("*").From(this.TableName)
	selectBuilder := this.SelectWhere(builder, exps)
	if limit > 0 {
		selectBuilder = selectBuilder.Limit(uint64(limit))
	}
	err := selectBuilder.OrderBy("`Order`").LoadStruct(&data)

	return data, err
}

func (this *Categorys) GetListJsonByExps(exps map[string]interface{}, limit int) ([]*CategorysJson, error) {
	data, err := this.GetListByExps(exps, limit)
	if err != nil {
		if err == dbr.ErrNotFound {
			return []*CategorysJson{}, nil
		}
		return nil, err
	}

	jsns := make([]*CategorysJson, len(data))
	for i, v := range data {
		jsn, err := this.GetJson(v)
		if err != nil {
			return nil, err
		}
		jsns[i] = jsn
	}

	return jsns, nil
}

func (this *Categorys) GetJson(c *Categorys) (*CategorysJson, error) {
	if c == nil {
		return nil, ErrParameterError
	}

	return &CategorysJson{
		ID:           IDEncrypt(c.ID),
		Assort:       c.Assort,
		CategoryName: c.CategoryName,
		Color:        c.Color,
		Order:        c.Order,
	}, nil
}

func (this *Categorys) GetCategoryNameByID(id int64) (string, error) {
	var name string

	exps := map[string]interface{}{
		"ID=?": id,
	}
	builder := this.Db.Select("CategoryName").From(this.TableName)
	err := this.SelectWhere(builder, exps).LoadStruct(&name)

	return name, err
}
