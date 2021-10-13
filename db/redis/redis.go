package redis

import (
	"log"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

type ConnPool struct {
	RedisPool *redis.Pool
	keyPrefix string
}

func InitRedisPool(prefix, host, password string, database, maxOpenConns, maxIdleConns int) *ConnPool {
	r := &ConnPool{keyPrefix: prefix}
	r.RedisPool = newPool(host, password, database, maxOpenConns, maxIdleConns)
	if _, err := r.Do("PING"); err != nil {
		log.Panicln("Init redis pool failed.", err.Error())
	}
	return r
}

func newPool(server, password string, database, maxOpenConns, maxIdleConns int) *redis.Pool {
	return &redis.Pool{
		MaxActive:   maxOpenConns, // max number of connections
		MaxIdle:     maxIdleConns,
		Wait: 		 true,
		IdleTimeout: 120 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("select", database); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (p *ConnPool) Close() error {
	err := p.RedisPool.Close()
	return err
}

func (p *ConnPool) RealKey(key string) string{
	if len(p.keyPrefix) > 0{
		if strings.HasPrefix(key, p.keyPrefix) {
			return key
		}
		return p.keyPrefix + key
	}
	return key
}

func (p *ConnPool) Do(command string, args ...interface{}) (interface{}, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	return conn.Do(command, args...)
}

func (p *ConnPool) SetString(key string, value interface{}) (interface{}, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	return conn.Do("SET", p.RealKey(key), value)
}

func (p *ConnPool) GetString(key string) (string, error) {
	// get one connection from pool
	conn := p.RedisPool.Get()
	// put connection to pool
	defer conn.Close()
	return redis.String(conn.Do("GET", p.RealKey(key)))
}

func (p *ConnPool) GetBytes(key string) ([]byte, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	return redis.Bytes(conn.Do("GET", p.RealKey(key)))
}

func (p *ConnPool) GetInt(key string) (int, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("GET", p.RealKey(key)))
}

func (p *ConnPool) GetInt64(key string) (int64, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("GET", p.RealKey(key)))
}

func (p *ConnPool) DelKey(key string) (interface{}, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	return conn.Do("DEL", p.RealKey(key))
}

func (p *ConnPool) ExpireKey(key string, seconds int64) (interface{}, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	return conn.Do("EXPIRE", p.RealKey(key), seconds)
}

func (p *ConnPool) Keys(pattern string) ([]string, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do("KEYS", p.RealKey(pattern)))
}

func (p *ConnPool) KeysByteSlices(pattern string) ([][]byte, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	return redis.ByteSlices(conn.Do("KEYS", p.RealKey(pattern)))
}

func (p *ConnPool) SetHashMap(key string, fieldValue map[string]interface{}) (interface{}, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	return conn.Do("HMSET", redis.Args{}.Add(p.RealKey(key)).AddFlat(fieldValue)...)
}

func (p *ConnPool) GetHashMapString(key string) (map[string]string, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	return redis.StringMap(conn.Do("HGETALL", p.RealKey(key)))
}

func (p *ConnPool) GetHashMapInt(key string) (map[string]int, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	return redis.IntMap(conn.Do("HGETALL", p.RealKey(key)))
}

func (p *ConnPool) GetHashMapInt64(key string) (map[string]int64, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	return redis.Int64Map(conn.Do("HGETALL", p.RealKey(key)))
}

func (p *ConnPool) Incr(key string) (int64, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("INCR", p.RealKey(key)))
}

func (p *ConnPool) IncrBy(key string, increment int64) (int64, error) {
	conn := p.RedisPool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("INCRBY", p.RealKey(key), increment))
}