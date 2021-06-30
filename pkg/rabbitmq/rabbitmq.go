package rabbitmq

import (
	"bytes"
	"encoding/gob"

	"github.com/streadway/amqp"
)

// todo implement this interface
// EventManager represents interface for working with rabbitmq or other message broker.
type EventManager interface {
	// Job sends message to the specified service as a job. It's work queue.
	Job(service string, message interface{}) error
	// Event sends message to the specified routingKey as an event. It's pub/sub pattern.
	Event(routingKey string, message interface{}) error
}

type rabbitmq struct {
	//err chan error
}

func (r rabbitmq) Job(service string, message interface{}) error {
	//select { //non blocking channel - if there is no error will go to default where we do nothing
	//case err := <-r.err:
	//	if err != nil {
	//		r.conn.Reconnect()
	//	}
	//default:
	//}

	panic("implement me")
}

func (r rabbitmq) Event(routingKey string, message interface{}) error {
	panic("implement me")
}

// ---------------

type Publisher interface {
	Publish(service string, msg interface{}) error
}

type eventPublisher struct {
	Exchange string
	Channel  *amqp.Channel
}

func (p eventPublisher) Publish(service string, msg interface{}) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(msg)

	return p.Channel.Publish(
		p.Exchange, // exchange
		service,    // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        buf.Bytes(),
		})
}

type jobPublisher struct {
	channel    *amqp.Channel
	Connection *amqp.Connection
}

func (p jobPublisher) Publish(service string, msg interface{}) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(msg)

	err := p.channel.Publish(
		"",      // exchange
		service, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/octet-stream",
			Body:         buf.Bytes(),
		})

	// todo in case of error - recreate channel? check it

	return err
}

func NewEventPublisher(exchange string, ch *amqp.Channel) Publisher {
	return eventPublisher{Exchange: exchange, Channel: ch}
}

func NewJobPublisher(connection *amqp.Connection) Publisher {
	// todo how to close this channel?
	jobChannel, err := connection.Channel()
	if err != nil {
		panic(err)
	}

	return jobPublisher{channel: jobChannel, Connection: connection}
}
