package registration

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/pkg/errors"

	"gitlab.com/slirx/newproj/pkg/queue/manager"
	"gitlab.com/slirx/newproj/pkg/template"
	"gitlab.com/slirx/newproj/pkg/tracer"
)

func TestServiceRegisterSuccess(t *testing.T) {
	rMock := repositoryMock{}
	rMock.RegisterFn = func(ctx context.Context, request RegisterRequest, code int) error {
		return nil
	}

	tMock := tracer.Mock{}
	tMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	m := manager.Mock{}
	m.SendFn = func(ctx context.Context, routingKey string, msg interface{}) error {
		return nil
	}

	g := template.Mock{}
	g.GenerateFn = func(t template.GeneratorType, fileName string, data interface{}) (string, error) {
		return "template-content", nil
	}

	s := NewService(tMock, rMock, m, g)

	ctx := context.Background()
	request := RegisterRequest{}

	msg, err := s.Register(ctx, request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err)
	}

	want := "you have been successfully registered. please, confirm your email"
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestServiceRegisterAlreadyUsedLoginError(t *testing.T) {
	rMock := repositoryMock{}
	rMock.RegisterFn = func(ctx context.Context, request RegisterRequest, code int) error {
		return nil
	}

	tMock := tracer.Mock{}
	tMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	m := manager.Mock{}
	m.SendFn = func(ctx context.Context, routingKey string, msg interface{}) error {
		return nil
	}

	g := template.Mock{}
	g.GenerateFn = func(t template.GeneratorType, fileName string, data interface{}) (string, error) {
		return "template-content", nil
	}

	s := NewService(tMock, rMock, m, g)

	ctx := context.Background()
	request := RegisterRequest{
		Login: "admin",
	}

	wantErr := "login is already in use"

	msg, err := s.Register(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: %s", wantErr)
	}

	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestServiceRegisterRepositoryError(t *testing.T) {
	wantErr := "some error"

	rMock := repositoryMock{}
	rMock.RegisterFn = func(ctx context.Context, request RegisterRequest, code int) error {
		return errors.New(wantErr)
	}

	tMock := tracer.Mock{}
	tMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	m := manager.Mock{}
	m.SendFn = func(ctx context.Context, routingKey string, msg interface{}) error {
		return nil
	}

	g := template.Mock{}
	g.GenerateFn = func(t template.GeneratorType, fileName string, data interface{}) (string, error) {
		return "template-content", nil
	}

	s := NewService(tMock, rMock, m, g)

	ctx := context.Background()
	request := RegisterRequest{
		Login: "test-user",
	}

	msg, err := s.Register(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: %s", wantErr)
	}

	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestServiceRegisterHTMLTemplateGeneratorError(t *testing.T) {
	rMock := repositoryMock{}
	rMock.RegisterFn = func(ctx context.Context, request RegisterRequest, code int) error {
		return nil
	}

	tMock := tracer.Mock{}
	tMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	m := manager.Mock{}
	m.SendFn = func(ctx context.Context, routingKey string, msg interface{}) error {
		return nil
	}

	wantErr := "generator error"

	g := template.Mock{}
	g.GenerateFn = func(t template.GeneratorType, fileName string, data interface{}) (string, error) {
		if t == template.TypeHTML {
			return "", errors.New(wantErr)
		}

		return "some-template", nil
	}

	s := NewService(tMock, rMock, m, g)

	ctx := context.Background()
	request := RegisterRequest{
		Login: "test-user",
	}

	msg, err := s.Register(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: %s", wantErr)
	}

	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestServiceRegisterTextTemplateGeneratorError(t *testing.T) {
	rMock := repositoryMock{}
	rMock.RegisterFn = func(ctx context.Context, request RegisterRequest, code int) error {
		return nil
	}

	tMock := tracer.Mock{}
	tMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	m := manager.Mock{}
	m.SendFn = func(ctx context.Context, routingKey string, msg interface{}) error {
		return nil
	}

	wantErr := "generator error"

	g := template.Mock{}
	g.GenerateFn = func(t template.GeneratorType, fileName string, data interface{}) (string, error) {
		if t == template.TypeText {
			return "", errors.New(wantErr)
		}

		return "some-template", nil
	}

	s := NewService(tMock, rMock, m, g)

	ctx := context.Background()
	request := RegisterRequest{
		Login: "test-user",
	}

	msg, err := s.Register(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: %s", wantErr)
	}

	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestServiceRegisterManagerError(t *testing.T) {
	rMock := repositoryMock{}
	rMock.RegisterFn = func(ctx context.Context, request RegisterRequest, code int) error {
		return nil
	}

	tMock := tracer.Mock{}
	tMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	wantErr := "generator error"

	m := manager.Mock{}
	m.SendFn = func(ctx context.Context, routingKey string, msg interface{}) error {
		return errors.New(wantErr)
	}

	g := template.Mock{}
	g.GenerateFn = func(t template.GeneratorType, fileName string, data interface{}) (string, error) {
		return "some-template", nil
	}

	s := NewService(tMock, rMock, m, g)

	ctx := context.Background()
	request := RegisterRequest{
		Login: "test-user",
	}

	msg, err := s.Register(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: %s", wantErr)
	}

	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestServiceConfirmSuccess(t *testing.T) {
	repositoryMock := repositoryMock{}
	repositoryMock.ConfirmationDataFn = func(ctx context.Context, email string) (*ConfirmationData, error) {
		return &ConfirmationData{
			Code:      12345,
			Login:     "john",
			CreatedAt: time.Now(),
		}, nil
	}
	repositoryMock.ConfirmFn = func(ctx context.Context, email string) error {
		return nil
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	m := manager.Mock{}
	m.SendFn = func(ctx context.Context, routingKey string, msg interface{}) error {
		return nil
	}

	g := template.Mock{}
	s := NewService(tracerMock, repositoryMock, m, g)

	ctx := context.Background()
	request := ConfirmRequest{
		Email:                "test@test.com",
		Code:                 12345,
		Password:             "test-test-123",
		PasswordConfirmation: "test-test-123",
	}

	msg, err := s.Confirm(ctx, request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err)
	}

	want := "you have successfully confirmed your email. you can log in now"
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestServiceConfirmEmailIsUnregisteredError(t *testing.T) {
	repositoryMock := repositoryMock{}
	repositoryMock.ConfirmationDataFn = func(ctx context.Context, email string) (*ConfirmationData, error) {
		return nil, sql.ErrNoRows
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	m := manager.Mock{}
	g := template.Mock{}
	s := NewService(tracerMock, repositoryMock, m, g)

	ctx := context.Background()
	request := ConfirmRequest{
		Email:                "test@test.com",
		Code:                 12345,
		Password:             "test-test-123",
		PasswordConfirmation: "test-test-123",
	}

	msg, err := s.Confirm(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: email is unregistered")
	}

	wantErr := "email is unregistered"
	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestServiceConfirmInvalidConfirmationCodeError(t *testing.T) {
	repositoryMock := repositoryMock{}
	repositoryMock.ConfirmationDataFn = func(ctx context.Context, email string) (*ConfirmationData, error) {
		return &ConfirmationData{
			Code:      11111,
			Login:     "john",
			CreatedAt: time.Now(),
		}, nil
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	m := manager.Mock{}
	g := template.Mock{}
	s := NewService(tracerMock, repositoryMock, m, g)

	ctx := context.Background()
	request := ConfirmRequest{
		Email:                "test@test.com",
		Code:                 12345,
		Password:             "test-test-123",
		PasswordConfirmation: "test-test-123",
	}

	msg, err := s.Confirm(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: confirmation code is invalid")
	}

	wantErr := "confirmation code is invalid"
	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestServiceConfirmRepositoryConfirmationDataError(t *testing.T) {
	wantErr := "repository error"

	repositoryMock := repositoryMock{}
	repositoryMock.ConfirmationDataFn = func(ctx context.Context, email string) (*ConfirmationData, error) {
		return nil, errors.New(wantErr)
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	m := manager.Mock{}
	g := template.Mock{}
	s := NewService(tracerMock, repositoryMock, m, g)

	ctx := context.Background()
	request := ConfirmRequest{
		Email:                "test@test.com",
		Code:                 12345,
		Password:             "test-test-123",
		PasswordConfirmation: "test-test-123",
	}

	msg, err := s.Confirm(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: %s", wantErr)
	}

	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestServiceConfirmRepositoryConfirmError(t *testing.T) {
	wantErr := "repository error"

	repositoryMock := repositoryMock{}
	repositoryMock.ConfirmationDataFn = func(ctx context.Context, email string) (*ConfirmationData, error) {
		return &ConfirmationData{
			Code:      12345,
			Login:     "john",
			CreatedAt: time.Now(),
		}, nil
	}
	repositoryMock.ConfirmFn = func(ctx context.Context, email string) error {
		return errors.New(wantErr)
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	m := manager.Mock{}
	g := template.Mock{}
	s := NewService(tracerMock, repositoryMock, m, g)

	ctx := context.Background()
	request := ConfirmRequest{
		Email:                "test@test.com",
		Code:                 12345,
		Password:             "test-test-123",
		PasswordConfirmation: "test-test-123",
	}

	msg, err := s.Confirm(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: %s", wantErr)
	}

	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestServiceConfirmManagerError(t *testing.T) {
	wantErr := "manager error"

	repositoryMock := repositoryMock{}
	repositoryMock.ConfirmationDataFn = func(ctx context.Context, email string) (*ConfirmationData, error) {
		return &ConfirmationData{
			Code:      12345,
			Login:     "john",
			CreatedAt: time.Now(),
		}, nil
	}
	repositoryMock.ConfirmFn = func(ctx context.Context, email string) error {
		return nil
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	m := manager.Mock{}
	m.SendFn = func(ctx context.Context, routingKey string, msg interface{}) error {
		return errors.New(wantErr)
	}

	g := template.Mock{}
	s := NewService(tracerMock, repositoryMock, m, g)

	ctx := context.Background()
	request := ConfirmRequest{
		Email:                "test@test.com",
		Code:                 12345,
		Password:             "test-test-123",
		PasswordConfirmation: "test-test-123",
	}

	msg, err := s.Confirm(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: %s", wantErr)
	}

	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestServiceConfirmTimedOutConfirmationCodeError(t *testing.T) {
	repositoryMock := repositoryMock{}
	repositoryMock.ConfirmationDataFn = func(ctx context.Context, email string) (*ConfirmationData, error) {
		return &ConfirmationData{
			Code:      12345,
			Login:     "john",
			CreatedAt: time.Now().Add(-25 * time.Hour),
		}, nil
	}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	m := manager.Mock{}
	g := template.Mock{}
	s := NewService(tracerMock, repositoryMock, m, g)

	ctx := context.Background()
	request := ConfirmRequest{
		Email:                "test@test.com",
		Code:                 12345,
		Password:             "test-test-123",
		PasswordConfirmation: "test-test-123",
	}

	msg, err := s.Confirm(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: confirmation code is timed out")
	}

	wantErr := "confirmation code is timed out"
	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestServiceConfirmEmptyPasswordError(t *testing.T) {
	repositoryMock := repositoryMock{}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	m := manager.Mock{}
	g := template.Mock{}
	s := NewService(tracerMock, repositoryMock, m, g)

	ctx := context.Background()
	request := ConfirmRequest{
		Email:                "test@test.com",
		Code:                 12345,
		Password:             "",
		PasswordConfirmation: "",
	}

	msg, err := s.Confirm(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: password should not be empty")
	}

	wantErr := "password should not be empty"
	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestServiceConfirmDifferentPasswordError(t *testing.T) {
	repositoryMock := repositoryMock{}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	m := manager.Mock{}
	g := template.Mock{}
	s := NewService(tracerMock, repositoryMock, m, g)

	ctx := context.Background()
	request := ConfirmRequest{
		Email:                "test@test.com",
		Code:                 12345,
		Password:             "test-test-123",
		PasswordConfirmation: "test-123-test",
	}

	msg, err := s.Confirm(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: passwords should match")
	}

	wantErr := "passwords should match"
	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}

func TestServiceConfirmMinLenPasswordError(t *testing.T) {
	repositoryMock := repositoryMock{}

	tracerMock := tracer.Mock{}
	tracerMock.RequestIDFn = func(ctx context.Context) string {
		return "req1"
	}

	m := manager.Mock{}
	g := template.Mock{}
	s := NewService(tracerMock, repositoryMock, m, g)

	ctx := context.Background()
	request := ConfirmRequest{
		Email:                "test@test.com",
		Code:                 12345,
		Password:             "test",
		PasswordConfirmation: "test",
	}

	msg, err := s.Confirm(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: passwords length should be more than 5 symbols")
	}

	wantErr := "passwords length should be more than 5 symbols"
	if err.Error() != wantErr {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	want := ""
	if msg != want {
		t.Fatalf("got: %s, want: %s", msg, want)
	}
}
