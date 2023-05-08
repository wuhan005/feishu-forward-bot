package main

import (
	"context"
	"net/http"

	"github.com/aliyun/fc-runtime-go-sdk/fc"

	handler "github.com/wuhan005/feishu-forward-bot/api"
)

func main() {
	fc.StartHttp(HandleHttpRequest)
}

func HandleHttpRequest(_ context.Context, w http.ResponseWriter, req *http.Request) error {
	handler.Handler(w, req)
	return nil
}
