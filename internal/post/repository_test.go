package post

import (
	"context"
	"database/sql"
	"log"
	"math"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"

	"gitlab.com/slirx/newproj/pkg/queue"
)

func newDatabaseMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestRepositoryListSuccess(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	wantTotal := 2
	userID := 1

	u := PostsUser{
		ID: 1,
	}

	posts := []Post{
		{
			ID:        2,
			Text:      "post 2",
			CreatedAt: 10050015,
			User:      u,
		},
		{
			ID:        1,
			Text:      "post 1",
			CreatedAt: 10050012,
			User:      u,
		},
	}

	rows := sqlmock.NewRows([]string{"COUNT(1)"}).AddRow(wantTotal)
	exec := regexp.QuoteMeta(`SELECT COUNT(1) FROM post WHERE user_id = $1`)
	mock.ExpectQuery(exec).WithArgs(userID).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id", "text", "created_at"}).
		AddRow(posts[0].ID, posts[0].Text, posts[0].CreatedAt).
		AddRow(posts[1].ID, posts[1].Text, posts[1].CreatedAt)
	exec = regexp.QuoteMeta(`SELECT id, text, extract(epoch from created_at)::INT AS created_at FROM "post"
				WHERE user_id = $1 AND id < $2
				ORDER BY id DESC
				LIMIT $3`)
	mock.ExpectQuery(exec).WithArgs(userID, math.MaxInt64, perPage).WillReturnRows(rows)

	ctx := context.Background()

	request := ListRequest{
		UserID:       userID,
		LatestPostID: 0,
	}

	response, err := repo.List(ctx, request, perPage)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}

	if wantTotal != response.Total {
		t.Fatalf("got: %d, want: %d", response.Total, wantTotal)
	}

	if wantTotal != len(response.Posts) {
		t.Fatalf("got: %d, want: %d", len(response.Posts), wantTotal)
	}

	if posts[0].ID != response.Posts[0].ID {
		t.Fatalf("got: %d, want: %d", response.Posts[0].ID, posts[0].ID)
	}
}

func TestRepositoryCreateSuccess(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	userID := 1
	users := []int{2, 3}
	postID := 1
	text := "post text #1"
	createdAt := 10050012

	mock.ExpectBegin()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(postID)

	exec := regexp.QuoteMeta(
		`INSERT INTO "post" (user_id, text) VALUES($1, $2) RETURNING id`,
	)
	mock.ExpectQuery(exec).
		WithArgs(userID, text).
		WillReturnRows(rows)

	stmt := mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO feed(user_id, post_id) VALUES ($1, $2)`))
	stmt.ExpectExec().
		WithArgs(userID, postID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	for _, id := range users {
		stmt.ExpectExec().
			WithArgs(id, postID).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	rows = sqlmock.NewRows([]string{"id", "text", "created_at"}).AddRow(postID, text, createdAt)
	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT id, text, extract(epoch from created_at)::INT AS created_at FROM "post" WHERE id = $1`,
		),
	).WithArgs(postID).WillReturnRows(rows)

	mock.ExpectCommit()

	ctx := context.Background()
	request := CreateRequest{
		Text: text,
	}

	response, err := repo.Create(ctx, userID, users, request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}

	if postID != response.ID {
		t.Fatalf("got: %d, want: %d", response.ID, postID)
	}

	if text != response.Text {
		t.Fatalf("got: %s, want: %s", response.Text, text)
	}
}

func TestRepositoryFeedSuccess(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	wantTotal := 4
	userID := 1

	u1 := PostsUser{
		ID: 1,
	}
	u2 := PostsUser{
		ID: 2,
	}

	posts := []Post{
		{
			ID:        4,
			Text:      "post 4",
			CreatedAt: 10050021,
			User:      u1,
		},
		{
			ID:        3,
			Text:      "post 3",
			CreatedAt: 10050018,
			User:      u1,
		},
		{
			ID:        2,
			Text:      "post 2",
			CreatedAt: 10050015,
			User:      u2,
		},
		{
			ID:        1,
			Text:      "post 1",
			CreatedAt: 10050012,
			User:      u1,
		},
	}

	rows := sqlmock.NewRows([]string{"COUNT(1)"}).AddRow(wantTotal)
	exec := regexp.QuoteMeta(`SELECT COUNT(1) FROM feed WHERE user_id = $1`)
	mock.ExpectQuery(exec).WithArgs(userID).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id", "text", "created_at", "user_id"}).
		AddRow(posts[0].ID, posts[0].Text, posts[0].CreatedAt, posts[0].User.ID).
		AddRow(posts[1].ID, posts[1].Text, posts[1].CreatedAt, posts[1].User.ID).
		AddRow(posts[2].ID, posts[2].Text, posts[2].CreatedAt, posts[2].User.ID).
		AddRow(posts[3].ID, posts[3].Text, posts[3].CreatedAt, posts[3].User.ID)
	exec = regexp.QuoteMeta(
		`SELECT post.id, post.text, extract(epoch from post.created_at)::INT AS created_at, post.user_id FROM feed 
			JOIN post ON post.id = feed.post_id
			WHERE feed.user_id = $1 AND post.id < $2
			ORDER BY post.id DESC
			LIMIT $3`)
	mock.ExpectQuery(exec).WithArgs(userID, math.MaxInt64, perPage).WillReturnRows(rows)

	ctx := context.Background()

	request := FeedRequest{
		LatestPostID: 0,
	}

	response, err := repo.Feed(ctx, userID, request, perPage)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}

	if wantTotal != response.Total {
		t.Fatalf("got: %d, want: %d", response.Total, wantTotal)
	}

	if wantTotal != len(response.Posts) {
		t.Fatalf("got: %d, want: %d", len(response.Posts), wantTotal)
	}

	if posts[0].ID != response.Posts[0].ID {
		t.Fatalf("got: %d, want: %d", response.Posts[0].ID, posts[0].ID)
	}
}

func TestRepositoryUnfollowSuccess(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	userID := 1
	unfollowUserID := 2

	exec := regexp.QuoteMeta(
		`DELETE FROM feed USING post WHERE post.id = feed.post_id AND post.user_id = $1 AND feed.user_id = $2`,
	)
	mock.ExpectExec(exec).
		WithArgs(unfollowUserID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()

	request := queue.PostUnfollow{
		UserID:         userID,
		UnfollowUserID: unfollowUserID,
	}

	err := repo.Unfollow(ctx, request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}
}

func TestRepositoryFollowSuccess(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	userID := 1
	followUserID := 2

	exec := regexp.QuoteMeta(`INSERT INTO feed (user_id, post_id) (SELECT $1, id from post where user_id = $2)`)
	mock.ExpectExec(exec).
		WithArgs(userID, followUserID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()

	request := queue.PostFollow{
		UserID:       userID,
		FollowUserID: followUserID,
	}

	err := repo.Follow(ctx, request)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}
}

func TestRepositorySearchSuccess(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	wantTotal := 3

	posts := []int{5, 3, 2}
	request := SearchRequest{
		Query:   "some text",
		QueryID: 0,
		Offset:  0,
	}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(posts[0]).
		AddRow(posts[1]).
		AddRow(posts[2])
	exec := regexp.QuoteMeta(
		`SELECT id
				FROM post, to_tsquery($1) q
				WHERE q @@ searchable_text
				ORDER BY ts_rank(searchable_text, q) DESC, created_at DESC
				LIMIT $2`)
	mock.ExpectQuery(exec).WithArgs(request.Query, maxSearchPerPage).WillReturnRows(rows)

	ctx := context.Background()

	response, err := repo.Search(ctx, request, maxSearchPerPage)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}

	if wantTotal != len(response) {
		t.Fatalf("got: %d, want: %d", len(response), wantTotal)
	}

	if posts[0] != response[0] {
		t.Fatalf("got: %d, want: %d", response[0], posts[0])
	}
}

func TestRepositoryPostsSuccess(t *testing.T) {
	db, mock := newDatabaseMock()
	repo := NewRepository(db)

	wantTotal := 4

	u1 := PostsUser{
		ID: 1,
	}
	u2 := PostsUser{
		ID: 2,
	}

	posts := []Post{
		{
			ID:        4,
			Text:      "post 4",
			CreatedAt: 10050021,
			User:      u1,
		},
		{
			ID:        3,
			Text:      "post 3",
			CreatedAt: 10050018,
			User:      u1,
		},
		{
			ID:        2,
			Text:      "post 2",
			CreatedAt: 10050015,
			User:      u2,
		},
		{
			ID:        1,
			Text:      "post 1",
			CreatedAt: 10050012,
			User:      u1,
		},
	}

	postsIDs := []int{4, 3, 2, 1}

	rows := sqlmock.NewRows([]string{"id", "text", "created_at", "user_id"}).
		AddRow(posts[0].ID, posts[0].Text, posts[0].CreatedAt, posts[0].User.ID).
		AddRow(posts[1].ID, posts[1].Text, posts[1].CreatedAt, posts[1].User.ID).
		AddRow(posts[2].ID, posts[2].Text, posts[2].CreatedAt, posts[2].User.ID).
		AddRow(posts[3].ID, posts[3].Text, posts[3].CreatedAt, posts[3].User.ID)
	exec := regexp.QuoteMeta(
		`SELECT post.id, post.text, extract(epoch from post.created_at)::INT AS created_at, post.user_id FROM post
				JOIN unnest($1::int[]) WITH ORDINALITY t(id, ord) USING (id)
				ORDER BY t.ord`)
	mock.ExpectQuery(exec).WithArgs(pq.Array(postsIDs)).WillReturnRows(rows)

	ctx := context.Background()

	response, err := repo.Posts(ctx, postsIDs)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err.Error())
	}

	if wantTotal != len(response) {
		t.Fatalf("got: %d, want: %d", len(response), wantTotal)
	}

	if posts[0].ID != response[0].ID {
		t.Fatalf("got: %d, want: %d", response[0].ID, posts[0].ID)
	}

	if posts[3].ID != response[3].ID {
		t.Fatalf("got: %d, want: %d", response[3].ID, posts[3].ID)
	}
}
