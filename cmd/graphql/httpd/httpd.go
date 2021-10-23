package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.elastic.co/apm"

	"gitlab.com/slirx/newproj/internal/graphql"
	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/http/apmmiddleware"
	"gitlab.com/slirx/newproj/pkg/http/jwtmiddleware"
	"gitlab.com/slirx/newproj/pkg/logger"
	"gitlab.com/slirx/newproj/pkg/tracer"
	"gitlab.com/slirx/newproj/pkg/utils"
)

func main() {
	// todo use signal.NotifyContext() instead. replace in all other executables
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()

	zapLogger, err := logger.NewZapLogger()
	if err != nil {
		log.Fatalf("can not initialize logger: %v", err)
	}

	conf, err := NewConfig("GRAPHQL_HTTPD_")
	if err != nil {
		zapLogger.Fatal(err)
	}

	t := tracer.NewAPMTracer()
	responseBuilder := api.NewResponseBuilder(t)

	service, err := graphql.NewService(conf.Endpoints, &conf.Server.ServiceConfig)
	if err != nil {
		zapLogger.Fatal(err)
	}

	handler := graphql.NewHandler(zapLogger, responseBuilder, service)
	apmTracer := apm.DefaultTracer

	recoveryFunc := utils.NewRecoveryFunc(zapLogger, responseBuilder)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   conf.Server.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Post(
		"/graphql",
		apmmiddleware.Wrap(
			jwtmiddleware.WrapInternal(
				handler.ServeHTTP,
				responseBuilder,
				zapLogger,
				conf.Server.JWT.InternalSecrets,
			),
			"/graphql",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)

	server := http.Server{
		Addr:    conf.Server.Addr,
		Handler: router,
	}

	go func() {
		_ = server.ListenAndServe()
	}()

	zapLogger.Debug("started")

	<-ctx.Done()
	_ = server.Shutdown(ctx)
}
