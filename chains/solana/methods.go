package solana

import "github.com/blockpilabs/solana-drpc/plugins/load_balancer"

var (
	methodMap = make(map[string]*load_balancer.Method)
)


const (
	MethodType_WRITE    = 0
	MethodType_READ     = 1

	MethodGroup_INDEX	= "INDEX"
)


func init() {
	load_balancer.AddMethodInfo( "getProgramAccounts", MethodType_READ, MethodGroup_INDEX)
}
