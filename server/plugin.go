package server

import (
	"avrilko-rpc/protocol"
	"context"
	"fmt"
	"net"
)

// 插件接口(在程序不同生命周期时切入程序，实现不同的功能)
type PluginContainer interface {
	Add(plugin Plugin)    // 添加插件
	Remove(plugin Plugin) // 移除插件
	All(plugin Plugin)    // 获取所有插件

	// 注册相关周期
	DoRegister(name string, object interface{}, metadata string) error                       // 反射注册对象时调用
	DoRegisterFunction(name, funcName string, funcObject interface{}, metadata string) error // 直接注册函数时调用
	DoUnregister(name string) error                                                          // 反注册时调用

	// 连接相关周期
	DoPostConnAccept(conn net.Conn) (net.Conn, bool) // 连接被listen accept后调用
	DoPostConnClose(conn net.Conn) bool              // 连接被关闭后调用

	// 读数据周期
	DoPreReadRequest(ctx context.Context) error                                // req数据转换为protocol.Message前调用
	DoPostReadRequest(ctx context.Context, message *protocol.Message, e error) // req数据转换为protocol.Message后调用

	// 处理请求周期
	DoPreHandleRequest(ctx context.Context, message *protocol.Message) error                                                           // 处理请求前（路由查找前）调用
	DoPreCall(ctx context.Context, serviceName, serviceMethod string, request interface{}) (interface{}, error)                        // 调用自定义方法前调用
	DoPostCall(ctx context.Context, serviceName, serviceMethod string, request interface{}, response interface{}) (interface{}, error) // 调用自定义方法后调用

	// 写数据相关周期
	DoPreWriteResponse(ctx context.Context, request, response *protocol.Message) error             // 写入数据之前调用
	DoPostWriteResponse(ctx context.Context, request, response *protocol.Message, err error) error // 写入数据之后调用
}

// 最原始的plugin(可以是任何类型)
type Plugin interface {
}

// PluginContainer默认实现
type pluginContainer struct {
	plugin []Plugin
}

func (p *pluginContainer) Add(plugin Plugin) {
	panic("implement me")
}

func (p *pluginContainer) Remove(plugin Plugin) {
	panic("implement me")
}

func (p *pluginContainer) All(plugin Plugin) {
	panic("implement me")
}

func (p *pluginContainer) DoRegister(name string, object interface{}, metadata string) error {
	fmt.Println(name)
	fmt.Println(object)
	fmt.Println(metadata)
	return nil
}

func (p *pluginContainer) DoRegisterFunction(name, funcName string, funcObject interface{}, metadata string) error {
	panic("implement me")
}

func (p *pluginContainer) DoUnregister(name string) error {
	panic("implement me")
}

func (p *pluginContainer) DoPostConnAccept(conn net.Conn) (net.Conn, bool) {
	panic("implement me")
}

func (p *pluginContainer) DoPostConnClose(conn net.Conn) bool {
	panic("implement me")
}

func (p *pluginContainer) DoPreReadRequest(ctx context.Context) error {
	panic("implement me")
}

func (p *pluginContainer) DoPostReadRequest(ctx context.Context, message *protocol.Message, e error) {
	panic("implement me")
}

func (p *pluginContainer) DoPreHandleRequest(ctx context.Context, message *protocol.Message) error {
	panic("implement me")
}

func (p *pluginContainer) DoPreCall(ctx context.Context, serviceName, serviceMethod string, request interface{}) (interface{}, error) {
	panic("implement me")
}

func (p *pluginContainer) DoPostCall(ctx context.Context, serviceName, serviceMethod string, request interface{}, response interface{}) (interface{}, error) {
	panic("implement me")
}

func (p *pluginContainer) DoPreWriteResponse(ctx context.Context, request, response *protocol.Message) error {
	panic("implement me")
}

func (p *pluginContainer) DoPostWriteResponse(ctx context.Context, request, response *protocol.Message, err error) error {
	panic("implement me")
}
