package load_balancer

import (
	"errors"
	"sync"

	"github.com/blockpilabs/solana-drpc/log"
	"github.com/blockpilabs/solana-drpc/rpc"
)

var logger = log.GetLogger("load_balancer")

type UpstreamItem struct {
	Id             int64
	TargetEndpoint string
	Weight         int64
}

var upstreamItemIdGen int64 = 0

func NewUpstreamItem(targetEndpoint string, weight int64) *UpstreamItem {
	defer func() {
		upstreamItemIdGen++
	}()
	return &UpstreamItem{
		Id:             upstreamItemIdGen,
		TargetEndpoint: targetEndpoint,
		Weight:         weight,
	}
}

type LoadBalanceMiddleware struct {
	selector      map[string]*WrrSelector
	UpstreamItems map[string][]*UpstreamItem

	mutex sync.Mutex
}

func NewLoadBalanceMiddleware() *LoadBalanceMiddleware {
	return &LoadBalanceMiddleware{
		selector: make(map[string]*WrrSelector),
		UpstreamItems: make(map[string][]*UpstreamItem),
	}
}

func (middleware *LoadBalanceMiddleware) AddUpstreamItem(groupName string, item *UpstreamItem) *LoadBalanceMiddleware {
	middleware.mutex.Lock()
	defer middleware.mutex.Unlock()
	if _, ok := middleware.selector[groupName]; !ok {
		middleware.selector[groupName] = &WrrSelector{}
		middleware.UpstreamItems[groupName] = []*UpstreamItem{}
	}

	middleware.UpstreamItems[groupName] = append(middleware.UpstreamItems[groupName], item)
	middleware.selector[groupName].AddNode(item.Weight, item)
	return middleware
}

func (middleware *LoadBalanceMiddleware) Clear() *LoadBalanceMiddleware {
	middleware.mutex.Lock()
	defer middleware.mutex.Unlock()
	for groupName, _ := range middleware.UpstreamItems {
		middleware.UpstreamItems[groupName] = []*UpstreamItem{}
	}
	for groupName, _ := range middleware.selector {
		middleware.selector[groupName].Clear()
	}
	//middleware.UpstreamItems = []*UpstreamItem{}
	//middleware.selector.Clear()
	return middleware
}

func (middleware *LoadBalanceMiddleware) Name() string {
	return "load_balance"
}

func (middleware *LoadBalanceMiddleware) selectTargetByWeight(groupName string) *UpstreamItem {
	selected, err := middleware.selector[groupName].Next()
	if err != nil {
		logger.Error("load balance selector next error", err)
		return nil
	}
	selectedUpStreamItem, ok := selected.(*UpstreamItem)
	if !ok {
		return nil
	}

	logger.Debugf("[load-balancer]select target %s", selectedUpStreamItem.TargetEndpoint)
	return selectedUpStreamItem
}

func (middleware *LoadBalanceMiddleware) OnStart() (err error) {
	return nil
}

func (middleware *LoadBalanceMiddleware) OnConnection(session *rpc.ConnectionSession) (err error) {
	return nil
}

func (middleware *LoadBalanceMiddleware) OnConnectionClosed(session *rpc.ConnectionSession) (err error) {
	return nil
}

func (middleware *LoadBalanceMiddleware) OnRpcRequest(session *rpc.JSONRpcRequestSession) (err error) {
	group := "default"

	methodInfo := GetMethodInfo(session.Request.Method)
	if methodInfo != nil {
		group = methodInfo.Group
	}

	selectedTargetItem := middleware.selectTargetByWeight(group)
	if selectedTargetItem == nil {
		err = errors.New("can't select one upstream target")
		return
	}
	logger.Debugf("selected upstream target item id#%d endpoint: %s\n", selectedTargetItem.Id, selectedTargetItem.TargetEndpoint)
	session.Conn.SelectedUpstreamTarget = &selectedTargetItem.TargetEndpoint

	return nil
}

func (middleware *LoadBalanceMiddleware) OnRpcResponse(session *rpc.JSONRpcRequestSession) (err error) {
	return nil
}

func (middleware *LoadBalanceMiddleware) ProcessRpcRequest(session *rpc.JSONRpcRequestSession) (err error) {
	return nil
}
