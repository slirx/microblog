package post

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmzap"

	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/logger"
)

// Handler represents methods for HTTP server.
type Handler interface {
	// List returns list of user's posts.
	List(w http.ResponseWriter, r *http.Request)
	// Create saves a new post.
	Create(w http.ResponseWriter, r *http.Request)
	// Feed returns user's feed.
	Feed(w http.ResponseWriter, r *http.Request)
	// Search searches across all posts in service.
	Search(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	Service         Service
	Logger          logger.Logger
	ResponseBuilder api.ResponseBuilder
}

// List returns list of posts.
func (h handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := ListRequest{}
	request.LatestPostID, _ = strconv.Atoi(r.URL.Query().Get("lpid"))

	var err error
	request.UserID, err = strconv.Atoi(chi.URLParam(r, "uid"))
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.RequestError)
		return
	}

	if request.UserID == 0 {
		err := api.NewRequestError(errors.New("invalid user id"))
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	response, err := h.Service.List(r.Context(), request)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	h.ResponseBuilder.DataResponse(ctx, w, response)
}

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := CreateRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.RequestError)
		return
	}

	response, err := h.Service.Create(r.Context(), request)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	h.ResponseBuilder.DataResponse(ctx, w, response)
}

func (h handler) Feed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := FeedRequest{}

	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.RequestError)
		return
	}

	request.LatestPostID, _ = strconv.Atoi(r.URL.Query().Get("lpid"))

	if userID == 0 {
		err := api.NewRequestError(errors.New("invalid user id"))
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	response, err := h.Service.Feed(r.Context(), userID, request)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	h.ResponseBuilder.DataResponse(ctx, w, response)
}

func (h handler) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := SearchRequest{}
	request.Offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	request.QueryID, _ = strconv.Atoi(r.URL.Query().Get("query_id"))
	request.Query = r.URL.Query().Get("query")

	response, err := h.Service.Search(r.Context(), request)
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
