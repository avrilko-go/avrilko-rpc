package server

import (
	"reflect"
	"sync"
)

// 反射方法得到的摘要
type methodType struct {
	sync.Mutex                  // 互斥锁
	rMethod      reflect.Method // 反射方法
	requestType  reflect.Type   // 方法请求类型
	responseType reflect.Type   // 方法返回类型
}

// 反射函数得到的摘要
type funcType struct {
	sync.Mutex                 // 互斥锁
	rFuc         reflect.Value // 反射函数
	requestType  reflect.Type  // 方法请求类型
	responseType reflect.Type  // 方法返回类型
}

// 单个服务提供者
type service struct {
	name     string                 // 服务提供者名称
	rValue   reflect.Value          // 反射值
	rType    reflect.Type           // 反射类型
	method   map[string]*methodType // 反射方法集合
	function map[string]*funcType   // 反射函数集合
}
