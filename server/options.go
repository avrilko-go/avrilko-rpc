package server

import (
	"avrilko-rpc/middleware"
	"time"
)

type OptionFuc func(options *Options)

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

// 服务监听地址
func WithAddress(address string) OptionFuc {
	return func(option *Options) {
		option.address = address
	}
}

// rpc传输协议
func WithProtocol(protocol string) OptionFuc {
	return func(options *Options) {
		options.protocol = protocol
	}
}

// 序列化方案
func WithSerializationType(serializationType string) OptionFuc {
	return func(options *Options) {
		options.serializationType = serializationType
	}
}

// 超时时间
func WithTimeout(timeout time.Duration) OptionFuc {
	return func(options *Options) {
		options.timeout = timeout
	}
}

// 服务发现地址
func WithSelectorAddr(address string) OptionFuc {
	return func(options *Options) {
		options.selectorAddr = address
	}
}

// 链路追踪的地址
func WithTracingAddr(address string) OptionFuc {
	return func(options *Options) {
		options.tracingAddr = address
	}
}

// 链路追踪span名称
func WithTracingSpanName(spanName string) OptionFuc {
	return func(options *Options) {
		options.tracingSpanName = spanName
	}
}

// 注册插件
func WithPluginName(pluginName ...string) OptionFuc {
	return func(options *Options) {
		options.pluginName = append(options.pluginName, pluginName...)
	}
}

// 注册中间件
func WithMiddlewares(middlewares ...middleware.ServerMiddleware) OptionFuc {
	return func(options *Options) {
		options.middlewares = append(options.middlewares, middlewares...)
	}
}
