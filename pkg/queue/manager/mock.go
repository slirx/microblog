package manager

import (
	"context"
)

type Mock struct {
	SendFn  func(ctx context.Context, routingKey string, msg interface{}) error
	CloseFn func() error
}

func (m Mock) Send(ctx context.Context, routingKey string, msg interface{}) error {
	return m.SendFn(ctx, routingKey, msg)
}

func (m Mock) Close() error {
	return m.CloseFn()
}
