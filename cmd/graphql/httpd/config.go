package main

import (
	"os"
	"strings"

	"gitlab.com/slirx/newproj/internal/api"
	"gitlab.com/slirx/newproj/pkg/rabbitmq"
)

// Config represents combined configuration.
type Config struct {
	EnvoyURL  string
	Server    Server
	Endpoints map[string]string
	RabbitMQ  rabbitmq.Config
}

type JWT struct {
	Secret []byte
}

// Server represents web server configuration.
type Server struct {
	// Addr represents address and port which server should listen to. It's specified in format host:port.
	Addr          string
	JWT           JWT
	ServiceConfig api.ServiceConfig
	// CORSAllowedOrigins is a list of origins a cross-domain request can be executed from.
	CORSAllowedOrigins []string
}

// NewConfig returns initialized instance of configuration. It reads configuration from environment variables.
// prefix represents prefix of environment variables' names.
func NewConfig(prefix string) (*Config, error) {
	//var rabbitmqMaxReconnections int
	//if rabbitmqMaxReconnections, err = strconv.Atoi(os.Getenv(prefix + "RABBITMQ_MAX_RECONNECTIONS")); err != nil {
	//	return nil, errors.WithStack(fmt.Errorf("invalid %s value", prefix+"RABBITMQ_MAX_RECONNECTIONS"))
	//}

	//var rabbitmqReconnectTimeoutSeconds int

	//rabbitmqReconnectTimeoutSeconds, err = strconv.Atoi(os.Getenv(prefix + "RABBITMQ_RECONNECT_TIMEOUT_SECONDS"))
	//if err != nil {
	//	return nil, errors.WithStack(fmt.Errorf("invalid %s value", prefix+"RABBITMQ_RECONNECT_TIMEOUT_SECONDS"))
	//}

	endpoints := make(map[string]string)
	endpoints["auth"] = os.Getenv(prefix + "ENDPOINT_AUTH")
	endpoints["user"] = os.Getenv(prefix + "ENDPOINT_USER")

	corsAllowedOrigins := make([]string, 0)
	tmpAllowedOrigins := strings.Split(os.Getenv(prefix+"SERVER_CORS_ALLOWED_ORIGINS"), ",")
	for _, origin := range tmpAllowedOrigins {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}

		corsAllowedOrigins = append(corsAllowedOrigins, origin)
	}

	config := Config{
		EnvoyURL: os.Getenv(prefix + "ENVOY_URL"),
		Server: Server{
			Addr: os.Getenv(prefix + "SERVER_ADDR"),
			JWT: JWT{
				Secret: []byte(os.Getenv(prefix + "SERVER_JWT_SECRET")),
			},
			ServiceConfig: api.ServiceConfig{
				InternalJWT: api.InternalJWT{
					Endpoint: os.Getenv(prefix + "SERVER_JWT_INTERNAL_ENDPOINT"),
					Login:    os.Getenv(prefix + "SERVER_JWT_INTERNAL_LOGIN"),
					Password: os.Getenv(prefix + "SERVER_JWT_INTERNAL_PASSWORD"),
				},
			},
			CORSAllowedOrigins: corsAllowedOrigins,
		},
		Endpoints: endpoints,
		//RabbitMQ: rabbitmq.Config{
		//	URI:                     os.Getenv(prefix + "RABBITMQ_AMQP_URI"),
		//	MaxReconnections:        rabbitmqMaxReconnections,
		//	ReconnectTimeoutSeconds: time.Duration(int64(rabbitmqReconnectTimeoutSeconds)) * time.Second,
		//},
	}

	return &config, nil
}
