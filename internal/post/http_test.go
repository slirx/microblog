package post

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/logger"
	"gitlab.com/slirx/newproj/pkg/tracer"
)

func TestHTTPListSuccess(t *testing.T) {
	u := PostsUser{
		ID:       2,
		Login:    "test-user",
		Name:     "Test",
		PhotoURL: "https://test.com/test.png",
	}

	serviceMock := serviceMock{}
	serviceMock.ListFn = func(ctx context.Context, request ListRequest) (*ListResponse, error) {
		return &ListResponse{
			Total: 2,
			Posts: []Post{
				{
					ID:        4,
					Text:      "post 4",
					CreatedAt: 100501,
					User:      u,
				},
				{
					ID:        3,
					Text:      "post 3",
					CreatedAt: 100500,
					User:      u,
				},
			},
		}, nil
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	userID := "2"

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("uid", userID)

	r := httptest.NewRequest(http.MethodGet, "/post/user/"+userID, nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routeCtx))

	responseBuilderMock := api.NewResponseBuilder(tracerMock)
	rec := httptest.NewRecorder()

	h := NewHandler(serviceMock, logger.NewNoop(), responseBuilderMock)
	h.List(rec, r)

	in := rec.Body.String()
	want := `{"request_id":"req1","type":"success","data":{"total":2,"posts":[{"id":4,"text":"post 4","created_at":100501,"comments_count":0,"likes_count":0,"reposts_count":0,"user":{"id":2,"login":"test-user","name":"Test","photo_url":"https://test.com/test.png"}},{"id":3,"text":"post 3","created_at":100500,"comments_count":0,"likes_count":0,"reposts_count":0,"user":{"id":2,"login":"test-user","name":"Test","photo_url":"https://test.com/test.png"}}]}}
`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("got: %d, want: %d", rec.Code, http.StatusOK)
	}
}

func TestHTTPListRequestError(t *testing.T) {
	serviceMock := serviceMock{}
	serviceMock.ListFn = func(ctx context.Context, request ListRequest) (*ListResponse, error) {
		return &ListResponse{}, nil
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req2"
	}

	r := httptest.NewRequest(http.MethodGet, "/post/user/", nil)

	responseBuilderMock := api.NewResponseBuilder(tracerMock)
	rec := httptest.NewRecorder()

	h := NewHandler(serviceMock, logger.NewNoop(), responseBuilderMock)
	h.List(rec, r)

	in := rec.Body.String()
	want := `{"request_id":"req2","type":"error","message":"invalid request data"}
`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got: %d, want: %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHTTPListInvalidUserIDError(t *testing.T) {
	serviceMock := serviceMock{}
	serviceMock.ListFn = func(ctx context.Context, request ListRequest) (*ListResponse, error) {
		return &ListResponse{}, nil
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	userID := "0"

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("uid", userID)

	r := httptest.NewRequest(http.MethodGet, "/post/user/"+userID, nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routeCtx))

	responseBuilderMock := api.NewResponseBuilder(tracerMock)
	rec := httptest.NewRecorder()

	h := NewHandler(serviceMock, logger.NewNoop(), responseBuilderMock)
	h.List(rec, r)

	in := rec.Body.String()
	want := `{"request_id":"req1","type":"error","message":"invalid user id"}
`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got: %d, want: %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHTTPListServiceError(t *testing.T) {
	wantErr := "service error"

	serviceMock := serviceMock{}
	serviceMock.ListFn = func(ctx context.Context, request ListRequest) (*ListResponse, error) {
		return nil, errors.New(wantErr)
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	userID := "2"

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("uid", userID)

	r := httptest.NewRequest(http.MethodGet, "/post/user/"+userID, nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routeCtx))

	responseBuilderMock := api.NewResponseBuilder(tracerMock)
	rec := httptest.NewRecorder()

	h := NewHandler(serviceMock, logger.NewNoop(), responseBuilderMock)
	h.List(rec, r)

	in := rec.Body.String()
	want := `{"request_id":"req1","type":"error","message":"oops, something went wrong"}
`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got: %d, want: %d", rec.Code, http.StatusInternalServerError)
	}
}
