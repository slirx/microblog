package auth

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"gitlab.com/slirx/newproj/pkg/queue"
)

type Repository interface {
	Create(ctx context.Context, request queue.AuthCreate) error
	UpdateUserID(ctx context.Context, login string, id int) error
	Auth(ctx context.Context, login string) (*Auth, error)
	InternalAuth(ctx context.Context, serviceName string) (*InternalAuth, error)
	AdminAuth(ctx context.Context, login string) (*Auth, error)
}

type repository struct {
	db *sql.DB
}

func (r repository) Create(ctx context.Context, request queue.AuthCreate) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO auth(user_id, email, login, password) VALUES(0, $1, $2, $3)",
		request.Email,
		request.Login,
		request.Password,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r repository) UpdateUserID(ctx context.Context, login string, id int) error {
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE auth SET user_id = $1 WHERE login = $2",
		id,
		login,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r repository) Auth(ctx context.Context, login string) (*Auth, error) {
	response := &Auth{}

	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, user_id, email, login, password, created_at from auth WHERE login = $1",
		login,
	).Scan(&response.ID, &response.UserID, &response.Email, &response.Login, &response.Password, &response.CreatedAt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

func (r repository) InternalAuth(ctx context.Context, serviceName string) (*InternalAuth, error) {
	response := &InternalAuth{}

	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, service_name, password, created_at from internal_auth WHERE service_name = $1",
		serviceName,
	).Scan(&response.ID, &response.ServiceName, &response.Password, &response.CreatedAt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

func (r repository) AdminAuth(ctx context.Context, login string) (*Auth, error) {
	response := &Auth{}

	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, user_id, email, login, password, created_at from auth WHERE login = $1 AND role='admin'",
		login,
	).Scan(&response.ID, &response.UserID, &response.Email, &response.Login, &response.Password, &response.CreatedAt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

func NewRepository(db *sql.DB) Repository {
	return repository{db: db}
}
