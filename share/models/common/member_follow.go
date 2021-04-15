package common

import (
    "fmt"
    "strconv"

/share/models"

	"stock
/share/gocraft/dbr"
	"stock
/share/store/redis"
)

type PageFollowJson struct {
	Pagination FollowPagination    `json:"pagination"`
	Rows       []*MemberFollowJson `json:"rows"`
}

type FollowPagination struct {
	Total    int `json:"total"`
	PageSize int `json:"page_size"`
}

type MemberFollow struct {
	Model        `db:"-" `
	ID           int64  // GUID
	CreateTime   int64  // 创建时间
	FriendlyName string // 会员名称
	MemberID     int64  // 会员ID
	RefID        int64  // 关注用户ID
}

type MemberFollowJson struct {
	Advisor    AdvisorJson `json:"advisor"`
	CreateTime int64       `json:"create_time"`
	LevelName  string      `json:"level_name"`
}

type FollowAdvisorInfo struct {
	MemberID  int64
	Intro     string
	LevelName dbr.NullString
}

type FollowAndFansCount struct {
	FollowIds []string `json:"follow_ids"`
	Fans      int      `json:"fans"`
}

func NewMemberFollow() *MemberFollow {
	return &MemberFollow{
		Model: Model{
			CacheKey:  REDIS_MEMBERS_FOLLOW,
			TableName: TABLE_MEMBER_FOLLOW,
			Db:        MyCat,
		},
	}
}

func NewMemberFollowTx(tx *dbr.Tx) *MemberFollow {
	return &MemberFollow{
		Model: Model{
			CacheKey:  REDIS_MEMBERS_FOLLOW,
			TableName: TABLE_MEMBER_FOLLOW,
			Db:        MyCat,
			Tx:        tx,
		},
	}
}

func (this *MemberFollow) GetSingleByExps(exps map[string]interface{}) error {
	builder := this.Db.Select("*").From(this.TableName)
	return this.SelectWhere(builder, exps).Limit(1).LoadStruct(this)
}

func (this *MemberFollow) GetTotalByExps(exps map[string]interface{}) (int, error) {
	var count int

	builder := this.Db.Select("COUNT(0)").From(this.TableName+" AS f").Join(TABLE_MEMBERS+" AS m", "f.RefID=m.ID")
	_, err := this.SelectWhere(builder, exps).
		Limit(1).
		LoadStructs(&count)

	return count, err
}

func (this *MemberFollow) GetPaginationByExps(exps map[string]interface{}, page int, limit int) (*PageFollowJson, error) {
	total, err := this.GetTotalByExps(exps)
	if err != nil {
		return nil, err
	}

	offset := limit * (page - 1)
	jsns, err := this.GetPageJsonListByExps(exps, uint64(limit), uint64(offset))
	if err != nil {
		return nil, err
	}

	return &PageFollowJson{
		Pagination: FollowPagination{
			Total:    total,
			PageSize: limit,
		},
		Rows: jsns,
	}, nil
}

func (this *MemberFollow) GetPageJsonListByExps(exps map[string]interface{}, limit uint64, offset uint64) ([]*MemberFollowJson, error) {
	data, err := this.GetPageListByExps(exps, limit, offset)
	if err != nil && err != dbr.ErrNotFound {
		return nil, err
	}

	jsns := make([]*MemberFollowJson, len(data))
	for i, v := range data {
		jsns[i] = this.GetJson(v)
	}

	followIds := make([]int64, len(data))
	for i, v := range data {
		followIds[i] = v.RefID
	}

	mp, err := this.getAdvisorMapByIds(followIds)
	if err != nil {
		return nil, err
	}

	for _, v := range jsns {
		a, ok := mp[IDDecrypt(v.Advisor.GUID)]
		if !ok {
			continue
		}

		v.LevelName = a.LevelName.String
	}

	return jsns, nil
}

func (this *MemberFollow) GetJson(f *MemberFollow) *MemberFollowJson {
	a, err := NewMember().GetSingleAdvisor(f.RefID)
	if err != nil {
		return nil
	}
	ajsn, err := NewMember().GetAdvisorJson(a)
	if err != nil {
		return nil
	}

	return &MemberFollowJson{
		Advisor:    *ajsn,
		CreateTime: f.CreateTime,
	}
}

func (this *MemberFollow) GetPageListByExps(exps map[string]interface{}, limit uint64, offset uint64) ([]*MemberFollow, error) {
	var data []*MemberFollow

	builder := this.Db.Select("f.*").From(this.TableName+" AS f").Join(TABLE_MEMBERS+" AS m", "f.RefID=m.ID")
	err := this.SelectWhere(builder, exps).
		Offset(offset).
		Limit(limit).
		LoadStruct(&data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (this *MemberFollow) GetFollowIdListByMemberId(id int64) ([]int64, error) {
	var ids []int64

	exps := map[string]interface{}{
		"MemberID=?": id,
	}

	builder := this.Db.Select("RefID").From(this.TableName)
	_, err := this.SelectWhere(builder, exps).
		LoadStructs(&ids)

	return ids, err
}

func (this *MemberFollow) FollowByMemberId(memberId int64, followId int64) error {
	cacheKey := fmt.Sprintf(this.CacheKey, memberId)
	_, err := redis.Sadd(cacheKey, []byte(fmt.Sprintf("%v", followId)))

	return err
}

func (this *MemberFollow) UnfollowByMemberId(memberId int64, followId int64) error {
	cacheKey := fmt.Sprintf(this.CacheKey, memberId)
	err := redis.Srem(cacheKey, []byte(fmt.Sprintf("%v", followId)))

	return err
}

func (this *MemberFollow) FanByMemberId(memberId int64) error {
	cacheKey := fmt.Sprintf(REDIS_MEMBERS_FANS, memberId)
	_, err := redis.Do("INCR", cacheKey)

	return err
}

func (this *MemberFollow) UnfanByMemberId(memberId int64) error {
	cacheKey := fmt.Sprintf(REDIS_MEMBERS_FANS, memberId)

	if b, err := redis.Do("EXISTS", cacheKey); err == nil {
		if exist, ok := b.(int64); ok && exist == 0 {
			return nil
		}
	}

	if s, err := redis.Get(cacheKey); err == nil {
		fans, err := strconv.Atoi(s)
		if err == nil && fans == 0 {
			return nil
		}
	}

	_, err := redis.Do("DECR", cacheKey)

	return err
}

func (this *MemberFollow) GetFollowIdsByMemberId(memberId int64) ([]int64, error) {
	cacheKey := fmt.Sprintf(this.CacheKey, memberId)
	ids, err := redis.Smembers(cacheKey)

	if err == nil {
		var err error
		follows := make([]int64, len(ids))
		for i, v := range ids {
			id, errParse := strconv.ParseInt(v, 10, 64)
			if errParse != nil {
				redis.Del(cacheKey)
				err = errParse

				break
			} else {
				follows[i] = int64(id)
			}
		}

		if err == nil {
			return follows, nil
		}
	}

	follows, err := this.GetFollowIdListByMemberId(memberId)
	if err != nil {
		return nil, err
	}

	l := len(follows)
	if l == 0 {
		return []int64{}, nil
	}

	// 重建缓存
	redis.Del(cacheKey)
	for _, v := range follows {
		if _, err := redis.Sadd(cacheKey, []byte(fmt.Sprintf("%v", v))); err != nil {
			redis.Del(cacheKey)

			return nil, err
		}
	}

	return follows, err
}

func (this *MemberFollow) GetFansCountByMemberId(memberId int64) (int, error) {
	cacheKey := fmt.Sprintf(REDIS_MEMBERS_FANS, memberId)
	if s, err := redis.Get(cacheKey); err == nil {
		fans, err := strconv.Atoi(s)
		if err == nil {
			return fans, nil
		}
	}

	exps := map[string]interface{}{
		"RefID=?": memberId,
	}
	fans, err := this.GetCount(exps)
	if err != nil {
		return 0, err
	}

	// 重建缓存
	redis.Del(cacheKey)
	if err := redis.Set(cacheKey, []byte(fmt.Sprintf("%v", fans))); err != nil {
		return 0, err
	}

	return fans, nil
}

func (this *MemberFollow) ResetCache(memberId int64) error {
	follows, err := this.GetFollowIdListByMemberId(memberId)
	if err != nil {
		return err
	}

	l := len(follows)
	if l == 0 {
		return nil
	}

	// 重建关注列表缓存
	cacheKeyFollow := fmt.Sprintf(this.CacheKey, memberId)

	if err := redis.Del(cacheKeyFollow); err != nil {
		return err
	}

	for _, v := range follows {
		if _, err := redis.Sadd(cacheKeyFollow, []byte(fmt.Sprintf("%v", v))); err != nil {
			redis.Del(cacheKeyFollow)

			return err
		}
	}

	// 重建粉丝数量缓存
	cacheKeyFans := fmt.Sprintf(REDIS_MEMBERS_FANS, memberId)
	if err := redis.Del(cacheKeyFans); err != nil {
		return err
	}

	exps := map[string]interface{}{
		"RefID=?": memberId,
	}
	fans, err := this.GetCount(exps)
	if err != nil {
		return err
	}

	if err := redis.Set(cacheKeyFans, []byte(fmt.Sprintf("%v", fans))); err != nil {
		return err
	}

	return nil
}

func (this *MemberFollow) GetFollowsAndFansCount(id int64) (*FollowAndFansCount, error) {
	follows, err := this.GetFollowIdsByMemberId(id)
	if err != nil {
		return nil, err
	}

	fans, err := this.GetFansCountByMemberId(id)
	if err != nil {
		return nil, err
	}

	followIds := make([]string, len(follows))
	for i, v := range follows {
		followIds[i] = IDEncrypt(v)
	}

	return &FollowAndFansCount{
		Fans:      fans,
		FollowIds: followIds,
	}, nil
}

// 获取指定投顾的关注者
func (this *MemberFollow) GetAdvisorFollowers(advisorID int64) ([]int64, error) {
	var ids []int64

	exps := map[string]interface{}{
		"RefID=?": advisorID,
	}

	builder := this.Db.Select("MemberID").From(this.TableName)
	_, err := this.SelectWhere(builder, exps).
		LoadStructs(&ids)

	return ids, err
}

// --------------------------------------------------------------------------------

func (this *MemberFollow) getAdvisorMapByIds(Ids []int64) (map[int64]*FollowAdvisorInfo, error) {
	if len(Ids) == 0 {
		return make(map[int64]*FollowAdvisorInfo), nil
	}

	data := []*FollowAdvisorInfo{}
	exps := map[string]interface{}{
		"a.`MemberID` IN ?": Ids,
	}
	builder := this.Db.Select("a.`MemberID`,g.`GroupName` AS LevelName").
		From(TABLE_MEMBER_ADVISORS+" AS a").
		LeftJoin(TABLE_GROUPS+" AS g", "g.ID=a.Level")

	_, err := this.SelectWhere(builder, exps).LoadStructs(&data)
	if err != nil {
		return nil, err
	}

	r := make(map[int64]*FollowAdvisorInfo, len(data))

	for _, v := range data {
		r[v.MemberID] = v
	}

	return r, nil
}
