package redis

import (
    "errors"
	"reflect"

	"stock/share/garyburd/redigo/redis"
)

type RedisPool struct {
	pool *redis.Pool
}

func NewRedisPool(addr, db, password string, timeout int) *RedisPool {
	return &RedisPool{
		pool: newRedisPool(addr, db, password, timeout),
	}
}

func (r *RedisPool) Close() {
	r.pool.Close()
}

func (r *RedisPool) Del(key string) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("DEL", key)
	return err
}

func (r *RedisPool) Get(key string) (string, error) {
	return r.GetString(key)
}
func (r *RedisPool) GetBytes(key string) ([]byte, error) {
	c := r.pool.Get()
	defer c.Close()

	return redis.Bytes(c.Do("GET", key))
}
func (r *RedisPool) GetString(key string) (string, error) {
	c := r.pool.Get()
	defer c.Close()

	return redis.String(c.Do("GET", key))
}

func (r *RedisPool) Hgetall(key string) (map[string]string, error) {
	c := r.pool.Get()
	defer c.Close()

	res, err := bytesSlice(c.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	writeToContainer(res, reflect.ValueOf(result))

	return result, err
}

func (r *RedisPool) Hmset(key string, mapping interface{}) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("HMSET", redis.Args{}.Add(key).AddFlat(mapping)...)
	return err
}

func (r *RedisPool) Keys(pattern string) ([]string, error) {
	c := r.pool.Get()
	defer c.Close()

	return redis.Strings(c.Do("KEYS", pattern))
}

func (r *RedisPool) Set(key string, val []byte) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("SET", key, string(val))
	return err
}

func (r *RedisPool) Setex(key string, expire int, val []byte) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("SETEX", key, expire, string(val))
	return err
}

func (r *RedisPool) Smembers(key string) ([]string, error) {
	c := r.pool.Get()
	defer c.Close()

	return redis.Strings(c.Do("SMEMBERS", key))
}

func (r *RedisPool) Sadd(key string, val []byte) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("SADD", key, string(val))
	return err
}

func (r *RedisPool) Srem(key string, val []byte) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("SREM", key, string(val))
	return err
}

// List commands

func (r *RedisPool) LRange(key string, start int, end int) ([]string, error) {
	c := r.pool.Get()
	defer c.Close()

	return redis.Strings(c.Do("LRANGE", key, start, end))
}

func (r *RedisPool) Lrem(key string, count int, val []byte) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("LREM", key, count, string(val))
	return err
}

func (r *RedisPool) Rpush(key string, val []byte) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("RPUSH", key, string(val))
	return err
}

func (r *RedisPool) Lpush(key string, val []byte) error {
	c := r.pool.Get()
	defer c.Close()

	_, err := c.Do("LPUSH", key, string(val))
	return err
}

func (r *RedisPool) Llen(key string) (int, error) {
	c := r.pool.Get()
	defer c.Close()

	res, err := c.Do("LLEN", key)
	if err != nil {
		return -1, err
	}
	return int(res.(int64)), nil
}

func (r *RedisPool) Lpop(key string) ([]byte, error) {
	c := r.pool.Get()
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

func (r *RedisPool) Rpop(key string) ([]byte, error) {
	c := r.pool.Get()
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

func (r *RedisPool) Blpop(key string, timeout int) (interface{}, error) {
	return bpop("BLPOP", key, timeout)
}

func (r *RedisPool) Brpop(key string, timeout int) (interface{}, error) {
	return bpop("BRPOP", key, timeout)
}

// General Commands

func (r *RedisPool) Do(cmd string, args ...interface{}) (reply interface{}, err error) {
	c := r.pool.Get()
	defer c.Close()

	return c.Do(cmd, args...)
}
