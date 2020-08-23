package client

import (
	"avrilko-rpc/protocol"
	"context"
	"net"
)

type Plugin interface {
}

// 客户端插件容器
type PluginContainer interface {
	Add(plugin Plugin)    // 添加一个插件
	Remove(plugin Plugin) // 移除一个插件
	All() []Plugin        // 获取所有插件

	DoConnCreated(conn net.Conn) (net.Conn, error)     // 在conn 创建之后执行
	DoClientConnected(conn net.Conn) (net.Conn, error) // 在客户端链接之后执行
	DoClientConnectClose(conn net.Conn) error          // 在客户端链接被关闭后执行

	DoPreCall(ctx context.Context, servicePath, serviceMethod string, request interface{}) error                       // 在执行调用前执行
	DoPostCall(ctx context.Context, servicePath, serviceMethod string, request, response interface{}, err error) error // 在执行调用后执行

	DoClientBeforeEncode(message *protocol.Message) error // 在客户端打包数据前调用
	DoClientAfterDecode(message protocol.Message) error   // 在客户端调用解包数据后调用

	DoWrapSelect(selectFunc SelectFunc) SelectFunc // 包装客户端负载均衡算法
}
