# SOLANA-DRPC

https://solana-drpc.blockpi.io/

## A Decentralized RPC service for Solana

Partially open source

Solana DRPC is an innovative Web3 based infrastructure,  establish a decentralized RPC service for Solana network.

The service will automatically scan the network and try to find out all the RPC nodes. Sort the RPC nodes and establish a DRPC proxy server with these nodes which have the lowest latency. Nodes are auto reweighted by latency and kicked out by status if not synced up. All the RPC requests will be distributed to different backends with fastest response.

Solana DRPC will give users a better experience and help developers to "Keep services always online".

## Inspiration

Solana provides a very fluent user experience. But RPC is an obstacle in daily use. A better PRC infrastructure is expected both from our own experience and feedback from the community. It ensures that Solana's network provides an improved user and developer experience.

## What it does

BlockPI DRPC service provides a distributed public RPC service network, it discovers the RPC nodes on the network automatically. Users' requests are allocated to the nearest RPC node or a node with the lowest latency. In this way, the RPC request pressure is spread over the entire network, maximizing resource utilization efficiency and improving user experience.

## How we built it

Our team has more than four years of blockchain industry development experience. Team members are proficient in multiple programming languages, especially have rich experience in network and middleware development. We deployed more than one thousand blockchain nodes. The team developed the DRPC system based on our expertise in the blockchain infrastructure industry.

The system is developed by Golang and integrates the components such as Request, Cache, Rate Limit, Load Balance, Upstream Backends, Response, etc. The system currently supports JSON-RPC over HTTP requests and routing of a particular request to a specific RPC node. In the future, it will also support more request methods such as WebSocket, GRPC, and so on.

## Challenges we ran into

During the development of the Solana DRPC renovation, we built RPC nodes by ourselves and found that Solana's RPC nodes required high hardware configurations.  With the gradual development process, we realized that an excellent distributed RPC scheme that can reduce the cost of all developers and users and improve the efficiency of using network resources is very important.



## Accomplishments that we're proud of

We finally finished the first phase of development! Thanks to all the developers who contributed resources and maintenance to Solana's infrastructure, BlockPI Solana DRPC can provide RPC services free of charge to the community.

## What we learned

The Infrastructure of the Web 3.0 should have distinctive innovations along with a more open and inclusive community spirit. We hope that Solana DRPC will benefit all users and developers, empowering them to use the Solana network more smoothly and focus on their own business.

## What's next for BlockPI Solana DRPC

We plan to add more utility functions, such as Cache Response, Load Balance, Offline Request, etc.  Ultimately, we plan to integrate Solana DRPC into BlockPI networks, which will provide additional incentives for all RPC operators to ensure that Solana users and developers can continue to enjoy more stable, high-speed RPC services.



<div style="page-break-after: always;"></div>


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
curl -D - -H "Content-Type: application/json" \
  --data '{"method":"getHealth","params":[],"id":1,"jsonrpc":"2.0"}' \
  -X POST http://localhost:9191/
  
HTTP/2 200 
date: Wed, 13 Oct 2021 03:50:38 GMT
content-type: application/json
content-length: 38
access-control-allow-methods: GET, POST, PATCH, PUT, DELETE, OPTIONS
access-control-allow-origin: *
upstream: http://86.109.15.141:8899

{"id":1,"jsonrpc":"2.0","result":"ok"}
```

```sh
curl -D - -H "Content-Type: application/json" \
  --data '{"method":"getBlockHeight","params":[],"id":1,"jsonrpc":"2.0"}' \
  -X POST http://localhost:9191/
  
HTTP/2 200 
date: Wed, 13 Oct 2021 03:51:55 GMT
content-type: application/json
content-length: 42
access-control-allow-methods: GET, POST, PATCH, PUT, DELETE, OPTIONS
access-control-allow-origin: *
upstream: http://147.75.xx.xxx:8899

{"id":1,"jsonrpc":"2.0","result":90744928}
```

```sh
curl -D - -H "Content-Type: application/json" \
  --data '{"method":"getBalance","params":["ADDRESS"],"id":1,"jsonrpc":"2.0"}' \
  -X POST http://localhost:9191/
  
HTTP/2 200 
date: Wed, 13 Oct 2021 03:53:34 GMT
content-type: application/json
content-length: 81
access-control-allow-methods: GET, POST, PATCH, PUT, DELETE, OPTIONS
access-control-allow-origin: *
upstream: http://203.90.xxx.xxx:8899

{"id":1,"jsonrpc":"2.0","result":{"context":{"slot":101199099},"value":5565072}}
```
