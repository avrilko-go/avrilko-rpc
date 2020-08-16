package server

import "context"

// 关闭http网关
func (s *Server) closeHTTP1APIGateway(ctx context.Context) error {
	s.connMu.Lock()
	defer s.connMu.Unlock()

	if s.gatewayHttpServer != nil {
		return s.gatewayHttpServer.Shutdown(ctx)
	}
	return nil
}
