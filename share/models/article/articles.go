package article

/share/models"

	"stock
/share/gocraft/dbr"
)

type Articles struct {
	Model        `db:"-"`
	ID           int64
	Affixs       string
	Author       string
	CategoryID   int64
	Content      dbr.NullString
	CreateTime   int64
	Introduction dbr.NullString
	IsRecommend  int
	Pubdate      int64
	Source       string
	Tags         int
	Thumbnail    string
	Title        string
	ViewCount    int //wdk 20170731 add
}

// 资讯简要Json
type ArticleBriefJson struct {
	ID           string `json:"_id"`
	Author       string `json:"author"`
	CategoryID   string `json:"category_id"`
	Introduction string `json:"introduction"`
	Pubdate      int64  `json:"pubdate"`
	Source       string `json:"source"`
	Tags         int    `json:"tags"`
	Thumbnail    string `json:"thumbnail"`
	Title        string `json:"title"`
}

// 资讯详情Json
type ArticleDetailJson struct {
	ID           string `json:"_id"`
	Affixs       string `json:"affixs"`
	Author       string `json:"author"`
	CategoryID   string `json:"category_id"`
	Content      string `json:"content"`
	Introduction string `json:"introduction"`
	Pubdate      int64  `json:"pubdate"`
	Source       string `json:"source"`
	Tags         int    `json:"tags"`
	Thumbnail    string `json:"thumbnail"`
	Title        string `json:"title"`
}

func NewArticles() *Articles {
	return &Articles{
		Model: Model{
			TableName: TABLE_ARTICLES,
			Db:        MyCat,
		},
	}
}

//
func (this *Articles) GetSingle(exps map[string]interface{}) error {
	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps).
		LoadStruct(this)
	return err
}

// 获取资讯简要
func (this *Articles) GetArticleBriefListJson(exps map[string]interface{}, limit int) ([]*ArticleBriefJson, error) {
	var lst []*Articles

	builder := this.Db.Select("*").From(this.TableName)
	_, err := this.SelectWhere(builder, exps).
		OrderBy("Pubdate DESC").
		Limit(uint64(limit)).
		LoadStructs(&lst)
	if err != nil {
		return nil, err
	}

	data := make([]*ArticleBriefJson, len(lst))
	for i, v := range lst {
		data[i] = this.getBriefJson(v)
	}
	return data, nil
}

// 获取资讯详情
func (this *Articles) GetArticleDetailJson(exps map[string]interface{}) (*ArticleDetailJson, error) {
	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps).
		LoadStruct(this)
	if err != nil {
		return nil, err
	}

	data := this.getDetailJson(this)
	return data, nil
}

func (this *Articles) getBriefJson(a *Articles) *ArticleBriefJson {
	var thumbnailUrl string
	if a.Thumbnail != "" {
		thumbnailUrl = AFFIX_URL + a.Thumbnail
	}

	data := &ArticleBriefJson{
		ID:           IDEncrypt(a.ID),
		Author:       a.Author,
		CategoryID:   IDEncrypt(a.CategoryID),
		Introduction: a.Introduction.String,
		Pubdate:      a.Pubdate,
		Source:       a.Source,
		Tags:         a.Tags,
		Thumbnail:    thumbnailUrl,
		Title:        a.Title,
	}

	return data
}

func (this *Articles) getDetailJson(a *Articles) *ArticleDetailJson {
	var affixUrl string
	var thumbnailUrl string
	if a.Affixs != "" {
		affixUrl = AFFIX_URL + a.Affixs
	}
	if a.Thumbnail != "" {
		thumbnailUrl = AFFIX_URL + a.Thumbnail
	}

	data := &ArticleDetailJson{
		ID:           IDEncrypt(a.ID),
		Affixs:       affixUrl,
		Author:       a.Author,
		CategoryID:   IDEncrypt(a.CategoryID),
		Content:      a.Content.String,
		Introduction: a.Introduction.String,
		Pubdate:      a.Pubdate,
		Source:       a.Source,
		Tags:         a.Tags,
		Thumbnail:    thumbnailUrl,
		Title:        a.Title,
	}

	return data
}

//咨询访问量更新到mysql wdk 20170731 add
func (this *Articles) updateViewCount(params map[string]interface{}, exps map[string]interface{}, conditions ...dbr.Condition) error {
	builder := this.Db.Update(this.TableName)

	this.UpdateParams(builder, params)

	_, err := this.UpdateWhere(builder, exps, conditions...).Exec()
	return err
}
