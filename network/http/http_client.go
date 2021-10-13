package http

import (
	"time"

	"github.com/blockpilabs/solana-drpc/log"
	"github.com/valyala/fasthttp"
)

const TIMEOUT = time.Second * 10

var logger = log.GetLogger("httpclient")

func PostJson(url string, jsonStr string) []byte {
	req := &fasthttp.Request{}
	req.SetRequestURI(url)
	requestBody := []byte(jsonStr)
	req.SetBody(requestBody)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod("POST")
	resp := &fasthttp.Response{}
	client := &fasthttp.Client{}
	if err := client.DoTimeout(req, resp, TIMEOUT);err != nil {
		//logger.Debug(err)
		return nil
	}
	return resp.Body()
}