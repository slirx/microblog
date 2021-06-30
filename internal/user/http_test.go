package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"

	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/logger"
	"gitlab.com/slirx/newproj/pkg/tracer"
)

func TestHTTPUpdate(t *testing.T) {
	serviceMock := serviceMock{}
	serviceMock.UpdateFn = func(ctx context.Context, request UpdateRequest) (string, error) {
		return "update ok", nil
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	request := UpdateRequest{}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err)
	}

	r := httptest.NewRequest(http.MethodPatch, "/user/", bytes.NewReader(jsonRequest))
	responseBuilderMock := api.NewResponseBuilder(tracerMock)
	rec := httptest.NewRecorder()

	h := NewHandler(serviceMock, logger.NewNoop(), responseBuilderMock)
	h.Update(rec, r)

	in := rec.Body.String()
	want := `{"request_id":"req1","type":"success","message":"update ok"}
`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("got: %d, want: %d", rec.Code, http.StatusOK)
	}
}

func TestHTTPUpdateRequestError(t *testing.T) {
	serviceMock := serviceMock{}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1-err"
	}

	jsonRequest := []byte("oops")
	r := httptest.NewRequest(http.MethodPatch, "/user/", bytes.NewReader(jsonRequest))
	responseBuilderMock := api.NewResponseBuilder(tracerMock)
	rec := httptest.NewRecorder()

	h := NewHandler(serviceMock, logger.NewNoop(), responseBuilderMock)
	h.Update(rec, r)

	in := rec.Body.String()
	want := `{"request_id":"req1-err","type":"error","message":"invalid request data"}
`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got: %d, want: %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHTTPUpdateError(t *testing.T) {
	serviceMock := serviceMock{}
	serviceMock.UpdateFn = func(ctx context.Context, request UpdateRequest) (string, error) {
		return "", errors.New("update error")
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req2-err"
	}

	request := UpdateRequest{}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err)
	}

	r := httptest.NewRequest(http.MethodPatch, "/user/", bytes.NewReader(jsonRequest))
	responseBuilderMock := api.NewResponseBuilder(tracerMock)
	rec := httptest.NewRecorder()

	h := NewHandler(serviceMock, logger.NewNoop(), responseBuilderMock)
	h.Update(rec, r)

	in := rec.Body.String()
	want := `{"request_id":"req2-err","type":"error","message":"oops, something went wrong"}
`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got: %d, want: %d", rec.Code, http.StatusInternalServerError)
	}
}
