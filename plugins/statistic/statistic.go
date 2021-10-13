package statistic

import (
	"fmt"

	"github.com/blockpilabs/solana-drpc/common"
	"github.com/blockpilabs/solana-drpc/db"
	"github.com/blockpilabs/solana-drpc/db/redis"
	"github.com/blockpilabs/solana-drpc/log"
	"github.com/blockpilabs/solana-drpc/rpc"
	"github.com/blockpilabs/solana-drpc/utils"
)

var logger = log.GetLogger("statistic")

var (
	REDIS_CONN_POOL_NAME = "STATISTIC"
)

type StatisticMiddleware struct {
	//rpcRequestsReceived  chan *rpc.JSONRpcRequestSession
	//rpcResponsesReceived chan *rpc.JSONRpcRequestSession
	redisConnPool *redis.ConnPool
}

func NewStatisticMiddleware(options ...common.Option) *StatisticMiddleware {
	return &StatisticMiddleware{
		redisConnPool: db.GetRedisPool(REDIS_CONN_POOL_NAME),
	}
}

func GetRedisConn() *redis.ConnPool {
	return db.GetRedisPool(REDIS_CONN_POOL_NAME)
}

func InitStatisticRedisPool(prefix, host, password string, database int)  {
	db.InitRedisPool(REDIS_CONN_POOL_NAME, prefix, host, password, database, -1, 5)
}

func (middleware *StatisticMiddleware) OnRpcRequest(session *rpc.JSONRpcRequestSession) (err error) {
	//middleware.rpcRequestsReceived <- session
	methodName := ""
	if session.MethodNameForCache != nil {
		methodName = *session.MethodNameForCache
	}else{
		methodName = session.Request.Method
	}

	time := utils.CurrentHourTimestamp()
	keys := []string{
		"REQ_TOTAL",
		fmt.Sprintf("REQ_T_%d", time),
		fmt.Sprintf("REQ_M_%s", methodName),
		fmt.Sprintf("REQ_MT_%s_%d", methodName, time),
	}
	for _, key := range keys {
		middleware.redisConnPool.Incr(key)
	}
	logger.Debug(methodName)
	return
}
func (middleware *StatisticMiddleware) OnRpcResponse(session *rpc.JSONRpcRequestSession) (err error) {
	//middleware.rpcResponsesReceived <- session


	return
}