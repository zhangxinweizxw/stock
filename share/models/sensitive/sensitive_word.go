package sensitive

/share/models"
)

type SensitiveWord struct {
	Model      `db:"-" `
	ID         int64  // ID
	Content    string // 内容
	IsDelete   int    // 删除标记
	UpdateTime int64  // 更新时间
}

func NewSensitiveWord() *SensitiveWord {
	return &SensitiveWord{
		Model: Model{
			TableName: TABLE_SENSITIVE_WORD_MSG, //wdk 20170713 modify （Andy说以后会有四种敏感词库，现在这里分开）
			Db:        MyCat,
		},
	}
}

// 获取内容列表
func (this *SensitiveWord) GetContentList() ([]string, error) {
	var data []string

	dbr := this.Db.Select("Content").From(this.TableName)
	_, err := this.SelectWhere(dbr, map[string]interface{}{}).LoadStructs(&data)

	return data, err
}
