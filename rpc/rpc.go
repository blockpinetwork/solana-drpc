package rpc

import "encoding/json"

const (
	RPC_INTERNAL_ERROR = 10001

	RPC_UPSTREAM_CONNECTION_CLOSED_ERROR = 50001
	RPC_UPSTREAM_TIMEOUT_ERROR           = 50002

	RPC_DISABLED_RPC_METHOD = 60001

	RPC_RESPONSE_TIMEOUT_ERROR = 70001
)

type JSONRpcRequest struct {
	Id      interface{} `json:"id"`
	JSONRpc string      `json:"jsonrpc,omitempty"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

func DecodeJSONRPCRequest(message []byte) (req *JSONRpcRequest, err error) {
	req = new(JSONRpcRequest)
	err = json.Unmarshal(message, &req)
	if err != nil {
		return
	}
	return
}

type JSONRpcResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewJSONRpcResponseError(code int, message string, data interface{}) *JSONRpcResponseError {
	return &JSONRpcResponseError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

type JSONRpcResponse struct {
	Id      interface{}           `json:"id"`
	JSONRpc string                `json:"jsonrpc,omitempty"`
	Error   *JSONRpcResponseError `json:"error,omitempty"`
	Result  interface{}           `json:"result,omitempty"`
}

func NewJSONRpcResponse(id interface{}, result interface{}, err *JSONRpcResponseError) *JSONRpcResponse {
	return &JSONRpcResponse{
		Id:      id,
		JSONRpc: "2.0",
		Error:   err,
		Result:  result,
	}
}

func CloneJSONRpcResponse(source *JSONRpcResponse) (result *JSONRpcResponse, err error) {
	if source == nil {
		result = nil
		return
	}
	bytes, err := json.Marshal(source)
	if err != nil {
		return
	}
	result = new(JSONRpcResponse)
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	return
}

func EncodeJSONRPCResponse(res *JSONRpcResponse) (data []byte, err error) {
	data, err = json.Marshal(res)
	return
}

func DecodeJSONRPCResponse(message []byte) (req *JSONRpcResponse, err error) {
	req = new(JSONRpcResponse)
	err = json.Unmarshal(message, &req)
	if err != nil {
		return
	}
	return
}

// RpcRequestDispatchType: dispatch type of rpc request
type RequestDispatchType int

const (
	RPC_REQUEST_CHANGE_TYPE_ADD_REQUEST  RequestDispatchType = 1 // add rpc request to dispatcher
	RPC_REQUEST_CHANGE_TYPE_ADD_RESPONSE RequestDispatchType = 2 // add rpc response to dispatcher
)

type RequestDispatchData struct {
	Type RequestDispatchType
	Data *JSONRpcRequestSession
}
