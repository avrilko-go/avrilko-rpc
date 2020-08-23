package client

import "context"

type SelectFunc func(ctx context.Context, servicePath, serviceMethod string, args interface{}) string // 负载均衡函数

// 负载均衡算法实现接口
type Selector interface {
	Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string // 负载均衡选择函数
	UpdateServer(servers map[string]string)                                                 // 更新负载均衡选择
}
