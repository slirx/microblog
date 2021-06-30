package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"gitlab.com/slirx/newproj/internal/api"
	"gitlab.com/slirx/newproj/pkg/rabbitmq"
	"gitlab.com/slirx/newproj/pkg/redis"
)

// Config represents combined configuration.
type Config struct {
	Server               Server
	RabbitMQ             rabbitmq.Config
	Redis                redis.Config
	Database             Database
	InternalAPIConfig    api.ServiceConfig
	InternalAPIEndpoints map[string]string
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type JWT struct {
	Secret []byte
}

// Server represents web server configuration.
type Server struct {
	// Addr represents address and port which server should listen to. It's specified in format host:port.
	Addr string
	JWT  JWT
}

// NewConfig returns initialized instance of configuration. It reads configuration from environment variables.
// prefix represents prefix of environment variables' names.
func NewConfig(prefix string) (*Config, error) {
	var err error

	//var rabbitmqMaxReconnections int
	//if rabbitmqMaxReconnections, err = strconv.Atoi(os.Getenv(prefix + "RABBITMQ_MAX_RECONNECTIONS")); err != nil {
	//	return nil, errors.WithStack(fmt.Errorf("invalid %s value", prefix+"RABBITMQ_MAX_RECONNECTIONS"))
	//}

	//var rabbitmqReconnectTimeoutSeconds int

	//rabbitmqReconnectTimeoutSeconds, err = strconv.Atoi(os.Getenv(prefix + "RABBITMQ_RECONNECT_TIMEOUT_SECONDS"))
	//if err != nil {
	//	return nil, errors.WithStack(fmt.Errorf("invalid %s value", prefix+"RABBITMQ_RECONNECT_TIMEOUT_SECONDS"))
	//}

	var databasePort int
	if databasePort, err = strconv.Atoi(os.Getenv(prefix + "DB_PORT")); err != nil {
		return nil, errors.WithStack(fmt.Errorf("invalid %s value", prefix+"DB_PORT"))
	}

	var redisDB int
	if redisDB, err = strconv.Atoi(os.Getenv(prefix + "REDIS_DB")); err != nil {
		return nil, errors.WithStack(fmt.Errorf("invalid %s value", prefix+"REDIS_DB"))
	}

	endpoints := make(map[string]string)
	endpoints["auth"] = os.Getenv(prefix + "ENDPOINT_AUTH")
	endpoints["user"] = os.Getenv(prefix + "ENDPOINT_USER")

	config := Config{
		Server: Server{
			Addr: os.Getenv(prefix + "SERVER_ADDR"),
			JWT: JWT{
				Secret: []byte(os.Getenv(prefix + "SERVER_JWT_SECRET")),
			},
		},
		InternalAPIConfig: api.ServiceConfig{
			InternalJWT: api.InternalJWT{
				Endpoint: os.Getenv(prefix + "SERVER_JWT_INTERNAL_ENDPOINT"),
				Login:    os.Getenv(prefix + "SERVER_JWT_INTERNAL_LOGIN"),
				Password: os.Getenv(prefix + "SERVER_JWT_INTERNAL_PASSWORD"),
			},
		},
		InternalAPIEndpoints: endpoints,
		//RabbitMQ: rabbitmq.Config{
		//	URI:                     os.Getenv(prefix + "RABBITMQ_AMQP_URI"),
		//	MaxReconnections:        rabbitmqMaxReconnections,
		//	ReconnectTimeoutSeconds: time.Duration(int64(rabbitmqReconnectTimeoutSeconds)) * time.Second,
		//},
		Database: Database{
			Host:     os.Getenv(prefix + "DB_HOST"),
			Port:     databasePort,
			User:     os.Getenv(prefix + "DB_USER"),
			Password: os.Getenv(prefix + "DB_PASSWORD"),
			Name:     os.Getenv(prefix + "DB_NAME"),
		},
		Redis: redis.Config{
			Addr:     os.Getenv(prefix + "REDIS_ADDR"),
			Password: os.Getenv(prefix + "REDIS_PASSWORD"),
			DB:       redisDB,
		},
	}

	return &config, nil
}