package server

import (
	"reflect"
	"sync"
)

var usePool bool

var ObjectPool = objectPool{
	pools: make(map[reflect.Type]*sync.Pool),
	New: func(p reflect.Type) interface{} {
		if p.Kind() == reflect.Ptr {
			p = p.Elem()
		}
		return reflect.New(p).Interface()
	},
}

// 如果全局开启对象链接池，则注入到链接池的对象必须实现此接口，不然会有变量污染
type ObjectReset interface {
	OReset() // 重置对象的属性为默认值
}

type objectPool struct {
	sync.RWMutex
	pools map[reflect.Type]*sync.Pool
	New   func(p reflect.Type) interface{}
}

func (o *objectPool) Init(p reflect.Type) {
	sPool := &sync.Pool{}
	sPool.New = func() interface{} {
		return o.New(p)
	}

	o.Lock()
	defer o.Unlock()

	o.pools[p] = sPool
}

func (o *objectPool) Get(p reflect.Type) interface{} {
	if !usePool {
		return o.New(p)
	}
	o.RLock()
	defer o.RUnlock()
	pool := o.pools[p]
	return pool.Get()
}

func (o *objectPool) Put(p reflect.Type, data interface{}) {
	if !usePool {
		return
	}

	if oReset, ok := data.(ObjectReset); ok {
		oReset.OReset()
	}

	o.RLock()
	defer o.RUnlock()
	pool := o.pools[p]
	pool.Put(data)
}
