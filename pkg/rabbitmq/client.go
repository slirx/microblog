package rabbitmq

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"gitlab.com/slirx/newproj/pkg/logger"
)

type Config struct {
	URI                     string
	MaxReconnections        int
	ReconnectTimeoutSeconds time.Duration
}

type Client struct {
	Config     Config
	Connection *connection
	Logger     logger.Logger
	QueueName  string
}

func (c Client) Messages() (<-chan amqp.Delivery, error) {
	q, err := c.Connection.Channel.QueueDeclare(
		c.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("can not declare queue: %w", err))
	}

	err = c.Connection.Channel.Qos(1, 0, false)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("can not apply Qos: %w", err))
	}

	messages, err := c.Connection.Channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("can not consume messages: %w", err))
	}

	return messages, nil
}

func (c Client) Events(exchangeName string) (<-chan amqp.Delivery, error) {
	err := c.Connection.Channel.ExchangeDeclare(
		exchangeName,
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	q, err := c.Connection.Channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	/*listenedServices := []string{"registration"}
	for _, service := range listenedServices {
		// todo which services listen to
		err = ch.QueueBind(
			q.Name,      // queue name
			service,     // routing key
			"microblog", // exchange
			false,
			nil,
		)
		if err != nil {
			panic(err)
		}
	}*/

	err = c.Connection.Channel.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	//c.Connection.Channel.Confirm()
	messages, err := c.Connection.Channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("can not consume messages: %w", err))
	}

	return messages, nil
}

func NewClient(conf Config, l logger.Logger, queueName string) *Client {
	conn := NewConnection(conf, l)
	return &Client{Connection: conn, Logger: l, QueueName: queueName, Config: conf}
}
