package common

import (
    "fmt"
    "strconv"

/share/models"

	"stock
/share/store/redis"
)

type Group struct {
	Model       `db:"-" `
	ID          int64
	Assort      int
	GroupName   string
	Description string
	IsProtected int
	Color       string
	Weight      int
}

type GroupJson struct {
	GUID      string `json:"_id"`
	Assort    int    `json:"assort"`
	Color     string `json:"color"`
	GroupName string `json:"name"`
	Weight    int    `json:"weight"`
}

func NewGroup() *Group {
	return &Group{
		Model: Model{
			CacheKey:  REDIS_GROUPS,
			TableName: TABLE_GROUPS,
			Db:        MyCat,
		},
	}
}

func (this *Group) GetList() ([]*Group, error) {
	var groups []*Group
	_, err := this.Db.Select("*").From(this.TableName).LoadStructs(&groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (this *Group) GetDefault() error {
	exps := map[string]interface{}{
		"IsProtected=?": 1,
	}
	builuder := this.Db.Select("*").From(this.TableName)
	return this.SelectWhere(builuder, exps).LoadValue(this)

}

// 根据组权重随机获取一个投顾ID
func (this *Group) GetAdvisorIDByWeight(weight int) (int64, error) {
	key := fmt.Sprintf(REDIS_MAJOR_ADVISOR_GROUP, weight)
	res, err := redis.Srandmember(key, 1)
	if err == nil && len(res) > 0 {
		id, err := strconv.Atoi(res[0])
		return int64(id), err
	}
	var aids []int64
	exps := map[string]interface{}{
		"g.Weight=?": weight,
	}
	builder := this.Db.Select("a.MemberID").From(this.TableName+" AS g").Join(TABLE_MEMBER_ADVISORS+" AS a", "g.ID=a.Level")
	lens, err := this.SelectWhere(builder, exps).LoadValues(&aids)
	if lens == 0 {
		return 0, fmt.Errorf("Not match advisor by weight")
	}
	if err != nil {
		return 0, err
	}
	for _, id := range aids {
		redis.Do("SADD", key, fmt.Sprintf("%d", id))
	}

	return aids[0], nil
}

func (this *Group) GetSingle(id int64, assort int) error {
	cacheKey := fmt.Sprintf(this.CacheKey, id)
	rec, err := redis.Hgetall(cacheKey)

	if err == nil && len(rec) > 0 {
		if err := MapToStruct(this, rec); err != nil {
			redis.Del(cacheKey)
		} else if this.Assort == assort {
			return nil
		}
	}

	exps := map[string]interface{}{
		"ID=?":     id,
		"Assort=?": assort,
	}
	builder := this.Db.Select("*").From(this.TableName)

	if err := this.SelectWhere(builder, exps).LoadStruct(&this); err != nil {
		return err
	}

	if err := redis.Hmset(cacheKey, StructToMap(this)); err != nil {
		return err
	}

	return nil
}

func (this *Group) GetWeightById(id int64, assort int) (int, error) {
	group := NewGroup()
	if err := group.GetSingle(id, assort); err != nil {
		return 0, err
	}

	return group.Weight, nil
}

func (this *Group) GetJson(g *Group) *GroupJson {
	return &GroupJson{
		GUID:      IDEncrypt(g.ID),
		Assort:    g.Assort,
		GroupName: g.GroupName,
		Color:     g.Color,
		Weight:    g.Weight,
	}
}

func (this *Group) GetJsonList() ([]*GroupJson, error) {
	jsns := []*GroupJson{}
	data, err := this.GetList()

	if err != nil {
		return nil, err
	}

	for _, g := range data {
		jsns = append(jsns, this.GetJson(g))
	}

	return jsns, nil
}
