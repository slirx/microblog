package user

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"gitlab.com/slirx/newproj/internal/api/media"
	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/http/jwtmiddleware"
	"gitlab.com/slirx/newproj/pkg/queue"
	"gitlab.com/slirx/newproj/pkg/queue/manager"
	"gitlab.com/slirx/newproj/pkg/tracer"
)

const (
	usersPerPage = 20
)

type Service interface {
	Update(ctx context.Context, request UpdateRequest) (string, error)
	Get(ctx context.Context, login string) (*GetResponse, error)
	Me(ctx context.Context) (*MeResponse, error)
	Follow(ctx context.Context, request FollowRequest) error
	Unfollow(ctx context.Context, request UnfollowRequest) error
	Followers(ctx context.Context, request FollowersRequest) (*FollowersResponse, error)
	Following(ctx context.Context, request FollowingRequest) (*FollowingResponse, error)
	InternalFollowers(ctx context.Context, uid int) ([]int, error)
	InternalUsers(ctx context.Context, request InternalUsersRequest) (*InternalUsersResponse, error)
	InternalGet(ctx context.Context, login string) (*GetResponse, error)
}

type service struct {
	Repository       Repository
	Tracer           tracer.Tracer
	Manager          manager.Manager
	InternalMediaAPI media.API
}

func (s service) Update(ctx context.Context, request UpdateRequest) (string, error) {
	uid, err := jwtmiddleware.UID(ctx)
	if err != nil {
		return "", err
	}

	if err = s.Repository.Update(ctx, uid, request); err != nil {
		return "", err
	}

	return "profile have been successfully updated", nil
}

func (s service) Get(ctx context.Context, login string) (*GetResponse, error) {
	uid, err := jwtmiddleware.UID(ctx)
	if err != nil {
		return nil, err
	}

	response, err := s.Repository.Get(ctx, login, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, api.NewNotFoundError(errors.New("user not found"))
		}

		return nil, err
	}

	var images map[int]string
	images, err = s.InternalMediaAPI.Images(ctx, "user", []int{response.ID})
	if err != nil {
		return nil, err
	}

	response.PhotoURL = images[response.ID]

	return response, nil
}

func (s service) Me(ctx context.Context) (*MeResponse, error) {
	uid, err := jwtmiddleware.UID(ctx)
	if err != nil {
		return nil, err
	}

	var response *MeResponse

	response, err = s.Repository.Me(ctx, uid)
	if err != nil {
		return nil, err
	}

	var images map[int]string
	images, err = s.InternalMediaAPI.Images(ctx, "user", []int{int(uid)})
	if err != nil {
		return nil, err
	}

	response.PhotoURL = images[uid]

	//response.PhotoURL = "https://i.imgur.com/ilf6CPH.png"

	return response, nil
}

func (s service) Follow(ctx context.Context, request FollowRequest) error {
	uid, err := jwtmiddleware.UID(ctx)
	if err != nil {
		return err
	}

	requestID := s.Tracer.RequestID(ctx)

	err = s.Manager.Send(ctx, queue.JobPostFollow, queue.PostFollow{
		RequestID:    requestID,
		UserID:       uid,
		FollowUserID: request.UserID,
	})
	if err != nil {
		return err
	}

	if err = s.Repository.Follow(ctx, uid, request); err != nil {
		return err
	}

	return nil
}

func (s service) Unfollow(ctx context.Context, request UnfollowRequest) error {
	uid, err := jwtmiddleware.UID(ctx)
	if err != nil {
		return err
	}

	requestID := s.Tracer.RequestID(ctx)

	err = s.Manager.Send(ctx, queue.JobPostUnfollow, queue.PostUnfollow{
		RequestID:      requestID,
		UserID:         uid,
		UnfollowUserID: request.UserID,
	})
	if err != nil {
		return err
	}

	if err = s.Repository.Unfollow(ctx, uid, request); err != nil {
		return err
	}

	return nil
}

func (s service) Followers(ctx context.Context, request FollowersRequest) (*FollowersResponse, error) {
	response, err := s.Repository.Followers(ctx, request, 20)
	if err != nil {
		return response, err
	}

	if len(response.Followers) > 0 {
		userIDs := make([]int, 0)
		for _, f := range response.Followers {
			userIDs = append(userIDs, f.UserID)
		}

		var images map[int]string
		images, err = s.InternalMediaAPI.Images(ctx, "user", userIDs)
		if err != nil {
			return nil, err
		}

		for i, f := range response.Followers {
			response.Followers[i].PhotoURL = images[f.UserID]
		}
	}

	return response, nil
}

func (s service) Following(ctx context.Context, request FollowingRequest) (*FollowingResponse, error) {
	response, err := s.Repository.Following(ctx, request, 20)
	if err != nil {
		return response, err
	}

	if len(response.Following) > 0 {
		userIDs := make([]int, 0)
		for _, f := range response.Following {
			userIDs = append(userIDs, f.UserID)
		}

		var images map[int]string
		images, err = s.InternalMediaAPI.Images(ctx, "user", userIDs)
		if err != nil {
			return nil, err
		}

		for i, f := range response.Following {
			response.Following[i].PhotoURL = images[f.UserID]
		}
	}

	return response, nil
}

func (s service) InternalFollowers(ctx context.Context, uid int) ([]int, error) {
	response, err := s.Repository.FollowersIDs(ctx, uid)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s service) InternalUsers(ctx context.Context, request InternalUsersRequest) (*InternalUsersResponse, error) {
	response, err := s.Repository.Users(ctx, request, usersPerPage)
	if err != nil {
		return response, err
	}

	userIDs := make([]int, 0)
	for _, user := range response.Users {
		userIDs = append(userIDs, user.ID)
	}

	var images map[int]string
	images, err = s.InternalMediaAPI.Images(ctx, "user", userIDs)
	if err != nil {
		return nil, err
	}

	for i, user := range response.Users {
		response.Users[i].PhotoURL = images[user.ID]
	}

	return response, nil
}

func (s service) InternalGet(ctx context.Context, login string) (*GetResponse, error) {
	response, err := s.Repository.Get(ctx, login, 0)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, api.NewNotFoundError(errors.New("user not found"))
		}

		return nil, err
	}

	var images map[int]string
	images, err = s.InternalMediaAPI.Images(ctx, "user", []int{response.ID})
	if err != nil {
		return nil, err
	}

	response.PhotoURL = images[response.ID]

	return response, nil
}

func NewService(
	repository Repository,
	tracer tracer.Tracer,
	m manager.Manager,
	internalMediaAPI media.API,
) Service {
	return service{
		Repository:       repository,
		Tracer:           tracer,
		Manager:          m,
		InternalMediaAPI: internalMediaAPI,
	}
}
