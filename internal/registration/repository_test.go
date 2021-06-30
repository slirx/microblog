package registration

import (
	"context"
	"database/sql"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
)

func newDatabaseMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestRepositoryRegisterSuccess(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	email := "test@test.com"
	login := "test"
	code := 1234

	exec := regexp.QuoteMeta(`INSERT INTO registration(email, login, code) VALUES($1, $2, $3)`)
	mock.ExpectExec(exec).
		WithArgs(email, login, code).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	request := RegisterRequest{
		Login: login,
		Email: email,
	}

	err := repo.Register(ctx, request, code)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}
}

func TestRepositoryRegisterError(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	email := "test@test.com"
	login := "test"
	code := 1234
	wantErr := "insert into error"

	exec := regexp.QuoteMeta(`INSERT INTO registration(email, login, code) VALUES($1, $2, $3)`)
	mock.ExpectExec(exec).
		WithArgs(email, login, code).
		WillReturnError(errors.New(wantErr))

	ctx := context.Background()
	request := RegisterRequest{
		Login: login,
		Email: email,
	}

	err := repo.Register(ctx, request, code)
	if err == nil {
		t.Fatalf("got: nil, want: %s", wantErr)
	}

	if wantErr != err.Error() {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}
}

func TestRepositoryConfirmationDataSuccess(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	code := 1234
	createdAt := time.Now()
	login := "test"
	email := "test@test.com"

	rows := sqlmock.NewRows([]string{"code", "created_at", "login"}).AddRow(code, createdAt, login)

	exec := regexp.QuoteMeta(`SELECT code, created_at, login FROM registration WHERE email = $1 LIMIT 1`)
	mock.ExpectQuery(exec).WithArgs(email).WillReturnRows(rows)

	ctx := context.Background()

	response, err := repo.ConfirmationData(ctx, email)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}

	if code != response.Code {
		t.Fatalf("got: %d, want: %d", response.Code, code)
	}

	if createdAt != response.CreatedAt {
		t.Fatalf("got: %s, want: %s", response.CreatedAt, createdAt)
	}

	if login != response.Login {
		t.Fatalf("got: %s, want: %s", response.Login, login)
	}
}

func TestRepositoryConfirmationDataError(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	email := "test@test.com"
	wantErr := "select error"

	exec := regexp.QuoteMeta(`SELECT code, created_at, login FROM registration WHERE email = $1 LIMIT 1`)
	mock.ExpectQuery(exec).WithArgs(email).WillReturnError(errors.New(wantErr))

	ctx := context.Background()

	response, err := repo.ConfirmationData(ctx, email)
	if err == nil {
		t.Fatalf("got: nil, want: %s", wantErr)
	}

	if wantErr != err.Error() {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	if response != nil {
		t.Fatalf("got: %v, want: nil", response)
	}
}

func TestRepositoryConfirmSuccess(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	email := "test@test.com"
	code := 0

	exec := regexp.QuoteMeta(`UPDATE registration SET code=$1, confirmed_at=CURRENT_TIMESTAMP WHERE email=$2`)
	mock.ExpectExec(exec).
		WithArgs(code, email).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()

	err := repo.Confirm(ctx, email)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}
}

func TestRepositoryConfirmError(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	email := "test@test.com"
	code := 0
	wantErr := "update error"

	exec := regexp.QuoteMeta(`UPDATE registration SET code=$1, confirmed_at=CURRENT_TIMESTAMP WHERE email=$2`)
	mock.ExpectExec(exec).
		WithArgs(code, email).
		WillReturnError(errors.New(wantErr))

	ctx := context.Background()

	err := repo.Confirm(ctx, email)
	if err == nil {
		t.Fatalf("got: nil, want: %s", wantErr)
	}

	if wantErr != err.Error() {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}
}
