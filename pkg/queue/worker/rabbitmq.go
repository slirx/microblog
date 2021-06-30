package worker

import (
	"context"
	"time"

	"github.com/streadway/amqp"
	"gitlab.com/slirx/newproj/pkg/logger"
	"gitlab.com/slirx/newproj/pkg/rabbitmq"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmzap"
)

func NewWorker(
	workerName string,
	client *rabbitmq.Client,
	l logger.Logger,
	tracer *apm.Tracer,
	handler Handler,
	queueName string,
) Worker {
	return &rabbitmqWorker{
		Name:      workerName,
		Client:    client,
		Error:     make(chan error),
		Logger:    l,
		Handler:   handler,
		Tracer:    tracer,
		QueueName: queueName,
	}
}

type rabbitmqWorker struct {
	Name      string // worker name
	Client    *rabbitmq.Client
	Error     chan error
	Logger    logger.Logger
	Tracer    *apm.Tracer
	Handler   Handler
	QueueName string
}

func (w *rabbitmqWorker) Run(ctx context.Context) {
	done := make(chan struct{})

	go func() {
		for err := range w.Error {
			w.Logger.Error(err)
		}

		done <- struct{}{}
	}()

	defer func() {
		close(w.Error)
		<-done
	}()

	err := w.Client.Connection.Reconnect(ctx)
	if err != nil {
		w.Error <- err
		return
	}

	defer func() {
		w.Error <- w.Client.Connection.Close()
	}()

	for {
		messages, err := w.Client.Messages()
		if err != nil {
			w.Error <- err
			break
		}

		go func(ctx context.Context) {
			for d := range messages {
				select {
				case <-ctx.Done():
					return
				default:
				}

				w.handleMessage(ctx, d)
			}
		}(ctx)

		amqpErr := make(chan *amqp.Error)

		select {
		case <-w.Client.Connection.Connection.NotifyClose(amqpErr):
			if err = w.Client.Connection.Reconnect(ctx); err != nil {
				w.Error <- err
				return
			}

			continue
		case <-ctx.Done():
			return
		}
	}
}

func (w rabbitmqWorker) EventListener(ctx context.Context, exchangeName string) {
	done := make(chan struct{})

	go func() {
		for err := range w.Error {
			w.Logger.Error(err)
		}

		done <- struct{}{}
	}()

	defer func() {
		close(w.Error)
		<-done
	}()

	err := w.Client.Connection.Reconnect(ctx)
	if err != nil {
		w.Error <- err
		return
	}

	defer func() {
		w.Error <- w.Client.Connection.Close()
	}()

	for {
		events, err := w.Client.Events(exchangeName)
		if err != nil {
			w.Error <- err
			break
		}

		go func(ctx context.Context) {
			for d := range events {
				select {
				case <-ctx.Done():
					return
				default:
				}

				w.handleMessage(ctx, d)
			}
		}(ctx)

		amqpErr := make(chan *amqp.Error)

		select {
		case <-w.Client.Connection.Connection.NotifyClose(amqpErr):
			if err = w.Client.Connection.Reconnect(ctx); err != nil {
				w.Error <- err
				return
			}

			continue
		case <-ctx.Done():
			return
		}
	}
}

func (w *rabbitmqWorker) handleMessage(ctx context.Context, msg amqp.Delivery) {
	var err error

	opts := apm.TransactionOptions{
		Start:        time.Now(),
		TraceContext: apm.TraceContext{},
	}
	tx := w.Tracer.StartTransactionOptions(w.Name, "rabbitmq_task", opts)
	tx.Context.SetLabel("rabbitmq_queue", w.QueueName)
	defer tx.End()

	ctx = apm.ContextWithTransaction(ctx, tx)

	if err = w.Handler.Handle(ctx, msg); err != nil {
		// todo check that error is bounded with context. so in log goes trace.id
		w.Logger.Error(err, apmzap.TraceContext(ctx)...)

		if err = msg.Nack(false, true); err != nil {
			w.Logger.Error(err, apmzap.TraceContext(ctx)...)
		}

		return
	}

	if err = msg.Ack(false); err != nil {
		w.Logger.Error(err, apmzap.TraceContext(ctx)...)
	}
}
