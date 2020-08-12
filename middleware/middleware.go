package middleware

import "context"

// 服务端中间件
type ServerMiddleware func(ctx context.Context, request interface{}, handler ServerCoreHandler) (interface{}, error)

// 服务端中间件最终执行的函数(洋葱圈最中间部分)
type ServerCoreHandler func(ctx context.Context, request interface{}) (interface{}, error)

// 服务端开始执行中间件
func SeverBeginMiddleware(ctx context.Context, request interface{}, middlewares []ServerMiddleware, handler ServerCoreHandler) (interface{}, error) {
	if len(middlewares) == 0 {
		return handler(ctx, request)
	}

	return middlewares[0](ctx, request, serverBegin(0, middlewares, handler))
}

// 包装生成ServerCoreHandler 中间件最核心逻辑
func serverBegin(index int, middlewares []ServerMiddleware, handler ServerCoreHandler) ServerCoreHandler {
	if index >= len(middlewares)-1 { // 中间件已经执行完成，直接返回服务端核心函数
		return handler
	}

	// 递归执行中间件
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return middlewares[index+1](ctx, request, serverBegin(index+1, middlewares, handler))
	}
}
