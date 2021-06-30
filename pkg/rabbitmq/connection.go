package rabbitmq

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.elastic.co/apm/module/apmzap"

	"gitlab.com/slirx/newproj/pkg/logger"
)

type connection struct {
	//URI                string
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Config     Config
	//MaxReconnections   int
	reconnectionsCount int
	Logger             logger.Logger
}

func (c *connection) Close() error {
	var err error

	if c.Channel != nil {
		if err = c.Channel.Close(); err != nil {
			return errors.WithStack(err)
		}
	}

	if c.Connection != nil {
		if err = c.Connection.Close(); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (c *connection) Reconnect(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errors.WithStack(ctx.Err())
	case <-time.After(c.Config.ReconnectTimeoutSeconds):
	}

	c.reconnectionsCount++

	if c.reconnectionsCount > c.Config.MaxReconnections {
		return errors.WithStack(fmt.Errorf("maximum reconnections reached"))
	}

	c.Logger.Debug("reconnecting.. try #"+strconv.Itoa(c.reconnectionsCount), apmzap.TraceContext(ctx)...)

	var err error
	if c.Connection, err = amqp.Dial(c.Config.URI); err != nil {
		return c.Reconnect(ctx)
	}

	c.Channel, err = c.Connection.Channel()
	if err != nil {
		return errors.WithStack(err)
	}

	c.Logger.Debug("reconnected", apmzap.TraceContext(ctx)...)

	return nil
}

func NewConnection(conf Config, l logger.Logger) *connection {
	return &connection{Config: conf, Logger: l}
}
