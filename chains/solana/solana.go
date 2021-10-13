package solana

import (
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/blockpilabs/solana-drpc/db"
	"github.com/blockpilabs/solana-drpc/log"
	"github.com/blockpilabs/solana-drpc/network/http"
	"github.com/blockpilabs/solana-drpc/plugins/statistic"
	"github.com/blockpilabs/solana-drpc/proxy"
	"github.com/blockpilabs/solana-drpc/queue"
	"github.com/blockpilabs/solana-drpc/rpc"
	"github.com/blockpilabs/solana-drpc/utils"
)

var logger = log.GetLogger("solana")

var (
	NODE_POOL_NAME = "SOLANA_NODES"
	CHECK_INTERVAL = int64(60 * 1000)
	MAX_NODE_RETENTION_TIME = int64(24 * 60 * 60 * 1000)

	rpcCheckQueue *queue.Queue

	backendRPCNodes *Nodes
	sortedRPCNodes *Nodes

	proxyServer *proxy.ProxyServer
	locker sync.Mutex
)


func InitNodeRedisPools(prefix, host, password string, database int) {
	db.InitRedisPool(NODE_POOL_NAME, prefix, host, password, database, -1, 5)
}

func saveClusterNodes(endpoint string)  {
	jsonReq := "{\"jsonrpc\":\"2.0\",\"id\":\"blockpi-drpc\", \"method\":\"getClusterNodes\",\"params\": []}"
	data := http.PostJson(endpoint, jsonReq)
	if data != nil {
		resp := rpc.JSONRpcResponse{}
		if err := json.Unmarshal(data, &resp); err == nil && resp.Error == nil {
			nodes := make([]Node, len( resp.Result.([]interface{})))
			if v, err := json.Marshal(resp.Result); err == nil{
				if err := json.Unmarshal(v, &nodes); err == nil {
					for i, _ := range nodes{
						//logger.Debug("saveClusterNodes ", nodes[i].PubKey)
						SaveNode(&nodes[i], true)
					}
				}
			}
		}
	}
}

func initNodesFromRedis() {
	conn := db.GetRedisPool(NODE_POOL_NAME)
	if keys, err := conn.Keys("*"); err == nil {
		for _, key := range keys {
			if v, err := conn.GetString(key); err == nil{
				var node Node
				err := json.Unmarshal([]byte(v), &node)
				if err == nil {
					MapSetNode(&node)
				}
			}
		}
	}
}

func UpdateClusterNodes()  {
	initNodesFromRedis()

	go func() {
		for {
			var seedEndpoints []string
			seedEndpoints = append(seedEndpoints, "http://api.mainnet-beta.solana.com")

			if sortedRPCNodes != nil {
				count := 0
				for _, node := range *sortedRPCNodes {
					count++
					seedEndpoints = append(seedEndpoints, node.GetRPCEndpoint())
					if count >= 5 {
						break
					}
				}
			}

			for _, endpoint := range seedEndpoints {
				saveClusterNodes(endpoint)
			}

			time.Sleep(time.Second * 30)
		}
	}()
}

func StartCheckQueue() {
	rpcCheckQueue = queue.NewQueue(1000)
	rpcCheckQueue.Run()

	go func() {
		for {
			time.Sleep(time.Second * 10)
			CheckGroupRpcNodes()
			nodeMap.Range(func(key, value interface{}) bool {
				node := value.(*Node)
				//logger.Debug("StartCheckQueue Push job ", node.PubKey)
				if !node.Checking && len(node.RPC) > 0 {
					if utils.CurrentTimestampMilli() - node.LastCheckTime >= CHECK_INTERVAL {
						//Check this node
						if !IsInGroupRpcNodes(node.PubKey) {

							rpcCheckQueue.PushJob(&CheckNodeJob{node.PubKey})
						}
					}
				}
				return true
			})

			//for _, node := range nodeMap{
			//	if !node.Checking && len(node.RPC) > 0 {
			//		if utils.CurrentTimestampMilli() - node.LastCheckTime >= CHECK_INTERVAL {
			//			//Check this node
			//			if !IsInGroupRpcNodes(node.PubKey) {
			//				rpcCheckQueue.PushJob(&CheckNodeJob{node.PubKey})
			//			}
			//		}
			//	}
			//}

			/*
			conn := db.GetRedisPool(NODE_POOL_NAME)
			if keys, err := conn.Keys("*"); err == nil {
				for _, key := range keys {
					if v, err := conn.GetString(key); err == nil{
						var node Node
						err := json.Unmarshal([]byte(v), &node)
						if err == nil {
							if !node.Checking && len(node.RPC) > 0 {
								if utils.CurrentTimestampMilli() - node.LastCheckTime >= CHECK_INTERVAL {
									//Check this node
									if !IsInGroupRpcNodes(node) {
										rpcCheckQueue.PushJob(&CheckNodeJob{node: &node})
									}
								}
							}
						}
					}
				}
			}
			*/


		}
	}()
}

type CheckNodeJob struct{
	NodePubKey string
}

func (job CheckNodeJob) Do() {
	if node, ok := MapGetNode(job.NodePubKey); ok {
		node.Check()
	}
	//logger.Debug("CheckNodeJob Do")
}

func SetProxyServer(server *proxy.ProxyServer) {
	proxyServer = server
}

func UpdateBestRPCNodes() {
	go func() {
		for {
			time.Sleep(time.Second * 10)
			getBestRPCNodes()
			setBackendRpcNodes()
			time.Sleep(time.Second * 20)
		}
	}()
}

func getBestRPCNodes() {
	locker.Lock()
	defer locker.Unlock()

	conn := db.GetRedisPool(NODE_POOL_NAME)
	var rpcNodes Nodes

	keys, err := conn.Keys("*")
	if err == nil {
		for _, key := range keys {
			if v, err := conn.GetString(key); err == nil{
				var node Node
				if err := json.Unmarshal([]byte(v), &node); err == nil {
					if node.Healthy && len(node.RPC) > 0 {
						rpcNodes = append(rpcNodes, &node)
					}
				}
			}
		}
	}

	if len(rpcNodes) == 0 {
		logger.Error("No RPC nodes found.")
		return
	}

	sort.Sort(rpcNodes)
	sortedRPCNodes = &rpcNodes

	var backends Nodes

	minLatency := int64(math.Max(100, float64(rpcNodes[0].Latency + 50)))
	for _, node := range rpcNodes {
		if node.Latency <= minLatency {
			backends = append(backends, node)
		}
	}
	if len(backends) > 0 {
		backendRPCNodes = &backends
	}
}

func setBackendRpcNodes() {
	locker.Lock()
	defer locker.Unlock()

	if backendRPCNodes != nil {
		backends := make(map[string]int64)
		for _, node := range *backendRPCNodes {
			backends[node.GetRPCEndpoint()] = 100
		}
		if len(backends) > 0{
			logger.Info("Setting proxy server backends: ", len(backends))
			proxyServer.SetBackends(backends)
		}
	}

	SetGroupRpcBackends()
}

func Summary() map[string]interface{} {
	locker.Lock()
	defer locker.Unlock()

	result := make(map[string]interface{})

	conn := db.GetRedisPool(NODE_POOL_NAME)
	connStatistic := statistic.GetRedisConn()

	keys, _ := conn.Keys("*");
	result["total_peers"] = len(keys)

	result["total_rpc_nodes"] = 0
	if sortedRPCNodes != nil {
		result["total_rpc_nodes"] = len(*sortedRPCNodes)
	}

	result["total_rpc_backends"] = 0
	if backendRPCNodes != nil {
		result["total_rpc_backends"] = len(*backendRPCNodes)
	}

	result["backend_rpc_nodes"] = backendRPCNodes
	result["all_rpc_nodes"] = sortedRPCNodes

	result["total_requests"] = 0
	if count, err := connStatistic.GetInt64("REQ_TOTAL"); err == nil {
		result["total_requests"] = count
	}

	reqTimeMap := make(map[int64]int64)
	result["requests"] = reqTimeMap

	ts := utils.CurrentHourTimestamp()
	hours := 1
	totalReq24h := int64(0)
	for {
		reqTimeMap[ts] = 0
		if count, err := connStatistic.GetInt64("REQ_T_"+ strconv.FormatInt(ts, 10)); err == nil {
			reqTimeMap[ts] = count
			totalReq24h += count
		}
		ts -= 3600
		hours++
		if hours > 24 {
			break
		}
	}

	result["total_requests_24h"] = totalReq24h
	result["hourly_avg_requests_24h"] = totalReq24h / 24

	return result
}