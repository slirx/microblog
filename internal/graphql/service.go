package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"gitlab.com/slirx/newproj/internal/api"
	"gitlab.com/slirx/newproj/internal/graphql/graph/model"
)

type Service struct {
	UserService UserService
}

type UserService interface {
	User(ctx context.Context, login string) (*model.User, error)
	Users(ctx context.Context, latestUserID int) (*model.UsersResponse, error)
}

func NewService(endpoints map[string]string, config *api.ServiceConfig) (*Service, error) {
	// todo fetch token from /internal/auth/login
	var err error

	a := api.GeneralAPI{
		ServiceConfig: config,
		Client: http.Client{
			Timeout: 3 * time.Second,
		},
		Endpoints: endpoints,
	}
	if err = a.Login("auth"); err != nil {
		return nil, err
	}

	s := &Service{
		UserService: userAPI{
			GeneralAPI: a,
		},
	}

	return s, nil
}

type userAPI struct {
	api.GeneralAPI
	Token string
}

type GetUserResponse struct {
	Data model.User `json:"data"`
}

type GetUsersResponse struct {
	Data model.UsersResponse `json:"data"`
}

func (a userAPI) User(ctx context.Context, login string) (*model.User, error) {
	body, err := a.SendRequest(ctx, "user", "GET", "internal/user/"+login, nil)
	if err != nil {
		return nil, err
	}

	response := GetUserResponse{}
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, errors.WithStack(err)
	}

	// todo handle error. for example where there are no data field in response and we've got 500 error

	return &response.Data, nil
}

func (a userAPI) Users(ctx context.Context, latestUserID int) (*model.UsersResponse, error) {
	body, err := a.SendRequest(
		ctx,
		"user",
		"GET",
		fmt.Sprintf("internal/user?latest_user_id=%d", latestUserID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	response := GetUsersResponse{}
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, errors.WithStack(err)
	}

	// todo handle error. for example where there are no data field in response and we've got 500 error

	return &response.Data, nil
}
