package media

import (
	"context"
)

type Mock struct {
	ImagesFn func(ctx context.Context, serviceName string, ids []int) (map[int]string, error)
}

func (m Mock) Images(ctx context.Context, serviceName string, ids []int) (map[int]string, error) {
	return m.ImagesFn(ctx, serviceName, ids)
}
