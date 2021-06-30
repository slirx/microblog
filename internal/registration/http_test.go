package registration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/logger"
	"gitlab.com/slirx/newproj/pkg/tracer"
)

func TestHTTPRegisterSuccess(t *testing.T) {
	serviceMock := serviceMock{}
	serviceMock.RegisterFn = func(ctx context.Context, request RegisterRequest) (string, error) {
		return "success registration", nil
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	request := RegisterRequest{}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err)
	}

	r := httptest.NewRequest(http.MethodPost, "/register/", bytes.NewReader(jsonRequest))
	responseBuilderMock := api.NewResponseBuilder(tracerMock)
	rec := httptest.NewRecorder()

	h := NewHandler(serviceMock, logger.NewNoop(), responseBuilderMock)
	h.Register(rec, r)

	in := rec.Body.String()
	want := `{"request_id":"req1","type":"success","message":"success registration"}
`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("got: %d, want: %d", rec.Code, http.StatusOK)
	}
}

func TestHTTPRegisterRequestError(t *testing.T) {
	serviceMock := serviceMock{}
	serviceMock.RegisterFn = func(ctx context.Context, request RegisterRequest) (string, error) {
		return "", errors.New("oops")
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1-err"
	}

	jsonRequest := []byte("oops")
	r := httptest.NewRequest(http.MethodPost, "/register/", bytes.NewReader(jsonRequest))

	rec := httptest.NewRecorder()

	h := NewHandler(serviceMock, logger.NewNoop(), api.NewResponseBuilder(tracerMock))
	h.Register(rec, r)

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

func TestHTTPRegisterServiceError(t *testing.T) {
	serviceMock := serviceMock{}
	serviceMock.RegisterFn = func(ctx context.Context, request RegisterRequest) (string, error) {
		return "", errors.New("oops")
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1-err"
	}

	request := RegisterRequest{}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err)
	}

	r := httptest.NewRequest(http.MethodPost, "/register/", bytes.NewReader(jsonRequest))
	rec := httptest.NewRecorder()

	h := NewHandler(serviceMock, logger.NewNoop(), api.NewResponseBuilder(tracerMock))
	h.Register(rec, r)

	in := rec.Body.String()
	want := `{"request_id":"req1-err","type":"error","message":"oops, something went wrong"}
`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got: %d, want: %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestHTTPConfirmSuccess(t *testing.T) {
	serviceMock := serviceMock{}
	serviceMock.ConfirmFn = func(ctx context.Context, request ConfirmRequest) (string, error) {
		return "success", nil
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	rec := httptest.NewRecorder()

	jsonRequest := []byte(`{
    "email": "test@test.com",
    "code": 12345,
    "password": "qwerty125",
    "password_confirmation": "qwerty125"
}`)
	r := httptest.NewRequest(http.MethodPost, "/register/confirm/", bytes.NewReader(jsonRequest))

	h := NewHandler(serviceMock, logger.NewNoop(), api.NewResponseBuilder(tracerMock))
	h.Confirm(rec, r)

	in := rec.Body.String()
	want := `{"request_id":"req1","type":"success","message":"success"}
`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("got: %d, want: %d", rec.Code, http.StatusOK)
	}
}

func TestHTTPConfirmRequestError(t *testing.T) {
	serviceMock := serviceMock{}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	rec := httptest.NewRecorder()

	jsonRequest := []byte("oops")
	r := httptest.NewRequest(http.MethodPost, "/register/confirm/", bytes.NewReader(jsonRequest))

	h := NewHandler(serviceMock, logger.NewNoop(), api.NewResponseBuilder(tracerMock))
	h.Confirm(rec, r)

	in := rec.Body.String()
	want := `{"request_id":"req1","type":"error","message":"invalid request data"}
`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got: %d, want: %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHTTPConfirmServiceError(t *testing.T) {
	serviceMock := serviceMock{}
	serviceMock.ConfirmFn = func(ctx context.Context, request ConfirmRequest) (string, error) {
		return "", errors.New("oops")
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	rec := httptest.NewRecorder()

	jsonRequest := []byte(`{
    "email": "test@test.com",
    "code": 12345,
    "password": "qwerty125",
    "password_confirmation": "qwerty125"
}`)
	r := httptest.NewRequest(http.MethodPost, "/register/confirm/", bytes.NewReader(jsonRequest))

	h := NewHandler(serviceMock, logger.NewNoop(), api.NewResponseBuilder(tracerMock))
	h.Confirm(rec, r)

	in := rec.Body.String()
	want := `{"request_id":"req1","type":"error","message":"oops, something went wrong"}
`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got: %d, want: %d", rec.Code, http.StatusBadRequest)
	}
}
