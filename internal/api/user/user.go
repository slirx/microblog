package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"gitlab.com/slirx/newproj/internal/api"
	"gitlab.com/slirx/newproj/internal/graphql/graph/model"
)

type API interface {
	Followers(ctx context.Context, uid int) ([]int, error)
	Users(ctx context.Context, userIDs []int) ([]model.User, error)
}

type userAPI struct {
	GeneralAPI api.GeneralAPI
}

type followersResponse struct {
	Data []int `json:"data"`
}

type usersResponse struct {
	Data []model.User `json:"data"`
}

func (u userAPI) Followers(ctx context.Context, uid int) ([]int, error) {
	body, err := u.GeneralAPI.SendRequest(
		ctx,
		"user",
		"GET",
		"internal/user/"+strconv.Itoa(uid)+"/followers",
		nil,
	)
	if err != nil {
		return nil, err
	}

	response := followersResponse{}
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, errors.WithStack(err)
	}

	return response.Data, nil
}

func (u userAPI) Users(ctx context.Context, userIDs []int) ([]model.User, error) {
	ids := ""
	for _, id := range userIDs {
		if ids != "" {
			ids += ","
		}

		ids += strconv.Itoa(id)
	}

	body, err := u.GeneralAPI.SendRequest(
		ctx,
		"user",
		"GET",
		"internal/user?ids="+ids,
		nil,
	)
	if err != nil {
		return nil, err
	}

	response := usersResponse{}
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, errors.WithStack(err)
	}

	return response.Data, nil
}

func NewAPI(endpoints map[string]string, config *api.ServiceConfig) (API, error) {
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

	s := userAPI{
		GeneralAPI: a,
	}

	return s, nil
}
