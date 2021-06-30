package user

import (
	"context"

	"gitlab.com/slirx/newproj/pkg/queue"
)

var _ Service = (*serviceMock)(nil)
var _ Repository = (*repositoryMock)(nil)

type serviceMock struct {
	UpdateFn            func(ctx context.Context, request UpdateRequest) (string, error)
	GetFn               func(ctx context.Context, login string) (*GetResponse, error)
	MeFn                func(ctx context.Context) (*MeResponse, error)
	FollowFn            func(ctx context.Context, request FollowRequest) error
	UnfollowFn          func(ctx context.Context, request UnfollowRequest) error
	FollowersFn         func(ctx context.Context, request FollowersRequest) (*FollowersResponse, error)
	FollowingFn         func(ctx context.Context, request FollowingRequest) (*FollowingResponse, error)
	InternalFollowersFn func(ctx context.Context, uid int) ([]int, error)
	InternalUsersFn     func(ctx context.Context, userIDs []int) ([]User, error)
}

func (s serviceMock) Get(ctx context.Context, login string) (*GetResponse, error) {
	return s.GetFn(ctx, login)
}

func (s serviceMock) Me(ctx context.Context) (*MeResponse, error) {
	return s.MeFn(ctx)
}

func (s serviceMock) Follow(ctx context.Context, request FollowRequest) error {
	return s.FollowFn(ctx, request)
}

func (s serviceMock) Unfollow(ctx context.Context, request UnfollowRequest) error {
	return s.UnfollowFn(ctx, request)
}

func (s serviceMock) Followers(ctx context.Context, request FollowersRequest) (*FollowersResponse, error) {
	return s.FollowersFn(ctx, request)
}

func (s serviceMock) Following(ctx context.Context, request FollowingRequest) (*FollowingResponse, error) {
	return s.FollowingFn(ctx, request)
}

func (s serviceMock) InternalFollowers(ctx context.Context, uid int) ([]int, error) {
	return s.InternalFollowersFn(ctx, uid)
}

func (s serviceMock) InternalUsers(ctx context.Context, userIDs []int) ([]User, error) {
	return s.InternalUsersFn(ctx, userIDs)
}

type repositoryMock struct {
	CreateFn       func(ctx context.Context, request queue.UserCreate) (int, error)
	UpdateFn       func(ctx context.Context, uid int, request UpdateRequest) error
	GetFn          func(ctx context.Context, login string, uid int) (*GetResponse, error)
	MeFn           func(ctx context.Context, uid int) (*MeResponse, error)
	FollowFn       func(ctx context.Context, uid int, request FollowRequest) error
	UnfollowFn     func(ctx context.Context, uid int, request UnfollowRequest) error
	FollowersFn    func(ctx context.Context, request FollowersRequest, perPage uint8) (*FollowersResponse, error)
	FollowersIDsFn func(ctx context.Context, uid int) ([]int, error)
	FollowingFn    func(ctx context.Context, request FollowingRequest, perPage uint8) (*FollowingResponse, error)
	UsersFn        func(ctx context.Context, userIDs []int) ([]User, error)
}

func (s serviceMock) Update(ctx context.Context, request UpdateRequest) (string, error) {
	return s.UpdateFn(ctx, request)
}

func (r repositoryMock) Create(ctx context.Context, request queue.UserCreate) (int, error) {
	return r.CreateFn(ctx, request)
}

func (r repositoryMock) Update(ctx context.Context, uid int, request UpdateRequest) error {
	return r.UpdateFn(ctx, uid, request)
}

func (r repositoryMock) Get(ctx context.Context, login string, uid int) (*GetResponse, error) {
	return r.GetFn(ctx, login, uid)
}

func (r repositoryMock) Me(ctx context.Context, uid int) (*MeResponse, error) {
	return r.MeFn(ctx, uid)
}

func (r repositoryMock) Follow(ctx context.Context, uid int, request FollowRequest) error {
	return r.FollowFn(ctx, uid, request)
}

func (r repositoryMock) Unfollow(ctx context.Context, uid int, request UnfollowRequest) error {
	return r.UnfollowFn(ctx, uid, request)
}

func (r repositoryMock) Followers(ctx context.Context, request FollowersRequest, perPage uint8) (*FollowersResponse, error) {
	return r.FollowersFn(ctx, request, perPage)
}

func (r repositoryMock) FollowersIDs(ctx context.Context, uid int) ([]int, error) {
	return r.FollowersIDsFn(ctx, uid)
}

func (r repositoryMock) Following(ctx context.Context, request FollowingRequest, perPage uint8) (*FollowingResponse, error) {
	return r.FollowingFn(ctx, request, perPage)
}

func (r repositoryMock) Users(ctx context.Context, userIDs []int) ([]User, error) {
	return r.UsersFn(ctx, userIDs)
}
