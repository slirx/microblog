package main

import (
	"os"
	"strings"
)

// Config represents combined configuration.
type Config struct {
	Server Server
}

type JWT struct {
	Secret          []byte
	InternalSecrets map[string][]byte // JWT secrets for microservices (service-to-service communication)
}

// Server represents web server configuration.
type Server struct {
	// Addr represents address and port which server should listen to. It's specified in format host:port.
	Addr string
	JWT  JWT
	// CORSAllowedOrigins is a list of origins a cross-domain request can be executed from.
	CORSAllowedOrigins []string
}

// NewConfig returns initialized instance of configuration. It reads configuration from environment variables.
// prefix represents prefix of environment variables' names.
func NewConfig(prefix string) (*Config, error) {
	internalSecrets := make(map[string][]byte)
	internalSecrets["user"] = []byte(os.Getenv(prefix + "SERVER_JWT_INTERNAL_USER_SECRET"))
	//internalSecrets["graphql"] = []byte(os.Getenv(prefix + "SERVER_JWT_INTERNAL_GRAPHQL_SECRET"))

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
		Server: Server{
			Addr: os.Getenv(prefix + "SERVER_ADDR"),
			JWT: JWT{
				Secret:          []byte(os.Getenv(prefix + "SERVER_JWT_SECRET")),
				InternalSecrets: internalSecrets,
			},
			CORSAllowedOrigins: corsAllowedOrigins,
		},
	}

	return &config, nil
}
