package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/slirx/newproj/pkg/rabbitmq"
)

// Config represents combined configuration.
type Config struct {
	RabbitMQ rabbitmq.Config
	Database Database
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// NewConfig returns initialized instance of configuration. It reads configuration from environment variables.
// prefix represents prefix of environment variables' names.
func NewConfig(prefix string) (*Config, error) {
	var err error

	var rabbitmqMaxReconnections int
	if rabbitmqMaxReconnections, err = strconv.Atoi(os.Getenv(prefix + "RABBITMQ_MAX_RECONNECTIONS")); err != nil {
		return nil, errors.WithStack(fmt.Errorf("invalid %s value", prefix+"RABBITMQ_MAX_RECONNECTIONS"))
	}

	var rabbitmqReconnectTimeoutSeconds int

	rabbitmqReconnectTimeoutSeconds, err = strconv.Atoi(os.Getenv(prefix + "RABBITMQ_RECONNECT_TIMEOUT_SECONDS"))
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("invalid %s value", prefix+"RABBITMQ_RECONNECT_TIMEOUT_SECONDS"))
	}

	var databasePort int
	if databasePort, err = strconv.Atoi(os.Getenv(prefix + "DB_PORT")); err != nil {
		return nil, errors.WithStack(fmt.Errorf("invalid %s value", prefix+"DB_PORT"))
	}

	config := Config{
		RabbitMQ: rabbitmq.Config{
			URI:                     os.Getenv(prefix + "RABBITMQ_AMQP_URI"),
			MaxReconnections:        rabbitmqMaxReconnections,
			ReconnectTimeoutSeconds: time.Duration(int64(rabbitmqReconnectTimeoutSeconds)) * time.Second,
		},
		Database: Database{
			Host:     os.Getenv(prefix + "DB_HOST"),
			Port:     databasePort,
			User:     os.Getenv(prefix + "DB_USER"),
			Password: os.Getenv(prefix + "DB_PASSWORD"),
			Name:     os.Getenv(prefix + "DB_NAME"),
		},
	}

	return &config, nil
}
