package post

import (
	"context"
	"errors"
	"testing"
	"time"

	"gitlab.com/slirx/newproj/internal/api/user"
	"gitlab.com/slirx/newproj/internal/graphql/graph/model"
	"gitlab.com/slirx/newproj/pkg/http/jwtmiddleware"
	"gitlab.com/slirx/newproj/pkg/redis"
)

func TestServiceListSuccess(t *testing.T) {
	wantTotal := 3
	wantUserID := 1
	wantLogin := "anon"
	wantName := "Anon"
	wantPhotoURL := "https://test.com/test.png"

	u := PostsUser{
		ID: wantUserID,
	}

	posts := []Post{
		{
			ID:        1,
			Text:      "post 1",
			CreatedAt: 0,
			User:      u,
		},
		{
			ID:        2,
			Text:      "post 2",
			CreatedAt: 0,
			User:      u,
		},
		{
			ID:        3,
			Text:      "post 3",
			CreatedAt: 0,
			User:      u,
		},
	}

	rMock := repositoryMock{}
	rMock.ListFn = func(ctx context.Context, request ListRequest, perPage int) (*ListResponse, error) {
		return &ListResponse{
			Total: wantTotal,
			Posts: posts,
		}, nil
	}

	internalUserAPIMock := user.Mock{}
	internalUserAPIMock.UsersFn = func(ctx context.Context, userIDs []int) ([]model.User, error) {
		return []model.User{
			{
				ID:       wantUserID,
				Login:    wantLogin,
				Name:     wantName,
				PhotoURL: wantPhotoURL,
			},
		}, nil
	}

	redisClientMock := redis.Mock{}
	s := NewService(rMock, internalUserAPIMock, redisClientMock)

	request := ListRequest{}

	response, err := s.List(context.Background(), request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err)
	}

	if response.Total != wantTotal {
		t.Fatalf("got: %d, want: %d", response.Total, wantTotal)
	}

	if wantUserID != response.Posts[0].User.ID {
		t.Fatalf("got: %d, want: %d", response.Posts[0].User.ID, wantUserID)
	}

	if wantLogin != response.Posts[0].User.Login {
		t.Fatalf("got: %s, want: %s", response.Posts[0].User.Login, wantLogin)
	}

	if wantName != response.Posts[0].User.Name {
		t.Fatalf("got: %s, want: %s", response.Posts[0].User.Name, wantName)
	}

	if wantPhotoURL != response.Posts[0].User.PhotoURL {
		t.Fatalf("got: %s, want: %s", response.Posts[0].User.PhotoURL, wantPhotoURL)
	}

	if wantTotal != len(response.Posts) {
		t.Fatalf("got: %d, want: %d", len(response.Posts), wantTotal)
	}
}

func TestServiceListRepositoryError(t *testing.T) {
	wantErr := "list error"

	u := PostsUser{
		ID: 1,
	}

	posts := []Post{
		{
			ID:        1,
			Text:      "post 1",
			CreatedAt: 0,
			User:      u,
		},
	}

	rMock := repositoryMock{}
	rMock.ListFn = func(ctx context.Context, request ListRequest, perPage int) (*ListResponse, error) {
		return &ListResponse{
			Total: 1,
			Posts: posts,
		}, nil
	}

	internalUserAPIMock := user.Mock{}
	internalUserAPIMock.UsersFn = func(ctx context.Context, userIDs []int) ([]model.User, error) {
		return nil, errors.New(wantErr)
	}

	redisClientMock := redis.Mock{}
	s := NewService(rMock, internalUserAPIMock, redisClientMock)

	request := ListRequest{}

	response, err := s.List(context.Background(), request)
	if err == nil {
		t.Fatalf("got: nil, want: %s", wantErr)
	}

	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	if response != nil {
		t.Fatalf("got: %v, want: nil", response.Posts)
	}
}

func TestServiceListUsersAPIError(t *testing.T) {
	wantErr := "users api error"

	rMock := repositoryMock{}
	rMock.ListFn = func(ctx context.Context, request ListRequest, perPage int) (*ListResponse, error) {
		return nil, errors.New(wantErr)
	}

	internalUserAPIMock := user.Mock{}
	internalUserAPIMock.UsersFn = func(ctx context.Context, userIDs []int) ([]model.User, error) {
		return nil, nil
	}

	redisClientMock := redis.Mock{}
	s := NewService(rMock, internalUserAPIMock, redisClientMock)

	request := ListRequest{}

	response, err := s.List(context.Background(), request)
	if err == nil {
		t.Fatalf("got: nil, want: %s", wantErr)
	}

	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	if response != nil {
		t.Fatalf("got: %v, want: nil", response.Posts)
	}
}

func TestServiceCreateSuccess(t *testing.T) {
	wantUserID := 1
	wantLogin := "anon"
	wantName := "Anon"
	wantPhotoURL := "https://test.com/test.png"

	var wantCreatedAt int64 = 100500123

	u := PostsUser{
		ID: wantUserID,
	}

	rMock := repositoryMock{}
	rMock.CreateFn = func(ctx context.Context, uid int, users []int, request CreateRequest) (*CreateResponse, error) {
		return &CreateResponse{
			ID:        1,
			Text:      "my post #1",
			CreatedAt: wantCreatedAt,
			User:      u,
		}, nil
	}

	internalUserAPIMock := user.Mock{}
	internalUserAPIMock.UsersFn = func(ctx context.Context, userIDs []int) ([]model.User, error) {
		return []model.User{
			{
				ID:       wantUserID,
				Login:    wantLogin,
				Name:     wantName,
				PhotoURL: wantPhotoURL,
			},
		}, nil
	}
	internalUserAPIMock.FollowersFn = func(ctx context.Context, uid int) ([]int, error) {
		return []int{2, 3}, nil
	}

	redisClientMock := redis.Mock{}
	s := NewService(rMock, internalUserAPIMock, redisClientMock)

	request := CreateRequest{
		Text: "my post #1",
	}

	ctx := context.WithValue(context.Background(), jwtmiddleware.ContextKeyUserID, 1)

	response, err := s.Create(ctx, request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err)
	}

	if request.Text != response.Text {
		t.Fatalf("got: %s, want: %s", response.Text, request.Text)
	}

	if wantUserID != response.User.ID {
		t.Fatalf("got: %d, want: %d", response.User.ID, wantUserID)
	}

	if wantLogin != response.User.Login {
		t.Fatalf("got: %s, want: %s", response.User.Login, wantLogin)
	}

	if wantName != response.User.Name {
		t.Fatalf("got: %s, want: %s", response.User.Name, wantName)
	}

	if wantPhotoURL != response.User.PhotoURL {
		t.Fatalf("got: %s, want: %s", response.User.PhotoURL, wantPhotoURL)
	}

	if request.Text != response.Text {
		t.Fatalf("got: %s, want: %s", response.Text, request.Text)
	}

	if wantCreatedAt != response.CreatedAt {
		t.Fatalf("got: %d, want: %d", response.CreatedAt, wantCreatedAt)
	}
}

func TestServiceFeedSuccess(t *testing.T) {
	wantTotal := 2
	wantUserID := 1
	wantLogin := "anon"
	wantName := "Anon"
	wantPhotoURL := "https://test.com/test.png"

	u := PostsUser{
		ID: wantUserID,
	}

	posts := []Post{
		{
			ID:        1,
			Text:      "post 1",
			CreatedAt: 0,
			User:      u,
		},
		{
			ID:        2,
			Text:      "post 2",
			CreatedAt: 0,
			User:      u,
		},
	}

	rMock := repositoryMock{}
	rMock.FeedFn = func(ctx context.Context, uid int, request FeedRequest, perPage int) (*FeedResponse, error) {
		return &FeedResponse{
			Total: wantTotal,
			Posts: posts,
		}, nil
	}

	internalUserAPIMock := user.Mock{}
	internalUserAPIMock.UsersFn = func(ctx context.Context, userIDs []int) ([]model.User, error) {
		return []model.User{
			{
				ID:       wantUserID,
				Login:    wantLogin,
				Name:     wantName,
				PhotoURL: wantPhotoURL,
			},
		}, nil
	}
	internalUserAPIMock.FollowersFn = func(ctx context.Context, uid int) ([]int, error) {
		return []int{2, 3}, nil
	}

	redisClientMock := redis.Mock{}
	s := NewService(rMock, internalUserAPIMock, redisClientMock)

	request := FeedRequest{
		LatestPostID: 0,
	}

	ctx := context.WithValue(context.Background(), jwtmiddleware.ContextKeyUserID, wantUserID)

	response, err := s.Feed(ctx, wantUserID, request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err)
	}

	if response.Total != wantTotal {
		t.Fatalf("got: %d, want: %d", response.Total, wantTotal)
	}

	if wantUserID != response.Posts[0].User.ID {
		t.Fatalf("got: %d, want: %d", response.Posts[0].User.ID, wantUserID)
	}

	if wantLogin != response.Posts[0].User.Login {
		t.Fatalf("got: %s, want: %s", response.Posts[0].User.Login, wantLogin)
	}

	if wantName != response.Posts[0].User.Name {
		t.Fatalf("got: %s, want: %s", response.Posts[0].User.Name, wantName)
	}

	if wantPhotoURL != response.Posts[0].User.PhotoURL {
		t.Fatalf("got: %s, want: %s", response.Posts[0].User.PhotoURL, wantPhotoURL)
	}

	if wantTotal != len(response.Posts) {
		t.Fatalf("got: %d, want: %d", len(response.Posts), wantTotal)
	}
}

func TestServiceSearchSuccess(t *testing.T) {
	wantTotal := 2
	wantUserID := 1
	wantLogin := "anon"
	wantName := "Anon"
	wantPhotoURL := "https://test.com/test.png"

	u := PostsUser{
		ID: wantUserID,
	}

	posts := []Post{
		{
			ID:        1,
			Text:      "post 1",
			CreatedAt: 0,
			User:      u,
		},
		{
			ID:        2,
			Text:      "post 2",
			CreatedAt: 0,
			User:      u,
		},
	}

	rMock := repositoryMock{}
	rMock.SearchFn = func(ctx context.Context, request SearchRequest, perPage int) ([]int, error) {
		return []int{1, 2}, nil
	}
	rMock.PostsFn = func(ctx context.Context, postsIDs []int) ([]Post, error) {
		return posts, nil
	}

	internalUserAPIMock := user.Mock{}
	internalUserAPIMock.UsersFn = func(ctx context.Context, userIDs []int) ([]model.User, error) {
		return []model.User{
			{
				ID:       wantUserID,
				Login:    wantLogin,
				Name:     wantName,
				PhotoURL: wantPhotoURL,
			},
		}, nil
	}

	redisClientMock := redis.Mock{}
	redisClientMock.HIncrByFn = func(ctx context.Context, key string, field string, incr int64) (int64, error) {
		return 1, nil
	}
	redisClientMock.SetIntSliceFn = func(ctx context.Context, key string, value []int, expiration time.Duration) error {
		return nil
	}

	s := NewService(rMock, internalUserAPIMock, redisClientMock)

	request := SearchRequest{
		Query:   "my query",
		QueryID: 0,
		Offset:  0,
	}

	ctx := context.WithValue(context.Background(), jwtmiddleware.ContextKeyUserID, wantUserID)

	response, err := s.Search(ctx, request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err)
	}

	if response.Total != wantTotal {
		t.Fatalf("got: %d, want: %d", response.Total, wantTotal)
	}

	if wantUserID != response.Posts[0].User.ID {
		t.Fatalf("got: %d, want: %d", response.Posts[0].User.ID, wantUserID)
	}

	if wantLogin != response.Posts[0].User.Login {
		t.Fatalf("got: %s, want: %s", response.Posts[0].User.Login, wantLogin)
	}

	if wantName != response.Posts[0].User.Name {
		t.Fatalf("got: %s, want: %s", response.Posts[0].User.Name, wantName)
	}

	if wantPhotoURL != response.Posts[0].User.PhotoURL {
		t.Fatalf("got: %s, want: %s", response.Posts[0].User.PhotoURL, wantPhotoURL)
	}

	if wantTotal != len(response.Posts) {
		t.Fatalf("got: %d, want: %d", len(response.Posts), wantTotal)
	}
}

func TestServiceSearchWithQueryIDSuccess(t *testing.T) {
	wantTotal := 2
	wantUserID := 1
	wantLogin := "anon"
	wantName := "Anon"
	wantPhotoURL := "https://test.com/test.png"

	u := PostsUser{
		ID: wantUserID,
	}

	posts := []Post{
		{
			ID:        1,
			Text:      "post 1",
			CreatedAt: 0,
			User:      u,
		},
		{
			ID:        2,
			Text:      "post 2",
			CreatedAt: 0,
			User:      u,
		},
	}

	rMock := repositoryMock{}
	rMock.PostsFn = func(ctx context.Context, postsIDs []int) ([]Post, error) {
		return posts, nil
	}

	internalUserAPIMock := user.Mock{}
	internalUserAPIMock.UsersFn = func(ctx context.Context, userIDs []int) ([]model.User, error) {
		return []model.User{
			{
				ID:       wantUserID,
				Login:    wantLogin,
				Name:     wantName,
				PhotoURL: wantPhotoURL,
			},
		}, nil
	}

	redisClientMock := redis.Mock{}
	redisClientMock.GetIntSliceFn = func(ctx context.Context, key string) ([]int, error) {
		return []int{1, 2}, nil
	}

	s := NewService(rMock, internalUserAPIMock, redisClientMock)

	request := SearchRequest{
		Query:   "my query",
		QueryID: 1,
		Offset:  0,
	}

	ctx := context.WithValue(context.Background(), jwtmiddleware.ContextKeyUserID, wantUserID)

	response, err := s.Search(ctx, request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err)
	}

	if response.Total != wantTotal {
		t.Fatalf("got: %d, want: %d", response.Total, wantTotal)
	}

	if wantUserID != response.Posts[0].User.ID {
		t.Fatalf("got: %d, want: %d", response.Posts[0].User.ID, wantUserID)
	}

	if wantLogin != response.Posts[0].User.Login {
		t.Fatalf("got: %s, want: %s", response.Posts[0].User.Login, wantLogin)
	}

	if wantName != response.Posts[0].User.Name {
		t.Fatalf("got: %s, want: %s", response.Posts[0].User.Name, wantName)
	}

	if wantPhotoURL != response.Posts[0].User.PhotoURL {
		t.Fatalf("got: %s, want: %s", response.Posts[0].User.PhotoURL, wantPhotoURL)
	}

	if wantTotal != len(response.Posts) {
		t.Fatalf("got: %d, want: %d", len(response.Posts), wantTotal)
	}
}
