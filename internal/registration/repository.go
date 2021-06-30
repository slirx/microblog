package registration

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type Repository interface {
	Register(ctx context.Context, request RegisterRequest, code int) error
	ConfirmationData(ctx context.Context, email string) (*ConfirmationData, error)
	Confirm(ctx context.Context, email string) error
}

type repository struct {
	db *sql.DB
}

func (r repository) Register(ctx context.Context, request RegisterRequest, code int) error {
	// todo handle case when email or login is already in use
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO registration(email, login, code) VALUES($1, $2, $3)",
		request.Email,
		request.Login,
		code,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r repository) ConfirmationData(ctx context.Context, email string) (*ConfirmationData, error) {
	response := ConfirmationData{}

	err := r.db.QueryRowContext(
		ctx,
		"SELECT code, created_at, login FROM registration WHERE email = $1 LIMIT 1",
		email,
	).Scan(&response.Code, &response.CreatedAt, &response.Login)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &response, nil
}

func (r repository) Confirm(ctx context.Context, email string) error {
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE registration SET code=$1, confirmed_at=CURRENT_TIMESTAMP WHERE email=$2",
		0,
		email,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func NewRepository(db *sql.DB) Repository {
	return repository{db: db}
}
