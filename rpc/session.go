package rpc

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type JSONRpcRequestBundle struct {
	MessageType        int
	Data               []byte
	Request            *JSONRpcRequest
	ResponseFutureChan chan *JSONRpcResponse
}

func NewJSONRpcRequestBundle(messageType int, data []byte, rpcRequest *JSONRpcRequest, rpcResponseFutureChan chan *JSONRpcResponse) *JSONRpcRequestBundle {
	return &JSONRpcRequestBundle{
		MessageType:        messageType,
		Data:               data,
		Request:            rpcRequest,
		ResponseFutureChan: rpcResponseFutureChan,
	}
}

type MessagePack struct {
	MessageType int
	Message     []byte
}

func NewMessagePack(messageType int, message []byte) *MessagePack {
	return &MessagePack{
		MessageType: messageType,
		Message:     message,
	}
}

type ConnectionSession struct {
	RequestConnection          *websocket.Conn
	HttpResponse               http.ResponseWriter
	HttpRequest                *http.Request
	ConnectionDone             chan struct{}
	RequestConnectionWriteChan chan *MessagePack

	// rpc fields shared in connection session
	RpcRequestsMap             map[interface{}]chan *JSONRpcResponse // rpc request id => channel notify of *JSONRpcResponse
	RpcRequestsDispatchChannel chan *RequestDispatchData        // rpc requests dispatching

	// same base middleware shared fields

	// upstream middleware shared fields in connection session
	SelectedUpstreamTarget       *string

	UpstreamTargetConnection     *websocket.Conn
	UpstreamTargetConnectionDone chan struct{}
	UpstreamRpcRequestsChan      chan *JSONRpcRequestBundle
}

func (connSession *ConnectionSession) Close() {
	close(connSession.ConnectionDone)
	connSession.ConnectionDone = nil
	close(connSession.RequestConnectionWriteChan)
	connSession.RequestConnectionWriteChan = nil
	close(connSession.RpcRequestsDispatchChannel)
	connSession.RpcRequestsDispatchChannel = nil
}

func NewConnectionSession() *ConnectionSession {
	return &ConnectionSession{
		RequestConnectionWriteChan: make(chan *MessagePack, 1000),
		ConnectionDone:             make(chan struct{}),
		RpcRequestsMap:             make(map[interface{}]chan *JSONRpcResponse),
		RpcRequestsDispatchChannel: make(chan *RequestDispatchData, 1000),
	}
}

type JSONRpcRequestSession struct {
	Conn         *ConnectionSession
	RequestBytes []byte
	Request      *JSONRpcRequest
	Response     *JSONRpcResponse
	Parameters   map[string]interface{}

	RpcResponseFutureChan chan *JSONRpcResponse

	// cache middleware shared fields
	ResponseSetByCache bool

	// before_cache middleware shared fields
	MethodNameForCache *string // only used to find cache in cache middleware

	// selected upstream target server url
	TargetServer string
}

func NewJSONRpcRequestSession(conn *ConnectionSession) *JSONRpcRequestSession {
	return &JSONRpcRequestSession{
		Conn:               conn,
		ResponseSetByCache: false,
	}
}

func (requestSession *JSONRpcRequestSession) FillRpcRequest(request *JSONRpcRequest, requestBytes []byte) {
	requestSession.Request = request
	requestSession.RequestBytes = requestBytes
}

func (requestSession *JSONRpcRequestSession) FillRpcResponse(response *JSONRpcResponse) {
	requestSession.Response = response
}
