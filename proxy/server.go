package proxy

import (
	"github.com/blockpilabs/solana-drpc/log"
	"github.com/blockpilabs/solana-drpc/plugins/http_upstream"
	"github.com/blockpilabs/solana-drpc/plugins/statistic"
	"github.com/blockpilabs/solana-drpc/providers"
	"github.com/blockpilabs/solana-drpc/rpc"
	"github.com/gorilla/websocket"
)

var logger = log.GetLogger("server")


type ProxyServer struct {
	Backend			*http_upstream.HttpUpstreamMiddleware
	Statistic		*statistic.StatisticMiddleware
	Provider        providers.RpcProvider
}

/**
 * NewProxyServer: init and return a new proxy server instance
 */
func NewProxyServer(provider providers.RpcProvider) *ProxyServer {
	server := &ProxyServer{
		Backend: 		http_upstream.NewHttpUpstreamMiddleware(),
		Statistic:		statistic.NewStatisticMiddleware(),
		Provider:        provider,
	}
	return server
}

func (server *ProxyServer) OnConnection(connSession *rpc.ConnectionSession) error {
	server.Backend.OnConnection(connSession)
	logger.Debug("OnConnection")
	return nil
}

func (server *ProxyServer) OnConnectionClosed(connSession *rpc.ConnectionSession) error {
	// must ensure middleware chain not change after calling OnConnection,
	// otherwise some removed middlewares may not call OnConnectionClosed
	server.Backend.OnConnectionClosed(connSession)
	logger.Debug("OnConnectionClosed")
	return nil
}

func (server *ProxyServer) OnRpcRequest(connSession *rpc.ConnectionSession, rpcSession *rpc.JSONRpcRequestSession) (err error) {
	logger.Debug("OnRpcRequest")
	server.Backend.OnRpcRequest(rpcSession)
	server.Statistic.OnRpcRequest(rpcSession)

	go func() {
		server.Backend.ProcessRpcRequest(rpcSession)
		if err != nil {
			logger.Warn("ProcessRpcRequest error", err)
			return
		}

		rpcRes := rpcSession.Response
		if rpcRes == nil {
			logger.Error("empty jsonrpc response, maybe no valid middleware added")
			return
		}

		server.Backend.OnRpcResponse(rpcSession)

		resBytes, err := rpc.EncodeJSONRPCResponse(rpcRes)
		if err != nil {
			logger.Error("encodeJSONRPCResponse err", err)
			return
		}


		connSession.RequestConnectionWriteChan <- rpc.NewMessagePack(websocket.TextMessage, resBytes)
	}()
	return
}

func (server *ProxyServer) SetBackends(backends map[string]int64) {
	server.SetBackendsByGroup("default", backends)
}

func (server *ProxyServer) SetBackendsByGroup(group string, backends map[string]int64) {
	for endpoint,weight := range backends {
		server.Backend.AddRpcEndPoint(group, endpoint, weight)
	}
}

func (server *ProxyServer) Start() {
	if server.Provider == nil {
		logger.Fatalln("please set provider to ProxyServer before start")
		return
	}
	server.Provider.SetRpcProcessor(server)
	logger.Fatal(server.Provider.ListenAndServe())
}

func (server *ProxyServer) Close() {
	logger.Debug("ProxyServer Close")
}
