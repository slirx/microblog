package user

import (
	"context"
	"database/sql"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"

	"gitlab.com/slirx/newproj/pkg/queue"
)

func newDatabaseMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestRepositoryCreate(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	email := "test@test.com"
	login := "test"
	id := 12

	rows := sqlmock.NewRows([]string{"id"}).AddRow(id)

	exec := regexp.QuoteMeta(`INSERT INTO "user" (email, login) VALUES($1, $2) RETURNING id`)
	mock.ExpectQuery(exec).WithArgs(email, login).WillReturnRows(rows)

	ctx := context.Background()
	request := queue.UserCreate{
		RequestID: "",
		Login:     login,
		Email:     email,
	}

	uid, err := repo.Create(ctx, request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}

	if id != uid {
		t.Fatalf("got: %d, want: %d", uid, id)
	}
}

func TestRepositoryCreateError(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	email := "test@test.com"
	login := "test"
	wantErr := "create error"

	exec := regexp.QuoteMeta(`INSERT INTO "user" (email, login) VALUES($1, $2) RETURNING id`)
	mock.ExpectQuery(exec).WithArgs(email, login).WillReturnError(errors.New(wantErr))

	ctx := context.Background()
	request := queue.UserCreate{
		RequestID: "",
		Login:     login,
		Email:     email,
	}

	uid, err := repo.Create(ctx, request)
	if err == nil {
		t.Fatalf("got: nil, want: %s", wantErr)
	}

	if wantErr != err.Error() {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}

	if 0 != uid {
		t.Fatalf("got: %d, want: 0", uid)
	}
}

func TestRepositoryUpdate(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	var uid = 1
	name := "John"
	bio := "some text.."

	exec := regexp.QuoteMeta(`UPDATE "user" SET bio = $1, name = $2 WHERE id = $3`)
	mock.ExpectExec(exec).WithArgs(bio, name, 1).WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	request := UpdateRequest{
		Bio:  bio,
		Name: name,
	}

	err := repo.Update(ctx, uid, request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}
}

func TestRepositoryUpdateError(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	var uid = 1
	bio := "some text.."
	name := "John"
	wantErr := "update error"

	exec := regexp.QuoteMeta(`UPDATE "user" SET bio = $1, name = $2 WHERE id = $3`)
	mock.ExpectExec(exec).WithArgs(bio, name, 1).WillReturnError(errors.New(wantErr))

	ctx := context.Background()
	request := UpdateRequest{
		Bio:  bio,
		Name: name,
	}

	err := repo.Update(ctx, uid, request)
	if err == nil {
		t.Fatalf("got: nil, want: %s", wantErr)
	}

	if wantErr != err.Error() {
		t.Fatalf("got: %s, want: %s", err.Error(), wantErr)
	}
}
