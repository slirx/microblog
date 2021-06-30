package post

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmzap"

	"gitlab.com/slirx/newproj/pkg/logger"
	"gitlab.com/slirx/newproj/pkg/queue"
	"gitlab.com/slirx/newproj/pkg/queue/worker"
)

type followHandler struct {
	Logger     logger.Logger
	Repository Repository
}

type unfollowHandler struct {
	Logger     logger.Logger
	Repository Repository
}

func (h followHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
	task := queue.PostFollow{}

	if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&task); err != nil {
		return errors.WithStack(err)
	}

	tx := apm.TransactionFromContext(ctx)
	tx.Context.SetLabel("request_id", task.RequestID)

	body, err := json.Marshal(task)
	if err != nil {
		return errors.WithStack(err)
	}

	tx.Context.SetCustom("request_body", string(body))

	h.Logger.Debug(
		fmt.Sprintf("following (generating feed) user %d/%d", task.UserID, task.FollowUserID),
		apmzap.TraceContext(ctx)...,
	)

	err = h.Repository.Follow(ctx, task)
	if err != nil {
		return err
	}

	tx.Result = "success"
	tx.Outcome = "success"

	return nil
}

func (h unfollowHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
	task := queue.PostUnfollow{}

	if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&task); err != nil {
		return errors.WithStack(err)
	}

	tx := apm.TransactionFromContext(ctx)
	tx.Context.SetLabel("request_id", task.RequestID)

	body, err := json.Marshal(task)
	if err != nil {
		return errors.WithStack(err)
	}

	tx.Context.SetCustom("request_body", string(body))

	h.Logger.Debug(
		fmt.Sprintf("unfollowing (deleting feed) user %d/%d", task.UserID, task.UnfollowUserID),
		apmzap.TraceContext(ctx)...,
	)

	err = h.Repository.Unfollow(ctx, task)
	if err != nil {
		return err
	}

	tx.Result = "success"
	tx.Outcome = "success"

	return nil
}

func NewFollowHandler(l logger.Logger, repository Repository) worker.Handler {
	return followHandler{
		Logger:     l,
		Repository: repository,
	}
}

func NewUnfollowHandler(l logger.Logger, repository Repository) worker.Handler {
	return unfollowHandler{
		Logger:     l,
		Repository: repository,
	}
}
