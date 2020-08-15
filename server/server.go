package server

import (
	"avrilko-rpc/protocol"
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"
)

// 核心服务类
type Server struct {
	ln           net.Listener  // 全局唯一的监听（可以多路复用）
	readTimeout  time.Duration // 读超时
	writeTimeout time.Duration // 写超时

	gatewayHttpServer     *http.Server // 当启用http网关时候被挂载
	disableHTTPGateway    bool         // 是否禁用http网关服务（开启时候方便测试和调试rpc服务）
	disableJSONRPCGateway bool         // 是否禁用json rpc网关服务 

	serviceMapMu sync.RWMutex        // 服务提供者map读写锁
	serviceMap   map[string]*service // 服务提供者集合map

	connMu     sync.RWMutex          // 各个活跃连接读写锁
	activeConn map[net.Conn]struct{} // 每个活跃的连接，map结构防止重复
	doneChan   chan struct{}         // 服务结束chan

	inShutdown bool               //服务是否关闭 1为关闭 0为正在运行
	onShutdown []func(s *service) // 服务结束后执行的钩子函数

	tlsConfig *tls.Config // tls证书配置

	Plugins PluginContainer // 插件容器（设计核心）

	AuthFunc func(ctx context.Context, request *protocol.Message, token string) error // 认证函数

	handlerMsgNum int32 // 正在处理的消息数量
}

// 初始化服务
func NewServer(opts ...OptionFunc) *Server {
	server := &Server{
		serviceMapMu: sync.RWMutex{},
		serviceMap:   make(map[string]*service),
		activeConn:   make(map[net.Conn]struct{}),
		doneChan:     make(chan struct{}),
		Plugins:      &pluginContainer{},
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(server)
		}
	}

	return server
}
