package http

import "github.com/valyala/fasthttp"

type HttpServerHandler func(ctx *fasthttp.RequestCtx)

type HttpServer struct {
	Endpoint     string
	Handler fasthttp.RequestHandler
}

func (server *HttpServer)ListenAndServe() error {
	s := &fasthttp.Server{Handler: server.Handler}
	return s.ListenAndServe(server.Endpoint)
}
