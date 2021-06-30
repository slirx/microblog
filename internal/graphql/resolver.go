package graphql

import (
	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/logger"
)

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	Service         *Service
	Logger          logger.Logger
	ResponseBuilder api.ResponseBuilder
}
