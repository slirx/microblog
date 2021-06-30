package auth

import (
	"encoding/json"
	"net/http"

	"go.elastic.co/apm/module/apmzap"

	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/logger"
)

// todo check corresponding to https://github.com/shieldfy/API-Security-Checklist

// Handler represents methods for auth HTTP server.
type Handler interface {
	// Login checks login/password and returns JWT.
	Login(w http.ResponseWriter, r *http.Request)
	// InternalLogin checks login/password and returns JWT. It's used for service-to-service authorization.
	InternalLogin(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	Service         Service
	Logger          logger.Logger
	ResponseBuilder api.ResponseBuilder
}

// Login checks login/password and returns JWT.
func (h handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.RequestError)
		return
	}

	response, err := h.Service.Login(r.Context(), request)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	h.ResponseBuilder.DataResponse(ctx, w, response)
}

func (h handler) InternalLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := InternalLoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.RequestError)
		return
	}

	response, err := h.Service.InternalLogin(r.Context(), request)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	h.ResponseBuilder.DataResponse(ctx, w, response)
}

// NewHandler returns instance of implemented Handler interface.
func NewHandler(s Service, l logger.Logger, rb api.ResponseBuilder) Handler {
	return handler{Service: s, Logger: l, ResponseBuilder: rb}
}
