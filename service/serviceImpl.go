package service

import (
	"avrilko-rpc/server"
	"context"
	"log"
	"sync"
)

// rpc服务默认实现
type service struct {
	serviceName string                    // 服务名称
	ctx         context.Context           // 上下文
	cancel      context.CancelFunc        // 取消函数
	handlers    map[string]HandlerService // 处理路由
	closing     bool                      // 单个服务是否正在关闭
}

// 单个服务关闭
func (s *service) Close(wait *sync.WaitGroup) {
	log.Printf("服务(%s)正在关闭中...", s.serviceName)
	s.closing = true
	if s.cancel != nil {
		s.cancel()
	}
	s.closing = false
	log.Printf("服务(%s)已经成功关闭...", s.serviceName)
	wait.Add(1)
}

func (s *service) Serve(options *server.Options) {
	panic("implement me")
}

// 直接注册服务的方法 相当于路由注册
func (s *service) Register(handlerName string, handlerService HandlerService) {
	if s.handlers == nil {
		s.handlers = make(map[string]HandlerService)
	}

	s.handlers[handlerName] = handlerService
}
