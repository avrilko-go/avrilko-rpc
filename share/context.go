package share

import (
	"context"
	"reflect"
)

// 自定义上下文默认实现context.Context
type Context struct {
	context.Context // 这是一个接口
	maps            map[interface{}]interface{}
}

func WithValue(ctx context.Context, key, value interface{}) *Context {
	if key == nil {
		panic("context key 传递非法")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("传入的key必须是可以比较大小的")
	}

	maps := make(map[interface{}]interface{})
	maps[key] = value
	return &Context{
		Context: ctx,
		maps:    maps,
	}
}
