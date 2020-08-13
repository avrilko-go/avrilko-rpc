package server

import (
	"avrilko-rpc/plugin"
	"avrilko-rpc/service"
)

type IServer interface {
}

type server struct {
	opts     *Options // 服务通用核心配置
	services map[string]service.IService
	plugins  []plugin.Plugin // 插件对象集合
	closing  bool            // 服务是否正在关闭
}

func NewServer(opts ...OptionFuc) IServer {
	s := &server{
		services: make(map[string]service.IService),
		closing:  false,
		opts:     &Options{},
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(s.opts)
		}
	}
	return s
}
