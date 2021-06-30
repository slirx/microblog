package registration

import (
	"context"
	"database/sql"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/queue"
	"gitlab.com/slirx/newproj/pkg/queue/manager"
	"gitlab.com/slirx/newproj/pkg/template"
	"gitlab.com/slirx/newproj/pkg/tracer"
)

type Service interface {
	Register(ctx context.Context, request RegisterRequest) (string, error)
	Confirm(ctx context.Context, request ConfirmRequest) (string, error)
}

type service struct {
	Tracer            tracer.Tracer
	Repository        Repository
	Manager           manager.Manager
	TemplateGenerator template.Generator
	ReservedWords     map[string]struct{}
}

func (s service) Register(ctx context.Context, request RegisterRequest) (string, error) {
	var err error

	// check reserved words
	if _, ok := s.ReservedWords[request.Login]; ok {
		return "", api.NewRequestError(errors.New("login is already in use"))
	}

	// generate random code
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := 10000 + r.Intn(89999)

	if err = s.Repository.Register(ctx, request, code); err != nil {
		return "", err
	}

	var htmlTemplate string
	var textTemplate string

	requestID := s.Tracer.RequestID(ctx)
	emailConfirmation := EmailConfirmation{Code: code}

	htmlTemplate, err = s.TemplateGenerator.Generate(
		template.TypeHTML,
		"template/registration/email/confirmation.html",
		emailConfirmation,
	)
	if err != nil {
		return "", err
	}

	textTemplate, err = s.TemplateGenerator.Generate(
		template.TypeText,
		"template/registration/email/confirmation.txt",
		emailConfirmation,
	)
	if err != nil {
		return "", err
	}

	err = s.Manager.Send(ctx, queue.JobEmailSend, queue.Email{
		RequestID:      requestID,
		RecipientEmail: request.Email,
		Subject:        "Registration Confirmation",
		HTML:           htmlTemplate,
		Text:           textTemplate,
	})
	if err != nil {
		return "", err
	}

	return "you have been successfully registered. please, confirm your email", nil
}

func (s service) Confirm(ctx context.Context, request ConfirmRequest) (string, error) {
	if err := request.Validate(); err != nil {
		return "", api.NewRequestError(err)
	}

	confirmationData, err := s.Repository.ConfirmationData(ctx, request.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", api.NewRequestError(errors.New("email is unregistered"))
		}

		return "", err
	}

	if request.Code == 0 || confirmationData.Code != request.Code {
		return "", api.NewRequestError(errors.New("confirmation code is invalid"))
	}

	if confirmationData.CreatedAt.Before(time.Now().Add(-24 * time.Hour)) {
		return "", api.NewRequestError(errors.New("confirmation code is timed out"))
	}

	var passwordHash []byte

	passwordHash, err = bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if err = s.Repository.Confirm(ctx, request.Email); err != nil {
		return "", err
	}

	requestID := s.Tracer.RequestID(ctx)

	err = s.Manager.Send(ctx, queue.JobAuthCreate, queue.AuthCreate{
		RequestID: requestID,
		Login:     confirmationData.Login,
		Email:     request.Email,
		Password:  string(passwordHash),
	})
	if err != nil {
		return "", err
	}

	return "you have successfully confirmed your email. you can log in now", nil
}

func NewService(
	tracer tracer.Tracer,
	repository Repository,
	m manager.Manager,
	tg template.Generator,
) Service {
	reserved := make(map[string]struct{})
	reserved["admin"] = struct{}{}
	reserved["edit"] = struct{}{}
	reserved["settings"] = struct{}{}
	reserved["profile"] = struct{}{}
	reserved["user"] = struct{}{}

	return service{
		Tracer:            tracer,
		Repository:        repository,
		Manager:           m,
		TemplateGenerator: tg,
		ReservedWords:     reserved,
	}
}
