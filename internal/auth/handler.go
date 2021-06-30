package auth

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
	"gitlab.com/slirx/newproj/pkg/queue/manager"
	"gitlab.com/slirx/newproj/pkg/queue/worker"
)

type CreateHandler struct {
	Logger     logger.Logger
	Repository Repository
	Manager    manager.Manager
}

type UpdateUserIDHandler struct {
	Logger     logger.Logger
	Repository Repository
}

func (h CreateHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
	task := queue.AuthCreate{}

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

	// todo do your dirty work here..
	h.Logger.Debug(fmt.Sprintf("creating auth for %s", task.Login), apmzap.TraceContext(ctx)...)

	err = h.Repository.Create(ctx, task)
	if err != nil {
		return err
	}

	err = h.Manager.Send(ctx, queue.JobUserCreate, queue.UserCreate{
		RequestID: task.RequestID,
		Login:     task.Login,
		Email:     task.Email,
	})
	if err != nil {
		return err
	}

	tx.Result = "success"
	tx.Outcome = "success"

	return nil
}

func (h UpdateUserIDHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
	task := queue.AuthUpdateUserID{}

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

	// fixme this happens before creating auth record..

	// todo do your dirty work here..
	h.Logger.Debug(fmt.Sprintf("updating user id for %s: %d", task.Login, task.UserID), apmzap.TraceContext(ctx)...)

	err = h.Repository.UpdateUserID(ctx, task.Login, task.UserID)
	if err != nil {
		return err
	}

	tx.Result = "success"
	tx.Outcome = "success"

	return nil
}

func NewCreateHandler(l logger.Logger, repository Repository, m manager.Manager) worker.Handler {
	return CreateHandler{
		Logger:     l,
		Repository: repository,
		Manager:    m,
	}
}

func NewUpdateUserIDHandler(l logger.Logger, repository Repository) worker.Handler {
	return UpdateUserIDHandler{
		Logger:     l,
		Repository: repository,
	}
}
