package providers

import (
	"github.com/blockpilabs/solana-drpc/log"
	"github.com/blockpilabs/solana-drpc/rpc"
)

var logger = log.GetLogger("provider")

type RpcProviderProcessor interface {
	OnConnection(connSession *rpc.ConnectionSession) error
	OnConnectionClosed(connSession *rpc.ConnectionSession) error
	OnRpcRequest(connSession *rpc.ConnectionSession, rpcSession *rpc.JSONRpcRequestSession) error
}

type RpcProvider interface {
	SetRpcProcessor(processor RpcProviderProcessor)
	ListenAndServe() error
}
