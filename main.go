package main

import (
	"github.com/blockpilabs/solana-drpc/chains/solana"
	"github.com/blockpilabs/solana-drpc/dashboard"
	"github.com/blockpilabs/solana-drpc/log"
	"github.com/blockpilabs/solana-drpc/plugins/statistic"
	"github.com/blockpilabs/solana-drpc/providers"
	"github.com/blockpilabs/solana-drpc/proxy"
)

var(
	logger = log.GetLogger("main")
)

func main() {
	solana.InitNodeRedisPools("SOLANA_NODE_", "127.0.0.1:6379","", 0)
	statistic.InitStatisticRedisPool("SOLANA_STATISTIC_", "127.0.0.1:6379","", 1)

	solana.UpdateClusterNodes()
	solana.StartCheckQueue()
	solana.UpdateBestRPCNodes()

	provider := providers.NewHttpJsonRpcProvider(":9191", "/", &providers.HttpJsonRpcProviderOptions{
		TimeoutSeconds: 30,
	})
	server := proxy.NewProxyServer(provider)

	solana.SetProxyServer(server)

	go server.Start()
	dashboard.ListenAndServ(dashboard.NewRouter())
}




