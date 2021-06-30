package registration

import "context"

var _ Service = (*serviceMock)(nil)
var _ Repository = (*repositoryMock)(nil)

type serviceMock struct {
	RegisterFn func(context.Context, RegisterRequest) (string, error)
	ConfirmFn  func(context.Context, ConfirmRequest) (string, error)
}

type repositoryMock struct {
	RegisterFn         func(ctx context.Context, request RegisterRequest, code int) error
	ConfirmationDataFn func(ctx context.Context, email string) (*ConfirmationData, error)
	ConfirmFn          func(ctx context.Context, email string) error
}

func (r serviceMock) Register(ctx context.Context, request RegisterRequest) (string, error) {
	return r.RegisterFn(ctx, request)
}

func (r serviceMock) Confirm(ctx context.Context, request ConfirmRequest) (string, error) {
	return r.ConfirmFn(ctx, request)
}

func (r repositoryMock) Register(ctx context.Context, request RegisterRequest, code int) error {
	return r.RegisterFn(ctx, request, code)
}

func (r repositoryMock) ConfirmationData(ctx context.Context, email string) (*ConfirmationData, error) {
	return r.ConfirmationDataFn(ctx, email)
}

func (r repositoryMock) Confirm(ctx context.Context, email string) error {
	return r.ConfirmFn(ctx, email)
}
