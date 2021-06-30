package registration

import (
	"encoding/json"
	"net/http"

	"go.elastic.co/apm/module/apmzap"

	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/logger"
)

// Handler represents methods for registration HTTP server.
type Handler interface {
	// Register creates a new registration record in database and sends confirmation request to the user.
	Register(w http.ResponseWriter, r *http.Request)
	// Confirm reads user's confirmation code and confirms registration in case code is correct.
	Confirm(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	Service         Service
	Logger          logger.Logger
	ResponseBuilder api.ResponseBuilder
}

// Register creates a new registration record in database and sends confirmation request to the user.
func (h handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := RegisterRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.RequestError)
		return
	}

	msg, err := h.Service.Register(r.Context(), request)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	h.ResponseBuilder.MessageResponse(ctx, w, msg)
}

// Confirm reads user's confirmation code and confirms registration in case code is correct.
func (h handler) Confirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := ConfirmRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.RequestError)
		return
	}

	msg, err := h.Service.Confirm(r.Context(), request)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	h.ResponseBuilder.MessageResponse(ctx, w, msg)
}

// NewHandler returns instance of implemented Handler interface.
func NewHandler(s Service, l logger.Logger, rb api.ResponseBuilder) Handler {
	return handler{Service: s, Logger: l, ResponseBuilder: rb}
}
