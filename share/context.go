package share

import (
	"context"
	"fmt"
	"reflect"
)

// 自定义上下文默认实现context.Context
type Context struct {
	context.Context // 这是一个接口
	maps            map[interface{}]interface{}
}

func NewContext(ctx context.Context) *Context {
	return &Context{
		Context: ctx,
		maps:    make(map[interface{}]interface{}),
	}
}

func (c *Context) Value(key interface{}) interface{} {
	if c.maps == nil {
		c.maps = make(map[interface{}]interface{})
	}
	if v, ok := c.maps[key]; ok {
		return v
	}

	return c.Context.Value(key)
}

func (c *Context) SetValue(key, value interface{}) {
	if c.maps == nil {
		c.maps = make(map[interface{}]interface{})
	}
	c.maps[key] = value
}

func (c *Context) String() string {
	return fmt.Sprintf("%v.WithValue(%v)", c.Context, c.maps)
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

func WithLocalValue(ctx *Context, key, value interface{}) *Context {
	if key == nil {
		panic("context key 传递非法")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("传入的key必须是可以比较大小的")
	}
	if ctx.maps == nil {
		ctx.maps = make(map[interface{}]interface{}, 10)
	}

	ctx.maps[key] = value
	return ctx
}
