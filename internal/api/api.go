package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type GeneralAPI struct {
	ServiceConfig *ServiceConfig
	Client        http.Client
	Endpoints     map[string]string
	Token         string
}

type InternalJWT struct {
	Endpoint string
	Login    string
	Password string
}

type ServiceConfig struct {
	InternalJWT InternalJWT
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Data struct {
		AccessToken string `json:"access_token"`
	} `json:"data"`
}

func (a *GeneralAPI) Login(serviceName string) error {
	request := LoginRequest{
		Login:    a.ServiceConfig.InternalJWT.Login,
		Password: a.ServiceConfig.InternalJWT.Password,
	}

	body, err := a.SendRequest(context.Background(), serviceName, "POST", "internal/auth/login", request)
	if err != nil {
		return err
	}

	response := LoginResponse{}
	if err = json.Unmarshal(body, &response); err != nil {
		return errors.WithStack(err)
	}

	if response.Data.AccessToken == "" {
		return errors.WithStack(errors.New("empty access token"))
	}

	a.Token = response.Data.AccessToken

	return nil
}

func (a *GeneralAPI) SendRequest(
	ctx context.Context,
	serviceName string,
	method string,
	endpoint string,
	request interface{},
) ([]byte, error) {
	var err error
	var r *http.Request

	baseURL, ok := a.Endpoints[serviceName]
	if !ok || baseURL == "" {
		return nil, errors.WithStack(fmt.Errorf("base URL is undefined for service %s", serviceName))
	}

	// todo track as APM request (child transaction)

	if request != nil {
		b, err := json.Marshal(request)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		body := bytes.NewReader(b)
		r, err = http.NewRequestWithContext(ctx, method, baseURL+endpoint, body)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		r, err = http.NewRequestWithContext(ctx, method, baseURL+endpoint, nil)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	r.Header.Set("Authorization", a.Token)
	var resp *http.Response

	resp, err = a.Client.Do(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer resp.Body.Close()

	var response []byte

	response, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// todo re-login in case token is expired

	// todo maybe I have to return *http.Response here
	return response, nil
}
