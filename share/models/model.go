package models

import (
	"database/sql"
    "fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/beevik/guid"
	"stock/share/gocraft/dbr"
	"stock/share/lib"
	"stock/share/lib/crypto"
	"stock/share/store/redis"
)

type Model struct {
	CacheKey  string       `json:"-"`
	Db        *dbr.Session `json:"-"`
	TableName string       `json:"-"`
	Tx        *dbr.Tx      `json:"-"`
	ViewName  string       `json:"-"`
}

func (this *Model) Delete(exps map[string]interface{}, conditions ...dbr.Condition) error {
	builder := this.Db.DeleteFrom(this.TableName)
	_, err := this.DeleteWhere(builder, exps, conditions...).Exec()

	return err
}

func (this *Model) DelCache(id int64) {
	if len(this.CacheKey) == 0 {
		return
	}
	redis.Del(fmt.Sprintf(this.CacheKey, id))
}

func (this *Model) GetIds(exps map[string]interface{}, conditions ...dbr.Condition) ([]int64, error) {
	var ids []int64
	builder := this.Db.Select("ID").From(this.TableName)
	err := this.SelectWhere(builder, exps, conditions...).LoadValue(&ids)
	return ids, err
}

func (this *Model) Insert(params map[string]interface{}) (int64, error) {
	builder := this.Db.InsertInto(this.TableName)
	result, err := this.InsertParams(builder, params).Exec()
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

func (this *Model) BatchInsert(colums []string, params []interface{}) (int64, error) {
	if len(params) > 500 {
		return 0, fmt.Errorf("Insert data can't more 500, current %v", len(params))
	}

	data := make([]string, len(params))
	for index, v := range params {
		val := reflect.ValueOf(v)
		if val.Kind() != reflect.Slice {
			return 0, fmt.Errorf("Insert data must is slice type")
		}

		_val := ""
		_data := make([]string, val.Len())

		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			itemVal := item.Interface()

			switch s := itemVal.(type) {
			case string:
				_val = fmt.Sprintf("'%s'", s)
			case dbr.NullString:
				_val = fmt.Sprintf("'%s'", s.String)
			default:
				_val = fmt.Sprintf("%d", s)
			}
			_data[i] = _val
		}
		data[index] = fmt.Sprintf("(%s)", strings.Join(_data, ","))
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		this.TableName,
		strings.Join(colums, ","),
		strings.Join(data, ","))

	tx, _ := this.Db.Begin()
	res, err := tx.InsertBySql(sql).Exec()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	id, _ := res.LastInsertId()

	return id, nil
}

func (this *Model) IsExist(exps map[string]interface{}, field string, value string, conditions ...dbr.Condition) (bool, error) {
	var key string

	if len(field) == 0 {
		return true, ErrParameterError
	}

	builder := this.Db.Select(field).From(this.TableName)
	err := this.SelectWhere(builder, exps, conditions...).Limit(1).LoadValue(&key)

	if err != nil {
		if err == dbr.ErrNotFound {
			return false, nil
		}
	} else {
		if key == value {
			return false, nil
		}
	}

	return true, err
}

func (this *Model) Update(params map[string]interface{}, exps map[string]interface{}, conditions ...dbr.Condition) error {
	builder := this.Db.Update(this.TableName)

	this.UpdateParams(builder, params)

	_, err := this.UpdateWhere(builder, exps, conditions...).Exec()
	return err
}

func (this *Model) GetCount(exps map[string]interface{}, conditions ...dbr.Condition) (int, error) {
	var count int

	builder := this.Db.Select("COUNT(0)").From(this.TableName)
	_, err := this.SelectWhere(builder, exps, conditions...).
		Limit(1).
		LoadStructs(&count)

	return count, err
}

func (this *Model) Increment(base string, value interface{}, field string, step int) error {
	cmd := fmt.Sprintf(`UPDATE %s SET %s = %s + %d WHERE %s = %v`,
		this.TableName, field, field, step, base, value)

	_, err := this.Db.UpdateBySql(cmd).Exec()
	return err
}

// --------------------------------------------------------------------------------

func (this *Model) InsertParams(builder *dbr.InsertBuilder, exps map[string]interface{}) *dbr.InsertBuilder {
	var column []string
	var value []interface{}
	for k, v := range exps {
		column = append(column, k)
		value = append(value, v)
	}
	builder.InsertStmt.Columns(column...)
	builder.InsertStmt.Values(value...)
	return builder
}

func (this *Model) UpdateParams(builder *dbr.UpdateBuilder, exps map[string]interface{}) *dbr.UpdateBuilder {
	for k, v := range exps {
		builder.UpdateStmt.Set(k, v)
	}
	return builder
}

func (this *Model) DeleteWhere(builder *dbr.DeleteBuilder, exps map[string]interface{}, conditions ...dbr.Condition) *dbr.DeleteBuilder {
	for k, v := range exps {
		values, ok := v.([]interface{})
		if ok {
			builder.DeleteStmt.Where(k, values...)
		} else {
			builder.DeleteStmt.Where(k, v)
		}
	}

	// 支持复杂表达式
	for _, condition := range conditions {
		builder.DeleteStmt.Where(condition)
	}
	return builder
}

func (this *Model) UpdateWhere(builder *dbr.UpdateBuilder, exps map[string]interface{}, conditions ...dbr.Condition) *dbr.UpdateBuilder {
	for k, v := range exps {
		values, ok := v.([]interface{})
		if ok {
			builder.UpdateStmt.Where(k, values...)
		} else {
			builder.UpdateStmt.Where(k, v)
		}
	}

	// 支持复杂表达式
	for _, condition := range conditions {
		builder.UpdateStmt.Where(condition)
	}
	return builder
}

func (this *Model) SelectWhere(builder *dbr.SelectBuilder, exps map[string]interface{}, conditions ...dbr.Condition) *dbr.SelectBuilder {
	// 支持AND表达式
	for k, v := range exps {
		values, ok := v.([]interface{})
		if ok {
			builder.SelectStmt.Where(k, values...)
		} else {
			builder.SelectStmt.Where(k, v)
		}

	}

	// 支持复杂表达式
	for _, condition := range conditions {
		builder.SelectStmt.Where(condition)

	}
	return builder
}

// --------------------------------------------------------------------------------

func FormatInt(id int64) string {
	return strconv.FormatInt(id, 10)
}

func ConvertInt(n interface{}) int {
	var result int
	switch n.(type) {
	case int:
		result = n.(int)
	case int64:
		v, _ := n.(int64)
		result = int(v)
	case float64:
		v, _ := n.(float64)
		result = int(v)
	}
	return result
}

func ConvertInt64(n interface{}) int64 {
	var result int64
	switch n.(type) {
	case int64:
		result = n.(int64)
	case float64:
		v, _ := n.(float64)
		result = int64(v)
	case int:
		i, _ := n.(int)
		result = int64(i)
	case string:
		result, _ = strconv.ParseInt(n.(string), 10, 64)
	default:
	}
	return result
}

func IDEncrypt(v int64) string {
	return lib.IDEncrypt(v)
}

func IDDecrypt(s string) int64 {
	return lib.IDDecrypt(s)
}

func NewUUID() string {
	return guid.New().String()
}

func NewUniqueID() string {
	return strings.Replace(guid.New().String(), "-", "", -1)
}

func NewAvatarNumber() string {
	return fmt.Sprintf("%v", rand.Intn(24))
}

func New6AuthCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%v", r.Intn(899999)+100000)
}

func New4BitCDKey() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%04x", r.Intn(65535))
}

func New6BitCDKey() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06x", r.Intn(16777215))
}

func New8BitCDKey() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%08x", r.Int63n(4294967295))
}

func NewBatchNo() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("3%v%v", time.Now().Format("060102150405"), r.Intn(89999)+10000)
}

func NewTimestamp() int64 {
	return time.Now().Unix()
}

// --------------------------------------------------------------------------------

func StructToMap(dst interface{}) map[string]interface{} {
	vv := reflect.ValueOf(dst)
	t := reflect.Indirect(vv).Type()
	v := vv.Elem()

	data := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		if tag == "-" {

			// ignore
			continue
		}
		if v.Field(i).Type().Kind() == reflect.Struct {
			data[field.Name] = v.Field(i).Field(0).Field(0).Interface()
			continue
		}
		data[field.Name] = v.Field(i).Interface()
	}
	return data
}

func MapToStruct(dst interface{}, src map[string]string) error {
	var err error
	for k, v := range src {
		err = setValue(dst, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetAvtar(memberId int64, lastTime int64, avatar string) string {
	if len(avatar) > 0 {
		return avatar
	}

	if memberId == 0 {
		return fmt.Sprintf("%su0.png", BASE_URL_AVATAR)
	}
	str := strconv.Itoa(int(memberId))
	return BASE_URL_AVATAR + fmt.Sprintf("%s_180.png?%v", crypto.GetMD5(str, false), lastTime)
}

// --------------------------------------------------------------------------------

func setValue(obj interface{}, name string, value string) error {
	elem := reflect.ValueOf(obj).Elem()
	fieldName := elem.FieldByName(name)

	if !fieldName.IsValid() {
		return fmt.Errorf("model: No such field %v", name)
	}

	if !fieldName.CanSet() {
		return fmt.Errorf("model: Cannot set %v field value", name)
	}

	fieldType := fieldName.Type()
	v := reflect.ValueOf(value)
	if fieldType != v.Type() {
		var typeValue interface{}
		switch fieldType.Kind() {
		case reflect.Int:
			typeValue, _ = strconv.Atoi(value)
		case reflect.Int64:
			typeValue, _ = strconv.ParseInt(value, 10, 64)
		case reflect.Float64:
			typeValue, _ = strconv.ParseFloat(value, 64)
		case reflect.Struct:
			switch fieldName.Interface().(type) {
			case dbr.NullString:
				typeValue = dbr.NullString{sql.NullString{value, true}}
			case dbr.NullInt64:
				nullInt64, _ := strconv.ParseInt(value, 10, 64)
				typeValue = dbr.NullInt64{sql.NullInt64{nullInt64, true}}
			default:
				fieldName.Field(0).Set(reflect.ValueOf(sql.NullString{value, true}))
				return nil
			}
		default:
			return fmt.Errorf("model: Type didn't match %v", name)
		}
		v = reflect.ValueOf(typeValue)
	}
	fieldName.Set(v)
	return nil
}
