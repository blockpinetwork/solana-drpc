package solana

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/blockpilabs/solana-drpc/db"
	"github.com/blockpilabs/solana-drpc/network/http"
	"github.com/blockpilabs/solana-drpc/rpc"
	"github.com/blockpilabs/solana-drpc/utils"
)

type Node struct {
	FeatureSet		int64	`json:"featureSet"`
	Gossip 			string	`json:"gossip"`
	PubKey 			string	`json:"pubkey"`
	RPC 			string	`json:"rpc"`
	ShredVersion 	int64	`json:"shredVersion"`
	TPU 			string	`json:"tpu"`
	Version 		string	`json:"version"`

	Latency 		int64	`json:"latency"`
	Checking		bool	`json:"checking"`
	LastCheckTime	int64	`json:"lastCheckTime"`
	LastOnlineTime	int64	`json:"lastOnlineTime"`
	Height			int64	`json:"height"`

	Healthy			bool	`json:"healthy"`
	Syncing			bool	`json:"syncing"`
}
var nodeMap sync.Map

type Nodes []*Node

func (nodes Nodes)Len() int {
	return len(nodes)
}

func (nodes Nodes)Less(i, j int) bool {
	return nodes[i].Latency <  nodes[j].Latency
}

func (nodes Nodes)Swap(i, j int) {
	tmp := nodes[i]
	nodes[i] = nodes[j]
	nodes[j] = tmp
}

func (p *Node) GetRPCEndpoint() string {
	if len(p.RPC) > 0 {
		if !(strings.HasPrefix(p.RPC, "http://") || strings.HasPrefix(p.RPC, "https://")){
			return "http://" + p.RPC
		}
	}
	return p.RPC
}

func MapGetNode(pubKey string) (*Node, bool) {
	node, ok := nodeMap.Load(pubKey)
	if node != nil {
		return node.(*Node), ok
	}
	return nil, false
}

func MapSetNode(node *Node) {
	nodeMap.Store(node.PubKey, node)
}


func SaveNode(node *Node, newNode bool) {
	defer func() {
		if err:=recover(); err != nil {
			logger.Error(err)
		}
	}()

	conn := db.GetRedisPool(NODE_POOL_NAME)
	if newNode {
		cachedNode, ok := MapGetNode(node.PubKey)
		if !ok {
			if v, err := conn.GetString(node.PubKey); err == nil {
				if err := json.Unmarshal([]byte(v), cachedNode); err == nil {
					cachedNode = nil
				}
			}
		}

		if cachedNode != nil {
			node.Latency = cachedNode.Latency
			node.Checking = cachedNode.Checking
			node.LastCheckTime = cachedNode.LastCheckTime
			node.LastOnlineTime = cachedNode.LastOnlineTime
			node.Height = cachedNode.Height
			node.Healthy = cachedNode.Healthy
			node.Syncing = cachedNode.Syncing
		}
	}

	MapSetNode(node)

	v, _ := json.Marshal(node)
	_, _ = conn.SetString(node.PubKey, v)
}

func (p *Node) Delete() {
	defer func() {
		if err:=recover(); err != nil {
			logger.Error(err)
		}
	}()
	_, _ = db.GetRedisPool(NODE_POOL_NAME).DelKey(p.PubKey)
}


func (p *Node) Check()  {
	if len(p.RPC) == 0 {
		return
	}

	startTime := utils.CurrentTimestampMilli()

	jsonReq := "{\"jsonrpc\":\"2.0\",\"id\":\"blockpi-drpc\", \"method\":\"getHealth\",\"params\": []}"
	data := http.PostJson(p.GetRPCEndpoint(), jsonReq)
	if data != nil {
		endTime := utils.CurrentTimestampMilli()
		resp := rpc.JSONRpcResponse{}
		err := json.Unmarshal(data, &resp)
		if err != nil {
			goto endCheck
		}

		p.LastOnlineTime = endTime
		p.Latency = endTime - startTime

		if resp.Error != nil {
			p.Healthy = false
			numSlotsBehind := resp.Error.Data.(map[string]interface{})["numSlotsBehind"]
			if numSlotsBehind != nil && numSlotsBehind.(float64) > 0{
				p.Syncing = true
				goto checkHeight
			}

			goto endCheck
		}

		p.Syncing = false
		p.Healthy = true
checkHeight:
		jsonReq = "{\"jsonrpc\":\"2.0\",\"id\":\"blockpi-drpc\", \"method\":\"getBlockHeight\",\"params\": []}"
		if data := http.PostJson(p.GetRPCEndpoint(), jsonReq); data != nil {
			resp := rpc.JSONRpcResponse{}
			if err := json.Unmarshal(data, &resp); err == nil && resp.Error == nil {
				p.Height = int64(resp.Result.(float64))
			}
		}
	}

endCheck:
	p.LastCheckTime = utils.CurrentTimestampMilli()
	p.Checking = false

	if !p.Healthy && !p.Syncing && p.LastOnlineTime > 0 {

		if utils.CurrentTimestampMilli() - p.LastOnlineTime > MAX_NODE_RETENTION_TIME {
			p.Delete()
			logger.Debug("Delete node ", p.RPC)
		}
	}

	if p.Healthy || p.Syncing {
		SaveNode(p, false)
	}
}
