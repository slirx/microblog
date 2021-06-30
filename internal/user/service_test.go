package user

import (
	"context"
	"testing"

	"github.com/pkg/errors"

	"gitlab.com/slirx/newproj/internal/api/media"
	"gitlab.com/slirx/newproj/pkg/http/jwtmiddleware"
	"gitlab.com/slirx/newproj/pkg/queue/manager"
	"gitlab.com/slirx/newproj/pkg/tracer"
)

func TestUpdate(t *testing.T) {
	rMock := repositoryMock{}
	rMock.UpdateFn = func(ctx context.Context, uid int, request UpdateRequest) error {
		if uid != 1 {
			t.Fatalf("got: %d, want: 1", uid)
		}

		return nil
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	managerMock := manager.Mock{}
	//tracerMock.RequestIDFn = func(ctx context.Context) string {
	//	return "req1"
	//}

	internalMediaAPIMock := media.Mock{}

	s := NewService(rMock, tracerMock, managerMock, internalMediaAPIMock)

	ctx := context.Background()
	ctx = context.WithValue(ctx, jwtmiddleware.ContextKeyUserID, 1)

	request := UpdateRequest{}

	msg, err := s.Update(ctx, request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err)
	}

	want := "profile have been successfully updated"
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestUpdateUIDError(t *testing.T) {
	rMock := repositoryMock{}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	managerMock := manager.Mock{}
	//tracerMock.RequestIDFn = func(ctx context.Context) string {
	//	return "req1"
	//}

	internalMediaAPIMock := media.Mock{}

	s := NewService(rMock, tracerMock, managerMock, internalMediaAPIMock)

	ctx := context.Background()
	request := UpdateRequest{}

	msg, err := s.Update(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: uid is not in context")
	}

	wantErr := "uid is not in context"
	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestUpdateRepositoryError(t *testing.T) {
	rMock := repositoryMock{}
	rMock.UpdateFn = func(ctx context.Context, uid int, request UpdateRequest) error {
		return errors.New("update test repository error")
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	managerMock := manager.Mock{}
	//tracerMock.RequestIDFn = func(ctx context.Context) string {
	//	return "req1"
	//}
	// todo remove commented code?
	internalMediaAPIMock := media.Mock{}

	s := NewService(rMock, tracerMock, managerMock, internalMediaAPIMock)

	ctx := context.Background()
	ctx = context.WithValue(ctx, jwtmiddleware.ContextKeyUserID, 1)

	request := UpdateRequest{}

	msg, err := s.Update(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: update test repository error")
	}

	wantErr := "update test repository error"
	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}
