package middleware

import (
	"context"
	"github.com/stretchr/testify/assert"
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
	assert.Nil(t, err)
	s, ok := result.(string)
	assert.True(t, ok)
	assert.Equal(t, s, "handlerfunc1func2")

}
