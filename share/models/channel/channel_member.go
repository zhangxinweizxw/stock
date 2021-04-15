package channel

import (
    "fmt"

/share/models"

	"stock
/share/gocraft/dbr"
	"stock
/share/store/redis"
)

type ChannelMember struct {
	Model      `db:"-"`
	ID         int64
	MemberID   int64 // 成员ID
	ChannelID  int64 // 频道ID
	CreateTime int64 // 更新时间
}

// --------------------------------------------------------------------------------

func NewChannelMember() *ChannelMember {
	return &ChannelMember{
		Model: Model{
			CacheKey:  REDIS_CHANNELS_MEMBERS,
			Db:        MyCat,
			TableName: TABLE_CHANNEL_MEMBERS,
		},
	}
}

func NewChannelMemberTx(tx *dbr.Tx) *ChannelMember {
	return &ChannelMember{
		Model: Model{
			CacheKey:  REDIS_CHANNELS_MEMBERS,
			Db:        MyCat,
			TableName: TABLE_CHANNEL_MEMBERS,
			Tx:        tx,
		},
	}
}

func (this *ChannelMember) DelCache(id int64, memberId int64) {
	redis.Del(fmt.Sprintf(this.CacheKey, id, memberId))
}

func (this *ChannelMember) GetMemberIdsByExps(exps map[string]interface{}, conditions ...dbr.Condition) ([]int64, error) {
	var ids []int64
	builder := this.Db.Select("MemberID").From(this.TableName)
	err := this.SelectWhere(builder, exps, conditions...).LoadValue(&ids)
	return ids, err
}

func (this *ChannelMember) GetSingleByRefID(id int64, memberId int64) error {
	cacheKey := fmt.Sprintf(this.CacheKey, id, memberId)
	rec, err := redis.Hgetall(cacheKey)
	if err == nil && len(rec) > 0 {
		if err := MapToStruct(this, rec); err != nil {
			redis.Del(cacheKey)
		} else {
			return nil
		}
	}

	exps := map[string]interface{}{
		"ChannelID=?": id,
		"MemberID=?":  memberId,
	}
	builder := this.Db.Select("*").From(this.TableName)
	err = this.SelectWhere(builder, exps).
		Limit(1).
		LoadStruct(&this)
	if err != nil {
		return err
	}

	return redis.Hmset(cacheKey, StructToMap(this))
}
