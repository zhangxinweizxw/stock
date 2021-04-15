package cluster

import (
    "errors"
    "fmt"
    "reflect"
    "strconv"
    "time"

    r "stock
/share/garyburd/redigo/redis"
	"stock
/share/logging"
)

var (
	c redis.Cluster
)

//初始化集群
func Init(addr, db, password string, timeout int) {
	var err error
	if c == nil {
		c, err = redis.NewCluster(
			&redis.Options{
				StartNodes:   []string{addr},
				ConnTimeout:  time.Duration(timeout) * time.Second,
				ReadTimeout:  time.Duration(timeout) * time.Second,
				WriteTimeout: time.Duration(timeout) * time.Second,
				KeepAlive:    16,
				AliveTime:    60 * time.Second,
			})

		if err != nil {
			logging.Fatal(err)
		}
	}
}

func Del(key string) error {
	_, err := c.Do("DEL", key)
	return err
}

func Get(key string) (string, error) {

	return redis.String(c.Do("GET", key))
}

func Hgetall(key string) (map[string]string, error) {

	res, err := bytesSlice(c.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	writeToContainer(res, reflect.ValueOf(result))

	return result, err
}

func Hmset(key string, mapping interface{}) error {
	_, err := c.Do("HMSET", r.Args{}.Add(key).AddFlat(mapping)...)
	return err
}

func Keys(pattern string) ([]string, error) {
	return r.Strings(c.Do("KEYS", pattern))
}

func Set(key string, val []byte) error {
	_, err := c.Do("SET", key, string(val))
	return err
}

func Smembers(key string) ([]string, error) {

	return redis.Strings(c.Do("SMEMBERS", key))
}

func Sadd(key string, val []byte) error {

	_, err := c.Do("SADD", key, string(val))
	return err
}

// List commands

func Rpush(key string, val []byte) error {

	_, err := c.Do("RPUSH", key, string(val))
	return err
}

func Lpush(key string, val []byte) error {

	_, err := c.Do("LPUSH", key, string(val))
	return err
}

func Llen(key string) (int, error) {

	res, err := c.Do("LLEN", key)
	if err != nil {
		return -1, err
	}
	return int(res.(int64)), nil
}

func Lpop(key string) ([]byte, error) {

	res, err := c.Do("LPOP", key)
	if err != nil {
		return nil, err
	}
	return res.([]byte), nil
}

func Rpop(key string) ([]byte, error) {

	res, err := c.Do("RPOP", key)
	if err != nil {
		return nil, err
	}
	return res.([]byte), nil
}

func Blpop(key string, timeout int) (interface{}, error) {
	return bpop("BLPOP", key, timeout)
}

func Brpop(key string, timeout int) (interface{}, error) {
	return bpop("BRPOP", key, timeout)
}

// General Commands

func Do(cmd string, args ...interface{}) (reply interface{}, err error) {
	return c.Do(cmd, args...)
}

func bpop(cmd string, key string, timeout int) (interface{}, error) {
	res, err := c.Do(cmd, key, timeout)
	if err != nil {
		return nil, err
	}

	// Get value from list
	if list, ok := res.([]interface{}); ok {
		for i, value := range list {
			if i == 1 {
				return value, nil
			}
		}
	}
	return nil, errors.New("EOF")
}

func bytesSlice(reply interface{}, err error) ([][]byte, error) {
	if err != nil {
		return nil, err
	}
	switch reply := reply.(type) {
	case []interface{}:
		result := make([][]byte, len(reply))
		for i := range reply {
			if reply[i] == nil {
				continue
			}
			p, ok := reply[i].([]byte)
			if !ok {
				return nil, fmt.Errorf("redigo: unexpected element type for []byte, got type %T", reply[i])
			}
			result[i] = p
		}
		return result, nil
	case nil:
		return nil, r.ErrNil
	case r.Error:
		return nil, reply
	}
	return nil, fmt.Errorf("redigo: unexpected type for []byte, got type %T", reply)
}

func writeTo(data []byte, val reflect.Value) error {
	s := string(data)
	switch v := val; v.Kind() {

	// if we're writing to an interace value, just set the byte data
	// TODO: should we support writing to a pointer?
	case reflect.Interface:
		v.Set(reflect.ValueOf(data))

	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		ui, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(ui)

	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.SetFloat(f)

	case reflect.String:
		v.SetString(s)

	case reflect.Slice:
		typ := v.Type()
		if typ.Elem().Kind() == reflect.Uint || typ.Elem().Kind() == reflect.Uint8 || typ.Elem().Kind() == reflect.Uint16 || typ.Elem().Kind() == reflect.Uint32 || typ.Elem().Kind() == reflect.Uint64 || typ.Elem().Kind() == reflect.Uintptr {
			v.Set(reflect.ValueOf(data))
		}
	}
	return nil
}

func writeToContainer(data [][]byte, val reflect.Value) error {
	switch v := val; v.Kind() {
	case reflect.Ptr:
		return writeToContainer(data, reflect.Indirect(v))
	case reflect.Interface:
		return writeToContainer(data, v.Elem())
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return errors.New("Invalid map type")
		}
		elemtype := v.Type().Elem()
		for i := 0; i < len(data)/2; i++ {
			mk := reflect.ValueOf(string(data[i*2]))
			mv := reflect.New(elemtype).Elem()
			writeTo(data[i*2+1], mv)
			v.SetMapIndex(mk, mv)
		}
	case reflect.Struct:
		for i := 0; i < len(data)/2; i++ {
			name := string(data[i*2])
			field := v.FieldByName(name)
			if !field.IsValid() {
				continue
			}
			writeTo(data[i*2+1], field)
		}
	default:
		return errors.New("redigo: invalid container type")
	}
	return nil
}
