package post

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"gitlab.com/slirx/newproj/pkg/queue"
)

type Repository interface {
	List(ctx context.Context, request ListRequest, perPage int) (*ListResponse, error)
	Create(ctx context.Context, uid int, users []int, request CreateRequest) (*CreateResponse, error)
	Feed(ctx context.Context, uid int, request FeedRequest, perPage int) (*FeedResponse, error)
	Unfollow(ctx context.Context, task queue.PostUnfollow) error
	Follow(ctx context.Context, task queue.PostFollow) error
	Search(ctx context.Context, request SearchRequest, perPage int) ([]int, error)
	Posts(ctx context.Context, postsIDs []int) ([]Post, error)
}

type repository struct {
	db *sql.DB
}

func (r repository) List(ctx context.Context, request ListRequest, perPage int) (*ListResponse, error) {
	if request.LatestPostID == 0 {
		request.LatestPostID = math.MaxInt64
	}

	response := ListResponse{ // todo use sync.pool?
		Posts: make([]Post, 0, perPage),
	}

	err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM post WHERE user_id = $1`, request.UserID).
		Scan(&response.Total)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, text, extract(epoch from created_at)::INT AS created_at FROM "post" 
				WHERE user_id = $1 AND id < $2
				ORDER BY id DESC
				LIMIT $3`,
		request.UserID,
		request.LatestPostID,
		perPage,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer rows.Close()

	post := Post{} // todo use sync.pool?
	for rows.Next() {
		err = rows.Scan(&post.ID, &post.Text, &post.CreatedAt)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		post.User.ID = request.UserID

		response.Posts = append(response.Posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	return &response, nil
}

func (r repository) Create(ctx context.Context, uid int, users []int, request CreateRequest) (*CreateResponse, error) {
	var id int

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer tx.Rollback()

	err = tx.QueryRowContext(
		ctx,
		`INSERT INTO "post" (user_id, text) VALUES($1, $2) RETURNING id`,
		uid,
		request.Text,
	).Scan(&id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if id == 0 {
		return nil, errors.WithStack(fmt.Errorf("invalid post id (0)"))
	}

	var feedStmt *sql.Stmt

	feedStmt, err = tx.PrepareContext(ctx, `INSERT INTO feed(user_id, post_id) VALUES ($1, $2)`)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = feedStmt.ExecContext(ctx, uid, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, userID := range users {
		_, err = feedStmt.ExecContext(ctx, userID, id)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	response := CreateResponse{}

	err = tx.QueryRowContext(
		ctx,
		`SELECT id, text, extract(epoch from created_at)::INT AS created_at FROM "post" WHERE id = $1`,
		id,
	).Scan(&response.ID, &response.Text, &response.CreatedAt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	response.User.ID = uid

	return &response, nil
}

func (r repository) Feed(ctx context.Context, uid int, request FeedRequest, perPage int) (*FeedResponse, error) {
	if request.LatestPostID == 0 {
		request.LatestPostID = math.MaxInt64
	}

	response := FeedResponse{ // todo use sync.pool?
		Posts: make([]Post, 0, perPage),
	}

	err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM feed WHERE user_id = $1`, uid).Scan(&response.Total)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT post.id, post.text, extract(epoch from post.created_at)::INT AS created_at, post.user_id FROM feed 
				JOIN post ON post.id = feed.post_id
				WHERE feed.user_id = $1 AND post.id < $2
				ORDER BY post.id DESC
				LIMIT $3`,
		uid,
		request.LatestPostID,
		perPage,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer rows.Close()

	post := Post{} // todo use sync.pool?
	for rows.Next() {
		err = rows.Scan(&post.ID, &post.Text, &post.CreatedAt, &post.User.ID)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		response.Posts = append(response.Posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	return &response, nil
}

func (r repository) Unfollow(ctx context.Context, task queue.PostUnfollow) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM feed USING post WHERE post.id = feed.post_id AND post.user_id = $1 AND feed.user_id = $2",
		task.UnfollowUserID,
		task.UserID,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r repository) Follow(ctx context.Context, task queue.PostFollow) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO feed (user_id, post_id) (SELECT $1, id from post where user_id = $2)",
		task.UserID,
		task.FollowUserID,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r repository) Search(ctx context.Context, request SearchRequest, perPage int) ([]int, error) {
	response := make([]int, 0)

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id
				FROM post, to_tsquery($1) q
				WHERE q @@ searchable_text
				ORDER BY ts_rank(searchable_text, q) DESC, created_at DESC
				LIMIT $2`,
		request.Query,
		perPage,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer rows.Close()

	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		response = append(response, id)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

func (r repository) Posts(ctx context.Context, postsIDs []int) ([]Post, error) {
	response := make([]Post, 0)

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT post.id, post.text, extract(epoch from post.created_at)::INT AS created_at, post.user_id FROM post
				JOIN unnest($1::int[]) WITH ORDINALITY t(id, ord) USING (id)
				ORDER BY t.ord`,
		pq.Array(postsIDs),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer rows.Close()

	post := Post{} // todo use sync.pool?
	for rows.Next() {
		err = rows.Scan(&post.ID, &post.Text, &post.CreatedAt, &post.User.ID)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		response = append(response, post)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

func NewRepository(db *sql.DB) Repository {
	return repository{db: db}
}
