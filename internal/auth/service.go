package auth

import (
	"context"
	"database/sql"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/slirx/newproj/pkg/api"
)

type Service interface {
	Login(ctx context.Context, request LoginRequest) (*LoginResponse, error)
	InternalLogin(ctx context.Context, request InternalLoginRequest) (*InternalLoginResponse, error)
	AdminLogin(ctx context.Context, request AdminLoginRequest) (*AdminLoginResponse, error)
}

type service struct {
	Repository    Repository
	Config        Config
	RandGenerator *rand.Rand
}

func (s service) Login(ctx context.Context, request LoginRequest) (*LoginResponse, error) {
	data, err := s.Repository.Auth(ctx, request.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, api.NewRequestError(errors.New("login/password is incorrect"))
		}

		return nil, err
	}

	if data.UserID == 0 {
		return nil, api.NewRequestError(errors.New("login/password is incorrect"))
	}

	if err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(request.Password)); err != nil {
		return nil, api.NewRequestError(errors.New("login/password is incorrect"))
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // todo decrease time to 1 hour
	claims["jti"] = s.RandGenerator.Uint64()
	claims["uid"] = data.UserID

	accessToken, err := token.SignedString([]byte(s.Config.Secret))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// todo add refresh token

	response := &LoginResponse{
		AccessToken: accessToken,
	}

	return response, nil
}

func (s service) InternalLogin(ctx context.Context, request InternalLoginRequest) (*InternalLoginResponse, error) {
	data, err := s.Repository.InternalAuth(ctx, request.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, api.NewRequestError(errors.New("login/password is incorrect"))
		}

		return nil, err
	}

	if data.ServiceName == "" {
		return nil, api.NewRequestError(errors.New("login/password is incorrect"))
	}

	if err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(request.Password)); err != nil {
		return nil, api.NewRequestError(errors.New("login/password is incorrect"))
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // todo decrease time to 1 hour
	claims["jti"] = s.RandGenerator.Uint64()
	claims["login"] = request.Login

	secret, ok := s.Config.InternalSecrets[request.Login]
	if !ok || secret == "" {
		return nil, api.NewRequestError(errors.New("service is incorrect"))
	}

	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// todo add refresh token?

	response := &InternalLoginResponse{
		AccessToken: accessToken,
	}

	return response, nil
}

func (s service) AdminLogin(ctx context.Context, request AdminLoginRequest) (*AdminLoginResponse, error) {
	data, err := s.Repository.AdminAuth(ctx, request.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, api.NewRequestError(errors.New("login/password is incorrect"))
		}

		return nil, err
	}

	if data.UserID == 0 {
		return nil, api.NewRequestError(errors.New("login/password is incorrect"))
	}

	if err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(request.Password)); err != nil {
		return nil, api.NewRequestError(errors.New("login/password is incorrect"))
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	claims["jti"] = s.RandGenerator.Uint64()
	claims["uid"] = data.UserID

	accessToken, err := token.SignedString([]byte(s.Config.Secret))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// todo add refresh token

	response := &AdminLoginResponse{
		AccessToken: accessToken,
	}

	return response, nil
}

func NewService(
	repository Repository,
	config Config,
) Service {
	return service{
		Repository:    repository,
		Config:        config,
		RandGenerator: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
