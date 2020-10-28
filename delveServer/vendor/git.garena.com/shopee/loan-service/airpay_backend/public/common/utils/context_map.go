package utils

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"sync"
)

var (
	contextMapLock sync.RWMutex
)

type contextMapKey struct {
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = addMapToCtx(ctx)
		return handler(ctx, req)
	}
}

func addMapToCtx(ctx context.Context) context.Context {
	v, ok := ctx.Value(contextMapKey{}).(map[string]interface{})
	if !ok {
		v = make(map[string]interface{})
	}
	ctx = context.WithValue(ctx, contextMapKey{}, v)
	return ctx
}

func Set(ctx context.Context, key string, value interface{}) error {
	v, ok := ctx.Value(contextMapKey{}).(map[string]interface{})
	if !ok {
		return errors.New("context map not found in the context")
	}
	contextMapLock.Lock()
	v[key] = value
	contextMapLock.Unlock()
	return nil
}

func Get(ctx context.Context, key string) (interface{}, bool, error) {
	m, ok := ctx.Value(contextMapKey{}).(map[string]interface{})
	if !ok {
		err := errors.New("context map not found in the context")
		return nil, false, err
	}
	contextMapLock.RLock()
	v, ok := m[key]
	contextMapLock.RUnlock()
	return v, ok, nil
}
