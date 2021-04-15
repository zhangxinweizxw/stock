package redis

import (
    "errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"stock/share/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
)

func Init(addr, db, password string, timeout int) {
	if pool == nil {
		pool = newRedisPool(addr, db, password, timeout)
	}
}

func Del(key string) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("DEL", key)
	return err
}

func Get(key string) (string, error) {
	c := pool.Get()
	defer c.Close()

	return redis.String(c.Do("GET", key))
}

func Hgetall(key string) (map[string]string, error) {
	c := pool.Get()
	defer c.Close()

	res, err := bytesSlice(c.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	writeToContainer(res, reflect.ValueOf(result))

	return result, err
}

func Hmset(key string, mapping interface{}) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("HMSET", redis.Args{}.Add(key).AddFlat(mapping)...)
	return err
}

func Keys(pattern string) ([]string, error) {
	c := pool.Get()
	defer c.Close()

	return redis.Strings(c.Do("KEYS", pattern))
}

func Set(key string, val []byte) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("SET", key, string(val))
	return err
}

func Setex(key string, expire int, val []byte) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("SETEX", key, expire, string(val))
	return err
}

func Smembers(key string) ([]string, error) {
	c := pool.Get()
	defer c.Close()

	return redis.Strings(c.Do("SMEMBERS", key))
}

func Sadd(key string, val []byte) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("SADD", key, string(val))
	return err
}

func Srem(key string, val []byte) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("SREM", key, string(val))
	return err
}

// List commands

func LRange(key string, start int, end int) ([]string, error) {
	c := pool.Get()
	defer c.Close()

	return redis.Strings(c.Do("LRANGE", key, start, end))
}
func LSet(key string, where int, val []byte) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("LSET", key, where, string(val))
	return err
}

func Lrem(key string, count int, val []byte) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("LREM", key, count, string(val))
	return err
}

func Rpush(key string, val []byte) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("RPUSH", key, string(val))
	return err
}

func Lpush(key string, val []byte) error {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("LPUSH", key, string(val))
	return err
}

func Llen(key string) (int, error) {
	c := pool.Get()
	defer c.Close()

	res, err := c.Do("LLEN", key)
	if err != nil {
		return -1, err
	}
	return int(res.(int64)), nil
}

func Lpop(key string) ([]byte, error) {
	c := pool.Get()
	defer c.Close()

	res, err := c.Do("LPOP", key)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("EOF")
	}
	return res.([]byte), nil
}

func Rpop(key string) ([]byte, error) {
	c := pool.Get()
	defer c.Close()

	res, err := c.Do("RPOP", key)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("EOF")
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
	c := pool.Get()
	defer c.Close()

	return c.Do(cmd, args...)
}

// ------------------------------------------------------------------------

func newRedisPool(addr string, db string, password string, timeout int) *redis.Pool {

	// Set dial options
	// Specifies the timeout for connecting to the Redis server
	dialOptions := make([]redis.DialOption, 1)
	dialOptions[0] = redis.DialConnectTimeout(time.Duration(timeout) * time.Second)

	return &redis.Pool{
		MaxIdle:     80,
		MaxActive:   10000,
		IdleTimeout: 600 * time.Second,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", addr, dialOptions...)
			if err != nil {
				return nil, err
			}
			_, err = con.Do("AUTH", password)
			if err == nil {
				con.Do("SELECT", db)
			}
			return con, err
		},
	}
}

func bpop(cmd string, key string, timeout int) (interface{}, error) {
	c := pool.Get()
	defer c.Close()

	res, err := c.Do(cmd, key, timeout)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, errors.New("EOF")
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
				return nil, fmt.Errorf("redigo: Unexpected element type for []byte, got type %T", reply[i])
			}
			result[i] = p
		}
		return result, nil
	case nil:
		return nil, redis.ErrNil
	case redis.Error:
		return nil, reply
	}
	return nil, fmt.Errorf("redigo: Unexpected type for []byte, got type %T", reply)
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
			return errors.New("redigo: Invalid map type")
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
		return errors.New("redigo: Invalid container type")
	}
	return nil
}
