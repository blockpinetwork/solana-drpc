package db

import "github.com/blockpilabs/solana-drpc/db/redis"

var (
	redisPools = make(map[string]*redis.ConnPool)
)

func InitRedisPool(poolName, prefix, host, password string, database, maxOpenConns, maxIdleConns int) *redis.ConnPool{
	pool := redis.InitRedisPool(prefix, host, password, database, maxOpenConns, maxIdleConns)
	redisPools[poolName] = pool

	return pool
}

func GetRedisPool(poolName string) *redis.ConnPool {
	return redisPools[poolName]
}
