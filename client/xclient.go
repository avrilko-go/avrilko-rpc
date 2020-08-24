package client

import (
	"avrilko-rpc/protocol"
	"context"
	"golang.org/x/sync/singleflight"
	"io"
	"net/url"
	"sync"
)

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
	SendRaw(ctx context.Context, r *protocol.Message) (map[string]string, []byte, error)                                     // 发送原始的数据
	SendFile(ctx context.Context, fileName string, rateInBytesPerSecond int64) error                                         // 发送文件
	DownloadFile(ctx context.Context, requestFileName string, saveTo io.Writer) error                                        // 下载文件
	Close() error                                                                                                            // 关闭服务
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

type xClient struct {
	failMode     FailMode             // 重试模式
	selectMode   SelectMode           // 负载均衡算法
	cachedClient map[string]RPCClient // 客户端缓存
	breakers     sync.Map             // 中断器
	servicePath  string               // 业务
	option       Option               // 客户端选项
	mu           sync.RWMutex         // 读写锁
	servers      map[string]string    // 服务列表
	discovery    ServiceDiscovery     // 服务发现
	selector     Selector             //负载均衡具体实现

	slGroup    singleflight.Group // 防止缓存击穿
	isShutdown bool               // 是否已经停止
	auth       string             // 鉴权验证

	Plugins PluginContainer // 插件容器

	ch                chan []*KVPair           // 服务发现通道
	serverMessageChan chan<- *protocol.Message // 消息通道
}

func (c *xClient) SetPlugins(plugins PluginContainer) {
	panic("implement me")
}

func (c *xClient) GetPlugin() PluginContainer {
	panic("implement me")
}

func (c *xClient) SetSelector(s Selector) {
	panic("implement me")
}

func (c *xClient) Auth(auth string) {
	panic("implement me")
}

func (c *xClient) Go(ctx context.Context, serviceMethod string, request interface{}, response interface{}, done chan *Call) (*Call, error) {
	panic("implement me")
}

func (c *xClient) Call(ctx context.Context, serviceMethod string, request interface{}, response interface{}) error {
	panic("implement me")
}

func (c *xClient) Broadcast(ctx context.Context, serviceMethod string, request interface{}, response interface{}) error {
	panic("implement me")
}

func (c *xClient) Fork(ctx context.Context, serviceMethod string, request interface{}, response interface{}) error {
	panic("implement me")
}

func (c *xClient) SendRaw(ctx context.Context, r *protocol.Message) (map[string]string, []byte, error) {
	panic("implement me")
}

func (c *xClient) SendFile(ctx context.Context, fileName string, rateInBytesPerSecond int64) error {
	panic("implement me")
}

func (c *xClient) DownloadFile(ctx context.Context, requestFileName string, saveTo io.Writer) error {
	panic("implement me")
}

func (c *xClient) Close() error {
	panic("implement me")
}

func (c *xClient) watch(ch chan []*KVPair) {
	for pairs := range ch { // 从通道中一直读取数据
		servers := make(map[string]string, len(pairs))
		for _, p := range pairs {
			servers[p.Key] = p.Value
		}
		c.mu.Lock()
		filterByStateAndGroup(c.option.Group, servers)
		c.servers = servers
		if c.selector != nil {
			c.selector.UpdateServer(servers)
		}

		c.mu.Unlock()
	}
}

func NewXClient(servicePath string, failMode FailMode, selectMode SelectMode, discovery ServiceDiscovery, option Option) XClient {
	client := &xClient{
		failMode:     failMode,
		selectMode:   selectMode,
		cachedClient: make(map[string]RPCClient),
		option:       option,
		mu:           sync.RWMutex{},
		discovery:    discovery,
	}

	pairs := discovery.GetServices()
	servers := make(map[string]string, len(pairs))
	for _, p := range pairs {
		servers[p.Key] = p.Value
	}

	client.servers = servers
	if selectMode != Closest && selectMode != SelectByUser {
		client.selector = newSelector(selectMode, servers)
	}

	client.Plugins = &pluginContainer{}

	ch := client.discovery.WatchService()
	if ch != nil {
		client.ch = ch
		go client.watch(ch)
	}

	return client
}

// 过滤掉指定的servers
func filterByStateAndGroup(group string, servers map[string]string) {
	for k, v := range servers {
		if values, err := url.ParseQuery(v); err == nil {
			if state := values.Get("state"); state == "inactive" {
				delete(servers, k)
			}

			if group != "" && group != values.Get("group") {
				delete(servers, k)
			}
		}
	}
}
