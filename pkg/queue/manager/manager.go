package manager

import (
	"context"
)

type Manager interface {
	Send(ctx context.Context, routingKey string, msg interface{}) error
	Close() error
	//EmitEvent(ctx context.Context, exchange string, msg interface{}) error
}
