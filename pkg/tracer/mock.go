package tracer

import (
	"context"
)

var _ Tracer = (*Mock)(nil)

// Mock represents Tracer interface mock.
type Mock struct {
	RequestIDFn func(ctx context.Context) string
}

func (m Mock) RequestID(ctx context.Context) string {
	return m.RequestIDFn(ctx)
}
