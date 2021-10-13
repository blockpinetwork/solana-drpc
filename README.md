# SOLANA-DRPC
## A Decentralized RPC service for Solana
Partially open source

Solana DRPC is an innovative Web3 based infrastructure,  establish a decentralized RPC service for Solana network.

The service will automatically scan the network and try to find out all the RPC nodes. Sort the RPC nodes and establish a DRPC proxy server with these nodes which have the lowest latency. Nodes are auto reweighted by latency and kicked out by status if not synced up. All the RPC requests will be distributed to different backends with fastest response.

Solana DRPC will give users a better experience and help developers to "Keep services always online.''


# Features
- [x] Auto rpc nodes discovery
- [x] Backends are sorted by latency
- [x] Upstream: http
- [x] Load-balance: WeightedRound-Robin algorithm implement
- [x] Route rpc methods to specified endpoints
- [x] Dashboard
- [ ] Upstream: websocket and grpc support
- [ ] Cache requests
- [ ] Offline requests


# Usage
## Getting Started
### Dependencies
- Install [Go](https://golang.org/doc/install)

### Building from source
```sh
git clone https://github.com/blockpilabs/solana-drpc.git
go build
```

<div style="page-break-after: always;"></div>

## DRPC status api
```sh
http://localhost:8181/api/status
```

## DRPC proxy endpoint
```sh
http://localhost:9191/
```
The backend upstream is set in the response headers.
```sh
curl \
  --data '{"method":"getHealth","params":[],"id":1,"jsonrpc":"2.0"}' \
  -H "Content-Type: application/json" \
  -X POST http://localhost:9191/ \
  -D -

HTTP/1.1 200 OK
Access-Control-Allow-Methods: GET, POST, PATCH, PUT, DELETE, OPTIONS
Access-Control-Allow-Origin: *
Content-Type: application/json
Upstream: http://54.169.xxx.xxx:8899
Date: Wed, 13 Oct 2021 03:26:32 GMT
Transfer-Encoding: chunked

{
    "id": 1,
    "jsonrpc": "2.0",
    "result": "ok"
}
```
