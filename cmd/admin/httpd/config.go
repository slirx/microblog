package main

import (
	"os"
)

// Config represents combined configuration.
type Config struct {
	Server Server
}

// Server represents web server configuration.
type Server struct {
	// Addr represents address and port which server should listen to. It's specified in format host:port.
	Addr string
}

// NewConfig returns initialized instance of configuration. It reads configuration from environment variables.
// prefix represents prefix of environment variables' names.
func NewConfig(prefix string) (*Config, error) {
	config := Config{
		Server: Server{
			Addr: os.Getenv(prefix + "SERVER_ADDR"),
		},
	}

	return &config, nil
}
