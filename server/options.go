package server

import (
	"crypto/tls"
	"time"
)

type OptionFunc func(server *Server)

func WithDisableHTTPGateway(flag bool) OptionFunc {
	return func(server *Server) {
		server.disableHTTPGateway = flag
	}
}

func WithDisableJSONRPCGateway(flag bool) OptionFunc {
	return func(server *Server) {
		server.disableJSONRPCGateway = flag
	}
}

func WithReadTimeout(timeout time.Duration) OptionFunc {
	return func(server *Server) {
		server.readTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) OptionFunc {
	return func(server *Server) {
		server.writeTimeout = timeout
	}
}

func WithTlsConfig(tls *tls.Config) OptionFunc {
	return func(server *Server) {
		server.tlsConfig = tls
	}
}
