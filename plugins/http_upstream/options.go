package http_upstream

import (
	"time"

	"github.com/blockpilabs/solana-drpc/common"
)

type httpUpstreamMiddlewareOptions struct {
	defaultTargetEndpoint string
	upstreamTimeout       time.Duration
}

func HttpDefaultTargetEndpoint(endpoint string) common.Option {
	return func(options common.Options) {
		mOptions := options.(*httpUpstreamMiddlewareOptions)
		mOptions.defaultTargetEndpoint = endpoint
	}
}

func HttpUpstreamTimeout(timeout time.Duration) common.Option {
	return func(options common.Options) {
		mOptions := options.(*httpUpstreamMiddlewareOptions)
		mOptions.upstreamTimeout = timeout
	}
}
