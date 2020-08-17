package server

import (
	"avrilko-rpc/protocol"
	"context"
	"github.com/soheilhy/cmux"
	"io"
	"net"
)

const (
	JsonRpcHeader = "X-JSONRPC-2.0"
)

// 关闭http网关
func (s *Server) closeHTTP1APIGateway(ctx context.Context) error {
	s.connMu.Lock()
	defer s.connMu.Unlock()

	if s.gatewayHttpServer != nil {
		return s.gatewayHttpServer.Shutdown(ctx)
	}
	return nil
}

// 根据配置文件开启多路复用的网关服务
func (s *Server) startGateway(network string, ln net.Listener) net.Listener {
	if network != "tcp" && network != "tcp4" && network != "tcp6" {
		// 不是tcp协议不能使用多路复用直接返回
		return ln
	}
	mu := cmux.New(ln)

	l := mu.Match(matchIsAvrilkoRpc())
	if !s.disableJSONRPCGateway { // 没有禁用json-rpc网关的话开启网关服务
		jsonRpcL := mu.Match(cmux.HTTP1HeaderField(JsonRpcHeader, "true"))
		go s.startJSONRPC2(jsonRpcL)
	}
	if !s.disableHTTPGateway { // 没有禁用http网关
		httpL := mu.Match(cmux.HTTP1Fast())
		go s.startHTTP1APIGateway(httpL)
	}
	return l
}

// 根据协议来判断是不是自定义的tcp协议
func matchIsAvrilkoRpc() cmux.Matcher {
	return func(reader io.Reader) bool {
		oneByte := make([]byte, 1)
		reader.Read(oneByte)
		if oneByte[0] == protocol.MagicNumber() {
			return true
		} else {
			return false
		}
	}
}
