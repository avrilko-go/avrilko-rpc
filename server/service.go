package server

import (
	"avrilko-rpc/log"
	"context"
	"errors"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

var (
	typeContext = reflect.TypeOf((*context.Context)(nil)).Elem()
	errorType   = reflect.TypeOf((*error)(nil)).Elem()
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

// 注册服务提供者(自定义名称)
func (s *Server) RegisterName(name string, object interface{}, metadata string) error {
	_, err := s.register(object, name, true)
	if err != nil {
		return err
	}
	return s.Plugins.DoRegister(name, object, metadata)
}

// 注册服务提供者(类型名称为结构体名称)
func (s *Server) Register(object interface{}, metadata string) error {
	name, err := s.register(object, "", false)
	if err != nil {
		return err
	}
	return s.Plugins.DoRegister(name, object, metadata)
}

// 反射注册服务
func (s *Server) register(object interface{}, name string, useName bool) (string, error) {
	// 读写锁
	s.serviceMapMu.Lock()
	defer s.serviceMapMu.Unlock()

	service := new(service)
	service.rType = reflect.TypeOf(object)
	service.rValue = reflect.ValueOf(object)

	// Indirect 这么做是防止传入的object不是一个指针类型直接取type.Elem()发生panic
	serviceName := reflect.Indirect(service.rValue).Type().Name()
	if useName {
		serviceName = name
	}

	if serviceName == "" { // 没有反射到name值或者外部传入为空字符串则不允许注册
		errStr := "rpc服务提供者注册失败，未传入服务名或者服务不能被反射 " + service.rType.String()
		log.Error(errStr)
		return serviceName, errors.New(errStr)
	}

	if !useName && !isExported(serviceName) { // 外面没有传进来服务提供者名称，反射出来的名称又不能导出，直接报错
		errStr := "rpc服务提供者名称" + serviceName + "不能被导出"
		log.Error(errStr)
		return serviceName, errors.New(errStr)
	}

	service.name = serviceName
	service.method = reflectMethod(service.rType, true)
	if len(service.method) == 0 {
		var errorStr string
		method := reflectMethod(reflect.PtrTo(service.rType), false)
		if len(method) != 0 {
			errorStr = "rpc服务提供者必须传入指针类型"
		} else {
			errorStr = "rpc服务提供者方法必须是可导出的"
		}

		return serviceName, errors.New(errorStr)
	}

	s.serviceMap[serviceName] = service

	return serviceName, nil
}

// 通过指定名称注册函数
func (s *Server) RegisterFuncName(function interface{}, name string, metadata string) error {
	_, err := s.registerFunction(function, name, true)
	if err != nil {
		return err
	}

	s.Plugins.DoRegister(name, function, metadata)

	return nil
}

// 通过反射注册函数
func (s *Server) RegisterFunc(function interface{}, metadata string) error {
	name, err := s.registerFunction(function, "", false)
	if err != nil {
		return err
	}

	s.Plugins.DoRegister(name, function, metadata)
	return nil
}

// 反射注册函数类型
func (s *Server) registerFunction(function interface{}, name string, useName bool) (string, error) {
	s.serviceMapMu.Lock()
	defer s.serviceMapMu.Unlock()

	// 首先断言传入进来的是不是一个反射好的类型，不是的话先进行反射
	f, ok := function.(reflect.Value)
	if !ok {
		f = reflect.ValueOf(function)
	}

	serviceName := formatFuncName(runtime.FuncForPC(reflect.Indirect(f).Pointer()).Name())

	if useName { // 需要取外面传递进来的
		serviceName = name
	}

	if serviceName == "" {
		errorStr := "服务提供者函数方式注册失败，反射不到函数名称" + f.Type().String()
		log.Error(errorStr)
		return serviceName, errors.New(errorStr)
	}

	t := f.Type()
	if t.NumIn() != 3 {
		errorStr := "服务提供者函数方式注册失败，函数入参必须有三个" + f.Type().String()
		log.Error(errorStr)
		return serviceName, errors.New(errorStr)
	}

	if !t.In(0).Implements(typeContext) {
		errorStr := "服务提供者函数方式注册失败，函数第一个参数必须实现context.Context" + f.Type().String()
		log.Error(errorStr)
		return serviceName, errors.New(errorStr)
	}

	requestType := t.In(1)
	if !isExportedOrBuildInType(requestType) {
		errorStr := "服务提供者函数方式注册失败，函数第二个参数必须为可导出的或者内建类型" + f.Type().String()
		log.Error(errorStr)
		return serviceName, errors.New(errorStr)
	}

	responseType := t.In(2)
	if responseType.Kind() != reflect.Ptr {
		errorStr := "服务提供者函数方式注册失败，函数第三个参数必须为指针类型" + f.Type().String()
		log.Error(errorStr)
		return serviceName, errors.New(errorStr)
	}

	if !isExportedOrBuildInType(responseType) {
		errorStr := "服务提供者函数方式注册失败，函数第三个参数必须为可导出的" + f.Type().String()
		log.Error(errorStr)
		return serviceName, errors.New(errorStr)
	}

	if t.NumOut() != 1 || t.Out(0).Implements(errorType) {
		errorStr := "服务提供者函数方式注册失败，函数返回值只能是一个error类型" + f.Type().String()
		log.Error(errorStr)
		return serviceName, errors.New(errorStr)
	}

	service := service{
		name:     serviceName,
		rType:    f.Type(),
		rValue:   f,
		function: make(map[string]*funcType),
	}

	funcType := &funcType{
		rFuc:         f,
		requestType:  requestType,
		responseType: responseType,
	}
	service.function[serviceName] = funcType
	ObjectPool.Init(requestType)
	ObjectPool.Init(responseType)

	return service.name, nil
}

func reflectMethod(rType reflect.Type, logError bool) map[string]*methodType {
	methods := make(map[string]*methodType)
	for i := 0; i < rType.NumMethod(); i++ {
		method := rType.Method(i)

		mType := method.Type
		mName := method.Name
		// PkgPath 不为空这说明这个方法是不可导出的
		if method.PkgPath != "" {
			continue
		}

		if mType.NumIn() != 4 { // 自定义方法的入参必须为3个（这个判断4是因为第0个参数为结构体本身，不包括入参）
			if logError {
				log.DebugF("自定义方法入参个数不对，方法名:%s,个数%d", mName, mType.NumIn())
			}
			continue
		}

		//第一个参数必须为context.Context
		if !mType.In(1).Implements(typeContext) {
			if logError {
				log.Debug("第一个参数必须实现了context.Context接口")
			}
			continue
		}

		requestType := mType.In(2)
		if !isExportedOrBuildInType(requestType) {
			if logError {
				log.DebugF("request 必须为可导出的或者内建类型")
			}
			continue
		}

		responseType := mType.In(3)
		// 第二个参数和第三个参数必须为指针
		if responseType.Kind() != reflect.Ptr {
			if logError {
				log.DebugF("response 必须为指针类型")
			}
			continue
		}

		if !isExportedOrBuildInType(responseType) {
			if logError {
				log.DebugF("response 必须为可导出的类型")
			}
			continue
		}

		if mType.NumOut() != 1 || !mType.Out(0).Implements(errorType) {
			if logError {
				log.DebugF("返回值必须为一个error类型")
			}
			continue
		}

		methodType := &methodType{
			rMethod:      method,
			requestType:  requestType,
			responseType: responseType,
		}

		methods[mName] = methodType

		ObjectPool.Init(requestType)
		ObjectPool.Init(responseType)

	}
	return methods
}

// 判断名称时候能导出（golang首字母必须大写才能导出）
func isExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

// 判断类型时候可以导出或者为内建类型
func isExportedOrBuildInType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return isExported(t.Name()) || t.PkgPath() == ""
}

// 格式话反射出来的函数名称，只取最后部分
func formatFuncName(funcName string) string {
	if funcName != "" {
		if index := strings.LastIndex(funcName, "."); index > 0 {
			funcName = funcName[index+1:]
		}
	}

	return funcName
}
