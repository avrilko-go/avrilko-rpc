package client

import "context"

// key value 键值对
type KVPair struct {
	Key   string
	Value string
}

type XClient interface {
	SetPlugins(plugins PluginContainer)
	GetPlugin() PluginContainer
	SetSelector(s Selector)
	Auth(auth string)

	Go(ctx context.Context, serviceMethod string, request interface{}, response interface{}, done chan *Call) (*Call, error) // 协程话调用Call
	Call(ctx context.Context, serviceMethod string, request interface{}, response interface{}) error                         // 同步调用call
	Broadcast(ctx context.Context, serviceMethod string, request interface{}, response interface{}) error                    // 广播发消息(直到所有服务端接收到消息)
	Fork(ctx context.Context, serviceMethod string, request interface{}, response interface{}) error                         // 广播发消息(有一个服务端响应即可)
}

type ServiceDiscoveryFilter func(kvp *KVPair) bool

type ServiceDiscovery interface {
	GetServices() []*KVPair                    // 获取所有服务发现的服务器
	WatchService() chan []*KVPair              // 监听服务发现的变动信息
	RemoveService(ch chan []*KVPair)           // 移除服务发现的变动信息
	Clone(servicePath string) ServiceDiscovery // 克隆一个服务发现的对象
	SetFilter(ServiceDiscoveryFilter)          // 设置过滤器
	Close()                                    // 关闭服务发现
}
