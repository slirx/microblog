package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"

	"gitlab.com/slirx/newproj/internal/graphql/graph/generated"
	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/logger"
)

// NewHandler returns instance of GraphQL endpoint handler.
func NewHandler(l logger.Logger, rb api.ResponseBuilder, s *Service) http.Handler {
	return handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{Resolvers: &Resolver{
				Service:         s,
				Logger:          l,
				ResponseBuilder: rb,
			}},
		),
	)
}
