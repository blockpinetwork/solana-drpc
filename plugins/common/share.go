package common

import (
	"errors"

	"github.com/blockpilabs/solana-drpc/rpc"
)

func GetSessionStringParam(session *rpc.JSONRpcRequestSession, paramName string, defaultValue *string) (result string, err error) {
	paramValue, ok := session.Parameters[paramName]
	var defaultValueStr = ""
	if defaultValue != nil {
		defaultValueStr = *defaultValue
	}
	if !ok {
		result = defaultValueStr
		return
	}
	result, ok = paramValue.(string)
	if !ok {
		err = errors.New("invalid " + paramName + " param in session")
		return
	}
	return
}

func SetSelectedUpstreamTargetEndpoint(session *rpc.ConnectionSession, value string) (err error) {
	session.SelectedUpstreamTarget = &value
	return
}

func GetSelectedUpstreamTargetEndpoint(session *rpc.ConnectionSession, defaultValue *string) (result string, err error) {
	var defaultValueStr = ""
	if defaultValue != nil {
		defaultValueStr = *defaultValue
	}
	if session.SelectedUpstreamTarget == nil {
		result = defaultValueStr
		return
	}
	result = *session.SelectedUpstreamTarget
	return
}
