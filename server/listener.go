package server

import (
	"crypto/tls"
	"errors"
	"net"
)

var makeListeners = make(map[string]MakeListener)

func init() {
	makeListeners["tcp"] = tcpMakeListener("tcp")
	makeListeners["tcp4"] = tcpMakeListener("tcp4")
	makeListeners["tcp6"] = tcpMakeListener("tcp6")
}

type MakeListener func(s *Server, address string) (ln net.Listener, err error)

// 注册监听生成者
func RegisterListener(network string, listener MakeListener) {
	makeListeners[network] = listener
}

// 生成Listener
func (s *Server) makeListener(network, address string) (net.Listener, error) {
	mu, ok := makeListeners[network]
	if !ok {
		return nil, errors.New("暂不支持该网络类型生成tcp listener" + network)
	}
	return mu(s, address)
}

// 生成监听tcp服务
func tcpMakeListener(network string) MakeListener {
	return func(s *Server, address string) (ln net.Listener, err error) {
		if s.tlsConfig != nil {
			ln, err = tls.Listen(network, address, s.tlsConfig)
		} else {
			ln, err = net.Listen(network, address)
		}
		return ln, err
	}
}
