package post

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gitlab.com/slirx/newproj/internal/api/user"
	"gitlab.com/slirx/newproj/internal/graphql/graph/model"
	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/http/jwtmiddleware"
	"gitlab.com/slirx/newproj/pkg/redis"
)

const (
	perPage          int = 20
	maxSearchPerPage int = 200
)

type Service interface {
	List(ctx context.Context, request ListRequest) (*ListResponse, error)
	Create(ctx context.Context, request CreateRequest) (*CreateResponse, error)
	Feed(ctx context.Context, userID int, request FeedRequest) (*FeedResponse, error)
	Search(ctx context.Context, request SearchRequest) (*SearchResponse, error)
}

type service struct {
	Repository      Repository
	InternalUserAPI user.API
	RedisClient     redis.Client
}

func (s service) List(ctx context.Context, request ListRequest) (*ListResponse, error) {
	response, err := s.Repository.List(ctx, request, 20)
	if err != nil {
		return nil, err
	}

	response.Posts, err = s.fetchUserInfo(ctx, response.Posts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s service) Create(ctx context.Context, request CreateRequest) (*CreateResponse, error) {
	uid, err := jwtmiddleware.UID(ctx)
	if err != nil {
		return nil, err
	}

	var users []int

	// fetch followers
	users, err = s.InternalUserAPI.Followers(ctx, uid)
	if err != nil {
		return nil, err
	}

	var response *CreateResponse
	response, err = s.Repository.Create(ctx, uid, users, request)
	if err != nil {
		return nil, err
	}

	var posts []Post

	posts, err = s.fetchUserInfo(ctx, []Post{Post(*response)})
	if err != nil {
		return nil, err
	}

	response.User = posts[0].User

	return response, nil
}

func (s service) Feed(ctx context.Context, userID int, request FeedRequest) (*FeedResponse, error) {
	uid, err := jwtmiddleware.UID(ctx)
	if err != nil {
		return nil, err
	}

	if uid != userID {
		return nil, api.NewRequestError(errors.New("invalid user id"))
	}

	response, err := s.Repository.Feed(ctx, uid, request, 10)
	if err != nil {
		return nil, err
	}

	response.Posts, err = s.fetchUserInfo(ctx, response.Posts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s service) Search(ctx context.Context, request SearchRequest) (*SearchResponse, error) {
	response := SearchResponse{
		Posts: make([]Post, 0),
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request.Query = reg.ReplaceAllString(request.Query, "")
	request.Query = strings.ReplaceAll(request.Query, "  ", " ")
	request.Query = strings.ReplaceAll(request.Query, " ", " | ")

	var postsIDs []int

	userID, err := jwtmiddleware.UID(ctx)
	if err != nil {
		return nil, err
	}

	if request.QueryID > 0 {
		response.QueryID = request.QueryID

		key := fmt.Sprintf("post:search:%d:%d", userID, request.QueryID)
		postsIDs, err = s.RedisClient.GetIntSlice(ctx, key)
		if err != nil {
			if !errors.Is(err, redis.ErrNoData) {
				return nil, err
			}

			return &response, nil
		}

		if len(postsIDs) == 0 {
			return &response, nil
		}

		response.Total = len(postsIDs)
	}

	if len(postsIDs) == 0 {
		request.Offset = 0

		postsIDs, err = s.Repository.Search(ctx, request, maxSearchPerPage)
		if err != nil {
			return nil, err
		}

		response.Total = len(postsIDs)

		var queryID int64

		// get query id
		queryID, err = s.RedisClient.HIncrBy(
			ctx,
			fmt.Sprintf("post:search.query:%d", userID),
			"query.id",
			1,
		)
		if err != nil {
			return nil, err
		}

		response.QueryID = int(queryID)

		key := fmt.Sprintf("post:search:%d:%d", userID, queryID)

		// save posts to redis
		err = s.RedisClient.SetIntSlice(ctx, key, postsIDs, time.Minute*30)
		if err != nil {
			return nil, err
		}
	}

	if len(postsIDs) > 0 {
		postsLen := len(postsIDs)
		if request.Offset > postsLen {
			return &response, nil
		}

		start := 0
		if request.Offset > 0 {
			start = request.Offset
		}

		end := postsLen - request.Offset
		if end > perPage {
			end = perPage
		}

		end = start + end
		postsIDs = postsIDs[start:end]

		response.Posts, err = s.Repository.Posts(ctx, postsIDs)
		if err != nil {
			return nil, err
		}

		response.Posts, err = s.fetchUserInfo(ctx, response.Posts)
		if err != nil {
			return nil, err
		}
	}

	return &response, nil
}

func (s service) fetchUserInfo(ctx context.Context, posts []Post) ([]Post, error) {
	userIDs := make([]int, 0)
	userIDsMap := make(map[int]struct{})

	for _, p := range posts {
		userIDsMap[p.User.ID] = struct{}{}
	}

	for id := range userIDsMap {
		userIDs = append(userIDs, id)
	}

	var users []model.User
	var err error

	users, err = s.InternalUserAPI.Users(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	usersMap := make(map[int]model.User)
	for _, u := range users {
		usersMap[u.ID] = u
	}

	var u model.User
	for i, p := range posts {
		u, _ = usersMap[p.User.ID]
		posts[i].User.Name = u.Name
		posts[i].User.Login = u.Login
		posts[i].User.PhotoURL = u.PhotoURL
	}

	return posts, nil
}

func NewService(repository Repository, internalUserAPI user.API, redisClient redis.Client) Service {
	return service{
		Repository:      repository,
		InternalUserAPI: internalUserAPI,
		RedisClient:     redisClient,
	}
}
