package post

import (
	"context"

	"gitlab.com/slirx/newproj/pkg/queue"
)

var _ Repository = (*repositoryMock)(nil)
var _ Service = (*serviceMock)(nil)

type serviceMock struct {
	ListFn   func(ctx context.Context, request ListRequest) (*ListResponse, error)
	CreateFn func(ctx context.Context, request CreateRequest) (*CreateResponse, error)
	FeedFn   func(ctx context.Context, userID int, request FeedRequest) (*FeedResponse, error)
	SearchFn func(ctx context.Context, request SearchRequest) (*SearchResponse, error)
}

type repositoryMock struct {
	ListFn     func(ctx context.Context, request ListRequest, perPage int) (*ListResponse, error)
	CreateFn   func(ctx context.Context, uid int, users []int, request CreateRequest) (*CreateResponse, error)
	FeedFn     func(ctx context.Context, uid int, request FeedRequest, perPage int) (*FeedResponse, error)
	UnfollowFn func(ctx context.Context, task queue.PostUnfollow) error
	FollowFn   func(ctx context.Context, task queue.PostFollow) error
	SearchFn   func(ctx context.Context, request SearchRequest, perPage int) ([]int, error)
	PostsFn    func(ctx context.Context, postsIDs []int) ([]Post, error)
}

func (s serviceMock) List(ctx context.Context, request ListRequest) (*ListResponse, error) {
	return s.ListFn(ctx, request)
}

func (s serviceMock) Create(ctx context.Context, request CreateRequest) (*CreateResponse, error) {
	return s.CreateFn(ctx, request)
}

func (s serviceMock) Feed(ctx context.Context, userID int, request FeedRequest) (*FeedResponse, error) {
	return s.FeedFn(ctx, userID, request)
}

func (s serviceMock) Search(ctx context.Context, request SearchRequest) (*SearchResponse, error) {
	return s.SearchFn(ctx, request)
}

func (r repositoryMock) List(ctx context.Context, request ListRequest, perPage int) (*ListResponse, error) {
	return r.ListFn(ctx, request, perPage)
}

func (r repositoryMock) Create(ctx context.Context, uid int, users []int, request CreateRequest) (*CreateResponse, error) {
	return r.CreateFn(ctx, uid, users, request)
}

func (r repositoryMock) Feed(ctx context.Context, uid int, request FeedRequest, perPage int) (*FeedResponse, error) {
	return r.FeedFn(ctx, uid, request, perPage)
}

func (r repositoryMock) Unfollow(ctx context.Context, task queue.PostUnfollow) error {
	return r.Unfollow(ctx, task)
}

func (r repositoryMock) Follow(ctx context.Context, task queue.PostFollow) error {
	return r.FollowFn(ctx, task)
}

func (r repositoryMock) Search(ctx context.Context, request SearchRequest, perPage int) ([]int, error) {
	return r.SearchFn(ctx, request, perPage)
}

func (r repositoryMock) Posts(ctx context.Context, postsIDs []int) ([]Post, error) {
	return r.PostsFn(ctx, postsIDs)
}
