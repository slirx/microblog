package worker

import (
	"context"

	"github.com/streadway/amqp"
)

type Worker interface {
	Run(ctx context.Context)
	EventListener(ctx context.Context, exchangeName string) // todo move to another interface?
}

type Handler interface {
	Handle(ctx context.Context, msg amqp.Delivery) error
}
