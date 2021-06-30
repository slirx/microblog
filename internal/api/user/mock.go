package user

import (
	"context"

	"gitlab.com/slirx/newproj/internal/graphql/graph/model"
)

type Mock struct {
	FollowersFn func(ctx context.Context, uid int) ([]int, error)
	UsersFn     func(ctx context.Context, userIDs []int) ([]model.User, error)
}

func (m Mock) Followers(ctx context.Context, uid int) ([]int, error) {
	return m.FollowersFn(ctx, uid)
}

func (m Mock) Users(ctx context.Context, userIDs []int) ([]model.User, error) {
	return m.UsersFn(ctx, userIDs)
}
