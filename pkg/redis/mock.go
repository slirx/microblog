package redis

import (
	"context"
	"time"
)

type Mock struct {
	SetFn         func(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetFn         func(ctx context.Context, key string) (string, error)
	SetIntSliceFn func(ctx context.Context, key string, value []int, expiration time.Duration) error
	GetIntSliceFn func(ctx context.Context, key string) ([]int, error)
	HIncrByFn     func(ctx context.Context, key string, field string, incr int64) (int64, error)
	CloseFn       func() error
}

func (m Mock) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return m.SetFn(ctx, key, value, expiration)
}

func (m Mock) Get(ctx context.Context, key string) (string, error) {
	return m.GetFn(ctx, key)
}

func (m Mock) SetIntSlice(ctx context.Context, key string, value []int, expiration time.Duration) error {
	return m.SetIntSliceFn(ctx, key, value, expiration)
}

func (m Mock) GetIntSlice(ctx context.Context, key string) ([]int, error) {
	return m.GetIntSliceFn(ctx, key)
}

func (m Mock) HIncrBy(ctx context.Context, key string, field string, incr int64) (int64, error) {
	return m.HIncrByFn(ctx, key, field, incr)
}

func (m Mock) Close() error {
	return m.CloseFn()
}
