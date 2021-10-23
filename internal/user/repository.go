package user

import (
	"context"
	"database/sql"
	"math"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"gitlab.com/slirx/newproj/pkg/queue"
)

type Repository interface {
	Create(ctx context.Context, request queue.UserCreate) (int, error)
	Update(ctx context.Context, uid int, request UpdateRequest) error
	Get(ctx context.Context, login string, uid int) (*GetResponse, error)
	Me(ctx context.Context, uid int) (*MeResponse, error)
	Follow(ctx context.Context, uid int, request FollowRequest) error
	Unfollow(ctx context.Context, uid int, request UnfollowRequest) error
	Followers(ctx context.Context, request FollowersRequest, perPage uint8) (*FollowersResponse, error)
	Following(ctx context.Context, request FollowingRequest, perPage uint8) (*FollowingResponse, error)
	FollowersIDs(ctx context.Context, uid int) ([]int, error)
	Users(ctx context.Context, request InternalUsersRequest, perPage int) (*InternalUsersResponse, error)
}

type repository struct {
	db *sql.DB
}

func (r repository) Create(ctx context.Context, request queue.UserCreate) (int, error) {
	var id int

	// todo add name column
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO "user" (email, login) VALUES($1, $2) RETURNING id`,
		request.Email,
		request.Login,
	).Scan(&id)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return id, nil
}

func (r repository) Update(ctx context.Context, uid int, request UpdateRequest) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE "user" SET bio = $1, name = $2 WHERE id = $3`,
		request.Bio,
		request.Name,
		uid,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r repository) Get(ctx context.Context, login string, uid int) (*GetResponse, error) {
	response := GetResponse{}

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, login, name, bio, followers, following FROM "user" WHERE login = $1`,
		login,
	).Scan(&response.ID, &response.Login, &response.Name, &response.Bio, &response.Followers, &response.Following)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if uid > 0 {
		err = r.db.QueryRowContext(
			ctx,
			"SELECT EXISTS (SELECT id FROM follower WHERE user_id = $1 AND follower_id = $2)",
			response.ID,
			uid,
		).Scan(&response.IsFollowed)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &response, nil
}

func (r repository) Me(ctx context.Context, uid int) (*MeResponse, error) {
	response := MeResponse{}

	err := r.db.QueryRowContext(ctx, `SELECT id, login, name, bio, followers, following FROM "user" WHERE id = $1`, uid).
		Scan(&response.ID, &response.Login, &response.Name, &response.Bio, &response.Followers, &response.Following)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &response, nil
}

func (r repository) Follow(ctx context.Context, uid int, request FollowRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO "follower" (user_id, follower_id) VALUES($1, $2)`,
		request.UserID,
		uid,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE "user" SET followers = followers + 1 WHERE id = $1`,
		request.UserID,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE "user" SET following = following + 1 WHERE id = $1`,
		uid,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = tx.Commit(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r repository) Unfollow(ctx context.Context, uid int, request UnfollowRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM "follower" where user_id = $1 AND follower_id = $2`,
		request.UserID,
		uid,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE "user" SET followers = followers - 1 WHERE id = $1`,
		request.UserID,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE "user" SET following = following - 1 WHERE id = $1`,
		uid,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = tx.Commit(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r repository) Followers(ctx context.Context, request FollowersRequest, perPage uint8) (*FollowersResponse, error) {
	if request.LatestFollowerID == 0 {
		request.LatestFollowerID = math.MaxInt64
	}

	response := FollowersResponse{
		Total:     0,
		Followers: make([]Follower, 0, perPage),
	}

	userID := 0
	err := r.db.QueryRowContext(ctx, `SELECT id FROM "user" WHERE login = $1`, request.Login).Scan(&userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM follower WHERE user_id = $1`, userID).Scan(&response.Total)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if response.Total == 0 {
		return &response, nil
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT follower.id AS follower_id, "user".id AS user_id, "user".login, "user".name, "user".bio FROM "user" 
			JOIN follower ON follower.follower_id = "user".id WHERE follower.user_id = $1
			AND follower.id < $2
			ORDER BY follower.id DESC
			LIMIT $3`,
		userID,
		request.LatestFollowerID,
		perPage,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer rows.Close()

	var f Follower
	for rows.Next() {
		if err = rows.Scan(&f.FollowerID, &f.UserID, &f.Login, &f.Name, &f.Bio); err != nil {
			return nil, errors.WithStack(err)
		}

		response.Followers = append(response.Followers, f)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	return &response, nil
}

func (r repository) Following(ctx context.Context, request FollowingRequest, perPage uint8) (*FollowingResponse, error) {
	if request.LatestFollowerID == 0 {
		request.LatestFollowerID = math.MaxInt64
	}

	response := FollowingResponse{
		Total:     0,
		Following: make([]Followed, 0, perPage),
	}

	userID := 0
	err := r.db.QueryRowContext(ctx, `SELECT id FROM "user" WHERE login = $1`, request.Login).Scan(&userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM follower WHERE follower_id = $1`, userID).
		Scan(&response.Total)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if response.Total == 0 {
		return &response, nil
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT follower.id AS follower_id, "user".id AS user_id, "user".login, "user".name, "user".bio FROM "user" 
			JOIN follower ON follower.user_id = "user".id WHERE follower.follower_id = $1
			AND follower.id < $2
			ORDER BY follower.id DESC
			LIMIT $3`,
		userID,
		request.LatestFollowerID,
		perPage,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer rows.Close()

	var f Followed
	for rows.Next() {
		if err = rows.Scan(&f.FollowerID, &f.UserID, &f.Login, &f.Name, &f.Bio); err != nil {
			return nil, errors.WithStack(err)
		}

		response.Following = append(response.Following, f)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	return &response, nil
}

func (r repository) FollowersIDs(ctx context.Context, uid int) ([]int, error) {
	response := make([]int, 0)
	var id int

	rows, err := r.db.QueryContext(ctx, "SELECT follower_id from follower where user_id = $1", uid)
	if err != nil {
		return response, errors.WithStack(err)
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return response, errors.WithStack(err)
		}

		response = append(response, id)
	}

	if err = rows.Err(); err != nil {
		return response, errors.WithStack(err)
	}

	return response, nil
}

func (r repository) Users(ctx context.Context, request InternalUsersRequest, perPage int) (*InternalUsersResponse, error) {
	response := &InternalUsersResponse{
		Total: 0,
		Users: make([]User, 0),
	}

	var rows *sql.Rows
	var err error

	if len(request.UserIDs) != 0 {
		// fetch users by IDs
		rows, err = r.db.QueryContext(
			ctx,
			`SELECT id, name, login FROM "user" WHERE id = ANY($1)`,
			pq.Array(request.UserIDs),
		)
		if err != nil {
			return response, errors.WithStack(err)
		}
	} else {
		err = r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM "user"`).Scan(&response.Total)
		if err != nil {
			return response, errors.WithStack(err)
		}

		// fetch all users (with pagination)
		rows, err = r.db.QueryContext(
			ctx,
			`SELECT id, name, login FROM "user" WHERE id < $1 ORDER BY id DESC LIMIT $2`,
			request.LatestUserID,
			perPage,
		)
		if err != nil {
			return response, errors.WithStack(err)
		}
	}

	defer rows.Close()

	var u User
	for rows.Next() {
		if err = rows.Scan(&u.ID, &u.Name, &u.Login); err != nil {
			return response, errors.WithStack(err)
		}

		response.Users = append(response.Users, u)
	}

	if err = rows.Err(); err != nil {
		return response, errors.WithStack(err)
	}

	return response, nil
}

func NewRepository(db *sql.DB) Repository {
	return repository{db: db}
}
