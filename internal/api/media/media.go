package media

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"gitlab.com/slirx/newproj/internal/api"
)

type API interface {
	Images(ctx context.Context, serviceName string, ids []int) (map[int]string, error)
}

type mediaAPI struct {
	GeneralAPI api.GeneralAPI
}

type imagesResponse struct {
	Data struct {
		Images map[string]string `json:"images"`
	} `json:"data"`
}

func (u mediaAPI) Images(ctx context.Context, serviceName string, ids []int) (map[int]string, error) {
	idsString := ""
	for _, id := range ids {
		if id == 0 {
			continue
		}

		if idsString != "" {
			idsString += ","
		}

		idsString += strconv.Itoa(id)
	}

	body, err := u.GeneralAPI.SendRequest(
		ctx,
		"media",
		"GET",
		"internal/media?service="+serviceName+"&ids="+idsString,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp := imagesResponse{}
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, errors.WithStack(err)
	}

	response := make(map[int]string)
	var id int

	for i, img := range resp.Data.Images {
		id, _ = strconv.Atoi(i)
		if id == 0 {
			continue
		}

		response[id] = img
	}

	return response, nil
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

	s := mediaAPI{
		GeneralAPI: a,
	}

	return s, nil
}
