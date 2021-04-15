package upload

import (
    "fmt"

/share/models"

	"stock
/share/gocraft/dbr"
	"stock
/share/store/redis"
)

type UploadAffixs struct {
	Model         `db:"-"`
	ID            int64
	Assort        int    // 附件分类
	Bucket        string // Aliyun Bucket
	CreateTime    int64  // 创建时间
	Creator       int64  // 创建人
	Duration      int    // 语音时长
	FileExt       string // 扩展名
	FileName      string // 源文件名称
	FilePath      string // 文件保存路径
	FileSize      int64  // 大小
	Height        int    // 图片高度
	RelevanceID   int64  // 附件关联ID
	RelevanceTime int64  // 创建时间
	Thumbnail     int    // 缩略图标记
	Title         string // 文件标题
	Width         int    // 图片宽度
	IsPay         int    // 支付标记
	RefType       int    // 类型
	RefID         int64  // 类型ID
}

type UploadAttachmentJson struct {
	GUID          string `json:"_id"`
	Assort        int    `json:"assort"`
	CreateTime    int64  `json:"created_time"`
	Creator       string `json:"creator"`
	Duration      int    `json:"duration"`
	Ext           string `json:"ext"`
	Height        int    `json:"height"`
	Relevance     string `json:"relevance"`
	RelevanceTime int64  `json:"relevance_time"`
	Size          int64  `json:"size"`
	Title         string `json:"title"`
	Width         int    `json:"width"`
}

// --------------------------------------------------------------------------------

func NewUploadAffixs() *UploadAffixs {
	return &UploadAffixs{
		Model: Model{
			CacheKey:  REDIS_UPLOAD_AFFIXS,
			TableName: TABLE_UPLOAD_AFFIXS,
			Db:        MyCat,
		},
	}
}

func NewUploadAffixsTx(tx *dbr.Tx) *UploadAffixs {
	return &UploadAffixs{
		Model: Model{
			CacheKey:  REDIS_UPLOAD_AFFIXS,
			TableName: TABLE_UPLOAD_AFFIXS,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

func (this *UploadAffixs) GetSingle(id int64) error {
	cacheKey := fmt.Sprintf(this.CacheKey, id)
	rec, err := redis.Hgetall(cacheKey)
	if err == nil && len(rec) > 1 {
		if err := MapToStruct(this, rec); err != nil {
			redis.Del(cacheKey)
		} else {
			//return nil
		}
	}
	exps := map[string]interface{}{
		"a.ID=?": id,
	}
	builder := this.Db.Select("a.*, r.ID AS RelevanceID, r.CreateTime AS RelevanceTime,r.IsPay as IsPay,r.RefType as RefType,r.RefID as RefID").From(this.TableName+" AS a").
		Join(TABLE_UPLOAD_RELEVANCE+" AS r", "a.ID=r.AffixID")
	err = this.SelectWhere(builder, exps).
		Limit(1).
		LoadStruct(&this)
	if err != nil {
		return err
	}

	err = redis.Hmset(cacheKey, StructToMap(this))
	if err != nil {
		return err
	}

	// 设置过期时间（时效15天）
	_, err = redis.Do("EXPIRE", cacheKey, 60*60*24*15)
	return err
}

func (this *UploadAffixs) GetListByRefId(refId int64, refType int, limit int, latestStamp int64) ([]UploadAffixs, error) {
	exps := map[string]interface{}{
		"r.RefID=?":   refId,
		"r.RefType=?": refType,
	}

	return this.GetListByExps(exps, limit, latestStamp)
}

func (this *UploadAffixs) GetListByExps(exps map[string]interface{}, limit int, latestStamp int64) ([]UploadAffixs, error) {
	var data []UploadAffixs

	if latestStamp > 0 {
		exps["a.ID<=?"] = latestStamp
	}

	builder := this.Db.Select("a.*, r.ID AS RelevanceID, r.CreateTime AS RelevanceTime").From(this.TableName+" AS a").
		Join(TABLE_UPLOAD_RELEVANCE+" AS r", "a.ID=r.AffixID")
	_, err := this.SelectWhere(builder, exps).
		OrderBy("r.ID DESC").
		Limit(uint64(limit + 1)).
		LoadStructs(&data)
	return data, err
}

func (this *UploadAffixs) GetSingleJson(u *UploadAffixs) (UploadAttachmentJson, error) {
	var jsn UploadAttachmentJson

	jsn.GUID = IDEncrypt(u.ID)
	jsn.Assort = u.Assort
	jsn.CreateTime = u.CreateTime
	jsn.Creator = IDEncrypt(u.Creator)
	jsn.Duration = u.Duration
	jsn.Ext = u.FileExt
	jsn.Height = u.Height
	jsn.Relevance = IDEncrypt(u.RelevanceID)
	jsn.RelevanceTime = u.RelevanceTime
	jsn.Size = u.FileSize
	jsn.Title = u.Title
	jsn.Width = u.Width
	return jsn, nil
}
