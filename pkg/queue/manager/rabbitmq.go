package manager

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"gitlab.com/slirx/newproj/pkg/logger"
	"gitlab.com/slirx/newproj/pkg/rabbitmq"
	"gitlab.com/slirx/newproj/pkg/utils"
)

type rabbitmqManager struct {
	Client         *rabbitmq.Client
	Error          chan error
	Logger         logger.Logger
	isReconnecting chan struct{}
}

func (m *rabbitmqManager) Close() error {
	err := m.Client.Connection.Close()
	if err != nil {
		return err
	}

	close(m.Error)

	if !utils.IsStructChanClosed(m.isReconnecting) {
		close(m.isReconnecting)
	}

	return nil
}

func (m *rabbitmqManager) handleErrors() {
	for err := range m.Error {
		// todo refactor error handling?
		m.Logger.Error(err)
	}
}

func (m *rabbitmqManager) reconnect(ctx context.Context) {
	m.isReconnecting = make(chan struct{})

	if err := m.Client.Connection.Reconnect(ctx); err != nil {
		close(m.isReconnecting)
		m.Error <- errors.WithStack(err)
		return
	}

	close(m.isReconnecting)

	for {
		amqpErr := make(chan *amqp.Error)

		select {
		case <-m.Client.Connection.Connection.NotifyClose(amqpErr):
			m.isReconnecting = make(chan struct{})

			if err := m.Client.Connection.Reconnect(ctx); err != nil {
				m.Error <- errors.WithStack(err)
				return
			}

			close(m.isReconnecting)
		case <-ctx.Done():
			return
		}
	}
}

func (m *rabbitmqManager) Send(ctx context.Context, routingKey string, msg interface{}) error {
	// wait in case reconnection is in progress
	select {
	case <-m.isReconnecting:
	case <-time.After(m.Client.Config.ReconnectTimeoutSeconds * 2):
		return errors.WithStack(errors.New("queue is not responding"))
	case <-ctx.Done():
		return nil
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(msg)
	if err != nil {
		return errors.WithStack(err)
	}

	err = m.Client.Connection.Channel.Publish(
		"",         // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/octet-stream",
			Body:         buf.Bytes(),
		})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// todo maybe combine with Send method internally to use one method
func (m *rabbitmqManager) EmitEvent(ctx context.Context, exchange string, msg interface{}) error {
	// wait in case reconnection is in progress
	select {
	case <-m.isReconnecting:
	case <-time.After(m.Client.Config.ReconnectTimeoutSeconds * 2):
		return errors.WithStack(errors.New("queue is not responding"))
	case <-ctx.Done():
		return nil
	}

	err := m.Client.Connection.Channel.ExchangeDeclare(
		exchange, // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return errors.WithStack(err)
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err = encoder.Encode(msg); err != nil {
		return errors.WithStack(err)
	}

	// todo I can call Channel.NotifyReturn to be notified when event is not published

	err = m.Client.Connection.Channel.Publish(
		exchange, // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // todo do I need this for pub/sub?
			ContentType:  "application/octet-stream",
			Body:         buf.Bytes(),
		})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func NewManager(ctx context.Context, l logger.Logger, conf rabbitmq.Config) Manager {
	m := &rabbitmqManager{
		Client:         rabbitmq.NewClient(conf, l, ""),
		Error:          make(chan error),
		Logger:         l,
		isReconnecting: make(chan struct{}),
	}

	go m.reconnect(ctx)
	go m.handleErrors()

	return m
}
