package client

import (
	"avrilko-rpc/protocol"
	"context"
	"crypto/tls"
	"time"
)

type Call struct {
	ServicePath   string
	ServiceMethod string
	Metadata      map[string]string // 需要传递到服务端的元数据
	ResMetadata   map[string]string // 接受服务端返回的元数据
	request       interface{}       // 客户端请求
	response      interface{}       // 响应
	Error         error             // 调用完成之后的错误
	Done          chan *Call        // 调用完整之后会讲数据塞到此通道中
	Raw           bool              // 是否发送原始数据
}

// 熔断器接口
type Breaker interface {
	Call(func() error, time.Duration) error
	Fail()
	Success()
	Ready() bool
}

// rpc 单个客户端
type RPCClient interface {
	Connect(network, address string) error                                                                           // 建立连接
	Go(ctx context.Context, servicePath, serviceMethod string, request, response interface{}, done chan *Call) *Call // 异步请求
	Call(ctx context.Context, servicePath, serviceMethod string, request, response interface{}) error                // 同步请求
	SendRaw(ctx context.Context, r *protocol.Message) (map[string]string, []byte, error)                             // 发送原始数据
	Close() error                                                                                                    // 关闭

	RegisterServerMessageChan(ch chan<- *protocol.Message) // 注册消息通道
	UnregisterServerMessageChan()                          // 卸载消息通道

	IsClosing() bool  // 是否正在关闭
	IsShutDown() bool // 是否已经关闭
}

type Option struct {
	Group string // 选择分组，相同组会被优先选择到，不设置将会不生效

	Retries int // 客户端重试次数

	TLSConfig *tls.Config // 证书设置

	RPCPath string // http 链接时候rpc的默认路径

	ConnectTimeout time.Duration // 连接超时时间

	ReadTimeout time.Duration // 客户端读取数据超时时间

	WriteTimeout time.Duration // 客户端写数据超时时间

	BackupLatency time.Duration // 这个用于Failbackup模式，发起请求BackupLatency后没有响应则向另外一台服务器发起请求，看看谁最新返回取谁的数据

	GenBreaker func() Breaker // 用来配置CircuitBreaker

	SerializeType protocol.SerializeType // 序列化的方式

	CompressType protocol.CompressType // 压缩的方式

	Heartbeat bool // 是否启用心跳检测

	HeartbeatInterval time.Duration // 心跳检测的间隔时间
}
