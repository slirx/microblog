package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gitlab.com/slirx/newproj/internal/auth"
	"gitlab.com/slirx/newproj/pkg/rabbitmq"
)

// Config represents combined configuration.
type Config struct {
	Server        Server
	RabbitMQ      rabbitmq.Config
	Database      Database
	ServiceConfig auth.Config
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// Server represents web server configuration.
type Server struct {
	// Addr represents address and port which server should listen to. It's specified in format host:port.
	Addr string
	// AdminAddr represents address and port which server should listen to. It's specified in format host:port.
	// It's used for admin panel.
	AdminAddr string
	// CORSAllowedOrigins is a list of origins a cross-domain request can be executed from.
	CORSAllowedOrigins []string
	// AdminCORSAllowedOrigins is a list of origins a cross-domain request can be executed from (for admin panel).
	AdminCORSAllowedOrigins []string
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

	internalSecrets := make(map[string]string)
	internalSecrets["post"] = os.Getenv(prefix + "SERVICE_INTERNAL_POST_SECRET")
	internalSecrets["user"] = os.Getenv(prefix + "SERVICE_INTERNAL_USER_SECRET")
	internalSecrets["graphql"] = os.Getenv(prefix + "SERVICE_INTERNAL_GRAPHQL_SECRET")

	corsAllowedOrigins := make([]string, 0)
	tmpAllowedOrigins := strings.Split(os.Getenv(prefix+"SERVER_CORS_ALLOWED_ORIGINS"), ",")
	for _, origin := range tmpAllowedOrigins {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}

		corsAllowedOrigins = append(corsAllowedOrigins, origin)
	}

	adminCORSAllowedOrigins := make([]string, 0)
	tmpAllowedOrigins = strings.Split(os.Getenv(prefix+"SERVER_ADMIN_CORS_ALLOWED_ORIGINS"), ",")
	for _, origin := range tmpAllowedOrigins {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}

		adminCORSAllowedOrigins = append(adminCORSAllowedOrigins, origin)
	}

	config := Config{
		Server: Server{
			Addr:                    os.Getenv(prefix + "SERVER_ADDR"),
			AdminAddr:               os.Getenv(prefix + "SERVER_ADMIN_ADDR"),
			CORSAllowedOrigins:      corsAllowedOrigins,
			AdminCORSAllowedOrigins: adminCORSAllowedOrigins,
		},
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
		ServiceConfig: auth.Config{
			Secret:          os.Getenv(prefix + "SERVICE_SECRET"),
			InternalSecrets: internalSecrets,
		},
	}

	return &config, nil
}
