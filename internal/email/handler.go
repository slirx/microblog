package email

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

	"gitlab.com/slirx/newproj/pkg/queue"

	"gitlab.com/slirx/newproj/pkg/logger"
	"gitlab.com/slirx/newproj/pkg/queue/worker"
)

type Handler struct {
	Logger logger.Logger
}

func (h Handler) Handle(ctx context.Context, msg amqp.Delivery) error {
	task := queue.Email{}

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
	h.Logger.Debug(fmt.Sprintf("sending email to %s; text: %s", task.RecipientEmail, task.Text), apmzap.TraceContext(ctx)...)

	tx.Result = "success"
	tx.Outcome = "success"

	return nil
}

func NewHandler(l logger.Logger) worker.Handler {
	return Handler{Logger: l}
}
