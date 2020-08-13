package service

import (
	"avrilko-rpc/middleware"
	"avrilko-rpc/serilazation"
	"avrilko-rpc/server"
	"context"
	"sync"
)

type IService interface {
	Close(wait *sync.WaitGroup)
	Serve(options *server.Options)
	Register(handlerName string, handlerService HandlerService)
}

// 服务端处理函数（注册服务时一个处理路由对应一个handler）
type HandlerService func(controller interface{}, ctx context.Context, unmarshal serilazation.WrapperUnmarshal, middlewares []middleware.ServerMiddleware) (interface{}, error)
