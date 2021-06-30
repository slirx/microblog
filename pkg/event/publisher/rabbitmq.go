package publisher

import (
	"context"
)

type rabbitmqPublisher struct {
}

func (r rabbitmqPublisher) Emit(ctx context.Context, exchange string, msg interface{}) error {
	// todo implement
	return nil
}
