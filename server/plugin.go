package server

import (
	"avrilko-rpc/protocol"
	"avrilko-rpc/share"
	"context"
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
	DoPreReadRequest(ctx context.Context) error                                      // req数据转换为protocol.Message前调用
	DoPostReadRequest(ctx context.Context, message *protocol.Message, e error) error // req数据转换为protocol.Message后调用

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
	return nil
}

func (p *pluginContainer) DoRegisterFunction(name, funcName string, funcObject interface{}, metadata string) error {
	return nil
}

func (p *pluginContainer) DoUnregister(name string) error {
	return nil
}

func (p *pluginContainer) DoPostConnAccept(conn net.Conn) (net.Conn, bool) {
	return conn, true
}

func (p *pluginContainer) DoPostConnClose(conn net.Conn) bool {
	return true
}

func (p *pluginContainer) DoPreReadRequest(ctx context.Context) error {
	return nil
}

func (p *pluginContainer) DoPostReadRequest(ctx context.Context, message *protocol.Message, e error) error {
	return nil
}

func (p *pluginContainer) DoPreHandleRequest(ctx context.Context, message *protocol.Message) error {
	return nil
}

func (p *pluginContainer) DoPreCall(ctx context.Context, serviceName, serviceMethod string, request interface{}) (interface{}, error) {
	return request, nil
}

func (p *pluginContainer) DoPostCall(ctx context.Context, serviceName, serviceMethod string, request interface{}, response interface{}) (interface{}, error) {
	return response, nil
}

func (p *pluginContainer) DoPreWriteResponse(ctx context.Context, request, response *protocol.Message) error {
	s := ctx.(*share.Context)
	m := make(map[string]string)
	m["name"] = "avrilko"
	s.SetValue(share.ResMetaDataKey, m)
	return nil
}

func (p *pluginContainer) DoPostWriteResponse(ctx context.Context, request, response *protocol.Message, err error) error {
	return nil
}
