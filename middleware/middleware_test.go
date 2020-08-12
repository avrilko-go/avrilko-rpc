package middleware

import (
	"context"
	"testing"
)

func TestMiddleware(t *testing.T) {
	var func1 ServerMiddleware = func(ctx context.Context, request interface{}, handler ServerCoreHandler) (interface{}, error) {
		s := request.(string)
		s = s + "func1"
		return handler(ctx, s)
	}

	var func2 ServerMiddleware = func(ctx context.Context, request interface{}, handler ServerCoreHandler) (interface{}, error) {
		s := request.(string)
		s = s + "func2"
		return handler(ctx, s)
	}

	var handler ServerCoreHandler = func(ctx context.Context, request interface{}) (interface{}, error) {
		return request, nil
	}

	middlewares := []ServerMiddleware{func1, func2}

	result, err := SeverBeginMiddleware(context.TODO(), "handler", middlewares, handler)
	if err != nil {
		t.Error(err)
	}

	s, ok := result.(string)
	if !ok {
		t.Error("断言失败")
	}

	if s != "handlerfunc1func2" {
		t.Error("中间件失败")
	}
}
