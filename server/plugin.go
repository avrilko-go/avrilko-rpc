package server

import (
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
	DoPreReadRequest(ctx context.Context) error // req数据转换为protocol.Message前调用
	DoPostReadRequest(ctx context.Context)
}

type Plugin struct {
}
