package user

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmzap"

	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/logger"
)

// Handler represents methods for user HTTP server.
type Handler interface {
	// Update updates user's information.
	Update(w http.ResponseWriter, r *http.Request)
	Me(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Follow(w http.ResponseWriter, r *http.Request)
	Unfollow(w http.ResponseWriter, r *http.Request)
	Followers(w http.ResponseWriter, r *http.Request)
	Following(w http.ResponseWriter, r *http.Request)
	InternalFollowers(w http.ResponseWriter, r *http.Request)
	InternalUsers(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	Service         Service
	Logger          logger.Logger
	ResponseBuilder api.ResponseBuilder
}

// Update updates user's information.
func (h handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := UpdateRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.RequestError)
		return
	}

	msg, err := h.Service.Update(r.Context(), request)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	h.ResponseBuilder.MessageResponse(ctx, w, msg)
}

func (h handler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	response, err := h.Service.Me(ctx)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	h.ResponseBuilder.DataResponse(ctx, w, response)
}

func (h handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	login := chi.URLParam(r, "login")
	if login == "" {
		err := api.NewRequestError(errors.New("invalid login"))
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	response, err := h.Service.Get(ctx, login)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	h.ResponseBuilder.DataResponse(ctx, w, response)
}

func (h handler) Follow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := FollowRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.RequestError)
		return
	}

	err := h.Service.Follow(r.Context(), request)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	h.ResponseBuilder.DataResponse(ctx, w, nil)
}

func (h handler) Unfollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := UnfollowRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.RequestError)
		return
	}

	err := h.Service.Unfollow(r.Context(), request)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	h.ResponseBuilder.DataResponse(ctx, w, nil)
}

func (h handler) Followers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	login := chi.URLParam(r, "login")
	if login == "" {
		err := api.NewRequestError(errors.New("invalid login"))
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	latestFollowerID, _ := strconv.Atoi(r.URL.Query().Get("lfid"))
	request := FollowersRequest{
		Login:            login,
		LatestFollowerID: latestFollowerID,
	}

	response, err := h.Service.Followers(ctx, request)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	h.ResponseBuilder.DataResponse(ctx, w, response)
}

func (h handler) Following(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	login := chi.URLParam(r, "login")
	if login == "" {
		err := api.NewRequestError(errors.New("invalid login"))
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	latestFollowerID, _ := strconv.Atoi(r.URL.Query().Get("lfid"))
	request := FollowingRequest{
		Login:            login,
		LatestFollowerID: latestFollowerID,
	}

	response, err := h.Service.Following(ctx, request)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	h.ResponseBuilder.DataResponse(ctx, w, response)
}

func (h handler) InternalFollowers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	uid, err := strconv.Atoi(chi.URLParam(r, "uid"))
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, api.RequestError)
		return
	}

	if uid == 0 {
		err := api.NewRequestError(errors.New("invalid user id"))
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	response, err := h.Service.InternalFollowers(ctx, uid)
	if err != nil {
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	h.ResponseBuilder.DataResponse(ctx, w, response)
}

func (h handler) InternalUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idsSlice := strings.Split(r.URL.Query().Get("ids"), ",")
	ids := make([]int, 0)

	var id int
	for _, s := range idsSlice {
		id, _ = strconv.Atoi(s)
		if id == 0 {
			continue
		}

		ids = append(ids, id)
	}

	if len(ids) == 0 {
		err := api.NewRequestError(errors.New("no user ids"))
		h.Logger.Error(err, apmzap.TraceContext(ctx)...)
		h.ResponseBuilder.ErrorResponse(ctx, w, err)
		return
	}

	response, err := h.Service.InternalUsers(ctx, ids)
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
