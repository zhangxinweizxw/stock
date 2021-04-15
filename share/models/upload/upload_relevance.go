package upload

/share/models"

	"stock
/share/gocraft/dbr"
)

type UploadRelevance struct {
	Model   `db:"-"`
	ID      int64
	AffixID int64 // 附件ID
	DataID  int64 // 数据ID
	RefID   int64 // 关联业务ID
	RefType int   // 关联业务类型
}

type UploadRelevanceJson struct {
	GUID    string `json:"_id"`
	RefID   string `json:"ref_id"`
	RefType int    `json:"ref_type"`
}

// --------------------------------------------------------------------------------

func NewUploadRelevance() *UploadRelevance {
	return &UploadRelevance{
		Model: Model{
			Db:        MyCat,
			TableName: TABLE_UPLOAD_RELEVANCE,
		},
	}
}

func NewUploadRelevanceTx(tx *dbr.Tx) *UploadRelevance {
	return &UploadRelevance{
		Model: Model{
			Db:        MyCat,
			TableName: TABLE_UPLOAD_RELEVANCE,
			Tx:        tx,
		},
	}
}

func (this *UploadRelevance) GetSingle(id int64) error {
	exps := map[string]interface{}{
		"ID=?": id,
	}
	builder := this.Db.Select("*").From(this.TableName)
	err := this.SelectWhere(builder, exps).
		Limit(1).
		LoadStruct(&this)
	return err
}

func (this *UploadRelevance) GetRelevanceIdByExps(exps map[string]interface{}) (int64, error) {
	var relevanceId int64
	builder := this.Db.Select("ID").From(this.TableName)
	err := this.SelectWhere(builder, exps).
		Limit(1).
		LoadValue(&relevanceId)
	return relevanceId, err
}

func (this *UploadRelevance) GetTotalByAffixId(affixId int64) int {
	var count int

	exps := map[string]interface{}{
		"AffixID=?": affixId,
	}

	builder := this.Db.Select("COUNT(0)").From(this.TableName)
	this.SelectWhere(builder, exps).
		Limit(1).
		LoadValue(&count)
	return count
}
