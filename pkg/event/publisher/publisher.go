package publisher

import (
	"context"
)

type Publisher interface {
	Emit(ctx context.Context, exchange string, msg interface{}) error
}
