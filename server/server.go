package server

import (
	"avrilko-rpc/log"
	"avrilko-rpc/protocol"
	"avrilko-rpc/share"
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var ErrServerClosed = errors.New("主服务已经关闭")

const (
	ReadBuffSize = 1024 // 读取消息时候缓冲区大小
)

type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return c.name
}

var (
	RemoteConnContextKey   = &contextKey{"remote_conn"}
	StartRequestContextKey = &contextKey{"start-parse-request"}
)

// 核心服务类
type Server struct {
	ln           net.Listener  // 全局唯一的监听（可以多路复用实现不同协议的转发）
	readTimeout  time.Duration // 读超时
	writeTimeout time.Duration // 写超时

	gatewayHttpServer     *http.Server // 当启用http网关时候被挂载
	disableHTTPGateway    bool         // 是否禁用http网关服务（开启时候方便测试和调试rpc服务）
	disableJSONRPCGateway bool         // 是否禁用json rpc网关服务 

	serviceMapMu sync.RWMutex        // 服务提供者map读写锁
	serviceMap   map[string]*service // 服务提供者集合map

	connMu     sync.RWMutex          // 各个活跃连接读写锁
	activeConn map[net.Conn]struct{} // 每个活跃的连接，map结构防止重复
	doneChan   chan struct{}         // 服务结束chan

	inShutdown int32             //服务是否关闭 1为关闭 0为正在运行
	onShutdown []func(s *Server) // 服务结束后执行的钩子函数

	tlsConfig *tls.Config // tls证书配置

	Plugins PluginContainer // 插件容器（设计核心）

	AuthFunc func(ctx context.Context, request *protocol.Message, token string) error // 认证函数

	handlerMsgNum int32 // 正在处理的消息数量
}

// 初始化服务
func NewServer(opts ...OptionFunc) *Server {
	server := &Server{
		serviceMapMu: sync.RWMutex{},
		serviceMap:   make(map[string]*service),
		activeConn:   make(map[net.Conn]struct{}),
		doneChan:     make(chan struct{}),
		Plugins:      &pluginContainer{},
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(server)
		}
	}

	return server
}

// 开启服务
func (s *Server) Serve(network, address string) error {
	var ln net.Listener
	var err error
	ln, err = s.makeListener(network, address)
	if err != nil {
		return err
	}
	return s.ServeListener(network, ln)
}

// 开启服务
func (s *Server) ServeListener(network string, ln net.Listener) error {
	// 开启信号量监听
	s.startShutdownServe()
	// 开启网关
	s.startGateway(network, ln)

	return s.serveListener(ln)
}

// 循环监听conn 并发给severConn处理
func (s *Server) serveListener(ln net.Listener) error {
	// 定义临时错误的延迟时间
	var tempDelay time.Duration

	s.connMu.Lock()
	s.ln = ln
	s.connMu.Unlock()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-s.doneChan:
				return ErrServerClosed
			default:
			}
			// 如果错误断言为网络错误，且是一个临时的（比如当时网络环境差，dns服务器不稳定引起的），稍后可能会自动恢复的
			// 不能直接返回错误，应该等待一段时间才返回错误
			// 等待的时间为上一次的两倍最大等待1s
			// 参考官方http包实现的
			if ne, ok := conn.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay = tempDelay * 2
				}

				if tempDelay > time.Second { // 大于1秒直接返回错误
					return err
				}
				time.Sleep(tempDelay)
				log.ErrorF("rpc服务接受conn异常，正在重试, 原因%v, sleep %d", err, tempDelay)
				continue
			}

			if strings.Contains(err.Error(), "listener closed") { // 服务关闭
				return ErrServerClosed
			}
			return err
		}

		// 成功请求延迟时间置为0
		tempDelay = 0

		if tc, ok := conn.(*net.TCPConn); ok { // tcp请求需要设置keepAlive保证链接的稳定性能
			tc.SetKeepAlive(true)
			tc.SetKeepAlivePeriod(time.Minute * 5) // 5分钟没有响应报错
			tc.SetLinger(10)                       // 关闭连接的行为 设置数据在断开时候也能在后台发送
		}

		conn, ok := s.Plugins.DoPostConnAccept(conn)
		if !ok { // 不允许链接则关闭（可能是限流没通过，验证没通过，业务方面的自己用插件扩展...）
			s.closeChannel(conn)
		}
		s.connMu.Lock()
		s.activeConn[conn] = struct{}{}
		s.connMu.Unlock()

		go s.serveConn(conn)
	}
}

// 开始处理消息
func (s *Server) serveConn(conn net.Conn) {
	// 单个conn协程中没有权限影响主进程panic，所有panic会这一层处理
	defer func() {
		if err := recover(); err != nil { // 发生panic
			buf := make([]byte, 65536)
			size := runtime.Stack(buf, false)
			if size > 65536 {
				size = 65536
			}
			buf = buf[:size]
			log.ErrorF("conn 发生panic,原因%s, 客户端地址%s, 堆栈信息 %s", err, conn.RemoteAddr(), buf)
		}
		s.connMu.Lock()
		delete(s.activeConn, conn)
		s.connMu.Unlock()
		s.Plugins.DoPostConnClose(conn)
	}()

	// 判断此时服务是否已经关闭
	if s.isShutdown() {
		s.closeChannel(conn)
		return
	}

	now := time.Now()
	// tls连接需要先握手
	if tlsL, ok := conn.(*tls.Conn); ok {
		if s.readTimeout != 0 {
			tlsL.SetReadDeadline(now.Add(s.readTimeout))
		}
		if s.writeTimeout != 0 {
			tlsL.SetWriteDeadline(now.Add(s.writeTimeout))
		}
		if err := tlsL.Handshake(); err != nil {
			log.ErrorF("tls尝试握手失败，原因：%s，addr:", err, tlsL.RemoteAddr())
			return
		}
	}
	// 初始化读取缓冲区
	rBuff := bufio.NewReaderSize(conn, ReadBuffSize)
	for {
		// 判断此时服务是否已经关闭
		if s.isShutdown() {
			s.closeChannel(conn)
			return
		}

		if s.readTimeout != 0 { // 设置读取的超时时间
			conn.SetReadDeadline(now.Add(s.readTimeout))
		}

		ctx := share.WithValue(context.Background(), RemoteConnContextKey, conn)
		request, err := s.readRequest(ctx, rBuff)
		if err != nil {
			if err == io.EOF {
				log.InfoF("客户端已经关闭链接c%s", conn.RemoteAddr())
			} else if strings.Contains(err.Error(), "use of closed network connection") {
				log.InfoF("连接已经被关闭%s", conn.RemoteAddr())
			} else {
				log.WarnF("rpc 读取数据失败，错误原因%v", err)
			}
			return
		}

		// 要开始写入了
		if s.writeTimeout != 0 {
			conn.SetWriteDeadline(now.Add(s.writeTimeout))
		}

		// 将开始时间写上下文中
		ctx = share.WithLocalValue(ctx, StartRequestContextKey, time.Now().UnixNano())
		if !request.IsHeartbeat() { // auth鉴权
			err := s.auth(ctx, request)
			if err != nil { // 鉴权失败
				if !request.IsOneway() { // 需要回复客户端鉴权失败
					response := request.Clone()                // 复制一个请求出来
					response.SetMessageType(protocol.Response) // 设置为response消息
					handleError(response, err)
					data := response.EncodeSlicePointer()
					_, err = conn.Write(*data)
					protocol.PutData(data)
					s.Plugins.DoPostWriteResponse(ctx, request, response, err)
					protocol.FreeMsg(response)
				} else { // 不需要回复
					s.Plugins.DoPreWriteResponse(ctx, request, nil)
				}
				protocol.FreeMsg(request)
				log.InfoF("连接鉴权失败，%s,错误原因%v", conn.RemoteAddr(), err)
				return
			}
		}

		// 下面需要处理消息了噢
		go func() {
			// 正在处理的消息数量+1
			atomic.AddInt32(&s.handlerMsgNum, 1)
			// 正在处理消息的数量-1
			defer atomic.AddInt32(&s.handlerMsgNum, - 1)

			if request.IsHeartbeat() { // 如果是客户端心跳
				request.SetMessageType(protocol.Response)
				data := request.EncodeSlicePointer()
				conn.Write(*data)
				protocol.PutData(data)
				return
			}

			// 不是心跳初始化返给客户端的meta
			responseMetadata := make(map[string]string)
			// 先将服务端的metadata方法进去
			ctx = share.WithLocalValue(ctx, share.ReqMetaDataKey, request.Metadata)
			// 再将客户端的metadata放进去
			ctx = share.WithLocalValue(ctx, share.ResMetaDataKey, responseMetadata)

			s.Plugins.DoPreHandleRequest(ctx, request) // 开始处理请求了

			response, err := s.handleRequest(ctx, request)
			if err != nil {
				log.WarnF("处理请求错误: %v", err)
				return
			}
			s.Plugins.DoPreWriteResponse(ctx, request, response)
			if !request.IsOneway() { // 需要回复客户端
				// 从ctx中拿出meta信息
				responseMetadataCtx := ctx.Value(share.ResMetaDataKey).(map[string]string)
				if len(responseMetadata) > 0 {
					meta := response.Metadata
					if meta == nil {
						meta = responseMetadataCtx
					} else {
						for k, v := range responseMetadataCtx {
							if meta[k] == "" {
								meta[k] = v
							}
						}
					}
				}

				if len(response.Payload) > 1024 && request.CompressType() != protocol.None {
					response.SetCompressType(request.CompressType())
				}
				data := response.EncodeSlicePointer()
				conn.Write(*data)
				protocol.PutData(data)
			}

			s.Plugins.DoPostWriteResponse(ctx, request, response, err)

			protocol.FreeMsg(response)
			protocol.FreeMsg(request)
		}()
	}
}

// 处理单个请求
func (s *Server) handleRequest(ctx context.Context, request *protocol.Message) (*protocol.Message, error) {
	var err error
	serviceName := request.ServicePath
	methodName := request.ServiceMethod

	response := request.Clone()
	response.SetMessageType(protocol.Response)

	s.serviceMapMu.RLock()
	service, ok := s.serviceMap[serviceName]
	s.serviceMapMu.RUnlock()
	if !ok { // 都没注册直接返回错误
		err = errors.New(fmt.Sprintf("不能找到服务发现者为%s的服务", serviceName))
		return handleError(response, err)
	}

	methodType, ok := service.method[methodName]
	if !ok { // 看看是否注册了函数的调用
		if _, ok := service.function[methodName]; ok {
			return s.handleRequestForFunction(ctx, request)
		}
		err = errors.New(fmt.Sprintf("不能找到服务提供者%s下方法名为%s的方法", serviceName, methodName))
		return handleError(response, err)
	}

	requestType := ObjectPool.Get(methodType.requestType)
	defer ObjectPool.Put(methodType.requestType, requestType)
	codec := share.Codecs[request.SerializeType()]
	if codec == nil {
		err = fmt.Errorf("不能找到对应的的序列化方式：%T", request.SerializeType())
		return handleError(response, err)
	}

	err = codec.Decode(request.Payload, requestType)
	if err != nil {
		return handleError(response, err)
	}

	responseType := ObjectPool.Get(methodType.responseType)
	defer ObjectPool.Put(methodType.responseType, responseType)

	responseType, err = s.Plugins.DoPreCall(ctx, serviceName, methodName, requestType)
	if err != nil {
		return handleError(response, err)
	}
	if methodType.requestType.Kind() != reflect.Ptr { // 不是指针
		err = service.call(ctx, methodType, reflect.ValueOf(requestType).Elem(), reflect.ValueOf(responseType))
	} else {
		err = service.call(ctx, methodType, reflect.ValueOf(requestType), reflect.ValueOf(responseType))
	}
	if err != nil {
		return handleError(response, err)
	}
	responseType, err = s.Plugins.DoPostCall(ctx, serviceName, methodName, requestType, responseType)
	if err != nil {
		return handleError(response, err)
	}

	if !request.IsOneway() {
		data, err := codec.Encode(responseType)
		if err != nil {
			return handleError(response, err)
		}
		response.Payload = data
	}
	return response, nil
}

func (s *Server) handleRequestForFunction(ctx context.Context, request *protocol.Message) (*protocol.Message, error) {

}

func (s *Server) auth(ctx context.Context, request *protocol.Message) error {
	if s.AuthFunc == nil {
		return nil
	}
	token := request.Metadata[share.AuthKey]
	return s.AuthFunc(ctx, request, token)
}

// 暴力关闭服务（生产环境不建议使用，建议使用Shutdown）
func (s *Server) Close() error {
	s.serviceMapMu.Lock()
	defer s.serviceMapMu.Unlock()
	var err error
	if s.ln != nil {
		err = s.ln.Close()
	}

	for conn, _ := range s.activeConn {
		err = conn.Close()
		delete(s.activeConn, conn)
		s.Plugins.DoPostConnClose(conn)
	}

	return err
}

// 优雅的关闭服务，
// 先关闭tcp监听，使得不再有conn连接进来
// 关闭每个conn的读端，使其不再收客户端数据
// 循环等待服务正在处理的消息数量变为0 (使得所有正在处理的消息都能处理完成)
// 关闭网关服务
// 依次关闭conn的读和写
func (s *Server) Shutdown(ctx context.Context) error {
	var err error
	if atomic.CompareAndSwapInt32(&s.inShutdown, 0, 1) { // 保证结束进程只执行一次
		log.Info("服务开始关闭...")
		// 先关闭tcp链接的读端（写端要等所有请求都结束后才能关闭）
		s.connMu.Lock()
		if s.ln != nil {
			s.ln.Close() // 关闭监听
		}
		for conn, _ := range s.activeConn {
			if lConn, ok := conn.(*net.TCPConn); ok {
				lConn.CloseRead()
			}
		}
		s.connMu.Unlock()

		ticker := time.NewTicker(time.Second) // 监听间隔
		defer ticker.Stop()

	outer:
		for {
			if s.checkMsgHandlerFinish() {
				break
			}
			select {
			case <-ctx.Done():
				break outer

			case <-ticker.C:

			}
		}

		if s.gatewayHttpServer != nil {
			if err = s.closeHTTP1APIGateway(ctx); err != nil {
				log.WarnF("关闭http网关时出错：%v", err)
			} else {
				log.Info("http网关服务已经关闭")
			}
		}

		s.connMu.Lock()
		for conn, _ := range s.activeConn {
			conn.Close()
			delete(s.activeConn, conn)
			s.Plugins.DoPostConnClose(conn)
		}
		s.closeDoneChanLocked()
		s.connMu.Unlock()

	}

	return err
}

// 检测服务处理的消息是否处理完成
func (s *Server) checkMsgHandlerFinish() bool {
	size := atomic.LoadInt32(&s.handlerMsgNum)
	log.InfoF("还需要处理%d条消息", size)
	return size == 0
}

// 关闭结束通道（如果别的协程已经关闭，则直接返回）
func (s *Server) closeDoneChanLocked() {
	select {
	case <-s.doneChan:
		return
	default:
		close(s.doneChan)
	}
}

// 监听结束服务事件(terminated)
func (s *Server) startShutdownServe() {
	go func(s *Server) {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		sg := <-c
		if sg.String() == "terminated" {
			if s.onShutdown != nil && len(s.onShutdown) > 0 {
				for _, shutdown := range s.onShutdown {
					shutdown(s)
				}
			}
			err := s.Shutdown(context.Background())
			if err != nil {
				log.Error(err.Error())
			}
		}
	}(s)
}

// 关闭链接
func (s *Server) closeChannel(conn net.Conn) {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	delete(s.activeConn, conn)
	conn.Close()
}

// 判断服务是否已经关闭了
func (s *Server) isShutdown() bool {
	return atomic.LoadInt32(&s.inShutdown) == 1
}

func (s *Server) readRequest(ctx context.Context, rBuff io.Reader) (*protocol.Message, error) {
	var err error
	err = s.Plugins.DoPreReadRequest(ctx)
	if err != nil {
		return nil, err
	}
	request := protocol.GetPooledMsg()
	err = s.Plugins.DoPreReadRequest(ctx)
	if err != nil {
		return nil, err
	}
	// 开始解码
	err = request.Decode(rBuff)
	if err == io.EOF { // io.EOF代表读完了
		return request, err
	}

	pErr := s.Plugins.DoPostReadRequest(ctx, request, err)
	if err == nil { // 看看插件里面的调用会报什么错误
		err = pErr
	}

	return request, err
}

// 处理错误
func handleError(response *protocol.Message, err error) (*protocol.Message, error) {
	response.SetMessageStatusType(protocol.Error)
	if response.Metadata == nil {
		response.Metadata = make(map[string]string, 10)
	}

	response.Metadata[protocol.ServiceError] = err.Error()

	return response, err
}
