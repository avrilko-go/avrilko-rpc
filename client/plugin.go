package client

import (
	"avrilko-rpc/protocol"
	"context"
	"net"
)

// 插件容器默认实现
type pluginContainer struct {
	plugins []Plugin
}

func (p pluginContainer) Add(plugin Plugin) {
	panic("implement me")
}

func (p pluginContainer) Remove(plugin Plugin) {
	panic("implement me")
}

func (p pluginContainer) All() []Plugin {
	panic("implement me")
}

func (p pluginContainer) DoConnCreated(conn net.Conn) (net.Conn, error) {
	panic("implement me")
}

func (p pluginContainer) DoClientConnected(conn net.Conn) (net.Conn, error) {
	panic("implement me")
}

func (p pluginContainer) DoClientConnectClose(conn net.Conn) error {
	panic("implement me")
}

func (p pluginContainer) DoPreCall(ctx context.Context, servicePath, serviceMethod string, request interface{}) error {
	panic("implement me")
}

func (p pluginContainer) DoPostCall(ctx context.Context, servicePath, serviceMethod string, request, response interface{}, err error) error {
	panic("implement me")
}

func (p pluginContainer) DoClientBeforeEncode(message *protocol.Message) error {
	panic("implement me")
}

func (p pluginContainer) DoClientAfterDecode(message protocol.Message) error {
	panic("implement me")
}

func (p pluginContainer) DoWrapSelect(selectFunc SelectFunc) SelectFunc {
	panic("implement me")
}

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
