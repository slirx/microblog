package user

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

func (h CreateHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
	task := queue.UserCreate{}

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

	var userID int

	userID, err = h.Repository.Create(ctx, task)
	if err != nil {
		return err
	}

	err = h.Manager.Send(ctx, queue.JobAuthUpdateUserIDAuth, queue.AuthUpdateUserID{
		RequestID: task.RequestID,
		Login:     task.Login,
		UserID:    userID,
	})
	if err != nil {
		return err
	}

	tx.Result = "success"
	tx.Outcome = "success"

	return nil
}

func NewCreateHandler(
	l logger.Logger,
	repository Repository,
	m manager.Manager,
) worker.Handler {
	return CreateHandler{
		Logger:     l,
		Repository: repository,
		Manager:    m,
	}
}
