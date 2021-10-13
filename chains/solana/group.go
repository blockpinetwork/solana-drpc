package solana

import (
	"math"
	"sort"
)

var (
	RPC_NODE_GROUPS = []string{"INDEX"}
	rpcGroupMap = make(map[string][]string)
)

func init() {
	rpcGroupMap["INDEX"] = []string{
		"CS9ZNjPG9qLrdhuL18A6tZ1NbQASiBKWo8ua1pyv5XT9",
		"882HshS9BnDQPoVSGaUzL2KTT25kRW7pErPns3oMj2Zd",
	}
}

func GetGroupRpcNodes(name string) Nodes {
	var nodes Nodes
	pubKeys := rpcGroupMap[name]
	for _, pubKey := range pubKeys{
		if node, ok := MapGetNode(pubKey); ok{
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func CheckGroupRpcNodes()  {
	for _, pubKeys := range rpcGroupMap{
		for _, pubKey := range pubKeys{
			//logger.Debug("CheckGroupRpcNodes Push job ", pubKey)
			rpcCheckQueue.PushJob(&CheckNodeJob{pubKey})
		}
	}
}

func IsInGroupRpcNodes(pubKey string) bool {
	for _, pubKeys := range rpcGroupMap{
		for _, key := range pubKeys{
			if key == pubKey {
				return true
			}
		}
	}
	return false
}

func SetGroupRpcBackends() {
	for _, group := range RPC_NODE_GROUPS{
		rpcNodes := GetGroupRpcNodes(group)
		if rpcNodes != nil {
			sort.Sort(rpcNodes)
			backends := make(map[string]int64)
			minLatency := int64(math.Max(100, float64(rpcNodes[0].Latency+100)))
			for _, node := range rpcNodes {
				if node.Latency <= minLatency {
					backends[node.GetRPCEndpoint()] = 100
				}
			}
			if len(backends) > 0 {
				logger.Info("Setting proxy server backends: [", group, "] ", len(backends))
				proxyServer.SetBackendsByGroup(group, backends)
			}
		}
	}
}
