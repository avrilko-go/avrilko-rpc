package server

import (
	"avrilko-rpc/middleware"
	"time"
)

type Options struct {
	address           string        // 服务监听地址，e.g "127.0.0.1:8080" , "192.168.1.3:9078"
	protocol          string        // rpc传输协议 e.g "http", "tcp(自定义的tcp协议，推荐使用)"
	serializationType string        // 序列化方案 e.g "json", "proto", "msgpack"
	timeout           time.Duration // 超时时间

	selectorAddr    string                        // 服务发现的地址
	tracingAddr     string                        // 链路追踪的地址
	tracingSpanName string                        // 链路追踪span名称
	pluginName      []string                      // 插件名称集合
	middlewares     []middleware.ServerMiddleware // 中间件集合
}
