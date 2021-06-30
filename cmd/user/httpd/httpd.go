package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/pq"

	"gitlab.com/slirx/newproj/internal/api/media"
	"gitlab.com/slirx/newproj/internal/user"
	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/http/apmmiddleware"
	"gitlab.com/slirx/newproj/pkg/http/jwtmiddleware"
	"gitlab.com/slirx/newproj/pkg/logger"
	"gitlab.com/slirx/newproj/pkg/queue/manager"
	"gitlab.com/slirx/newproj/pkg/tracer"
	"gitlab.com/slirx/newproj/pkg/utils"
)

func main() {
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

	conf, err := NewConfig("USER_HTTPD_")
	if err != nil {
		zapLogger.Fatal(err)
	}

	m := manager.NewManager(ctx, zapLogger, conf.RabbitMQ)
	defer func() {
		if err := m.Close(); err != nil {
			zapLogger.Error(err)
		}
	}()

	dbDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.Database.User, conf.Database.Password, conf.Database.Host, conf.Database.Port, conf.Database.Name,
	)
	db, err := apmsql.Open("postgres", dbDSN)
	if err != nil {
		zapLogger.Fatal(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			zapLogger.Error(err)
		}
	}()

	t := tracer.NewAPMTracer()
	responseBuilder := api.NewResponseBuilder(t)

	internalMediaAPI, err := media.NewAPI(conf.InternalAPIEndpoints, &conf.InternalAPIConfig)
	if err != nil {
		zapLogger.Fatal(err)
	}

	service := user.NewService(user.NewRepository(db), t, m, internalMediaAPI)
	handler := user.NewHandler(service, zapLogger, responseBuilder)

	apmTracer := apm.DefaultTracer

	recoveryFunc := utils.NewRecoveryFunc(zapLogger, responseBuilder)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://microblog.local:8001*"}, // todo move domain to config
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Patch(
		"/user",
		apmmiddleware.Wrap(
			jwtmiddleware.Wrap(
				handler.Update,
				responseBuilder,
				zapLogger,
				conf.Server.JWT.Secret,
			),
			"/user",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)
	router.Get(
		"/user/me",
		apmmiddleware.Wrap(
			jwtmiddleware.Wrap(
				handler.Me,
				responseBuilder,
				zapLogger,
				conf.Server.JWT.Secret,
			),
			"/user/me",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)
	router.Post(
		"/user/follow",
		apmmiddleware.Wrap(
			jwtmiddleware.Wrap(
				handler.Follow,
				responseBuilder,
				zapLogger,
				conf.Server.JWT.Secret,
			),
			"/user/follow",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)
	router.Post(
		"/user/unfollow",
		apmmiddleware.Wrap(
			jwtmiddleware.Wrap(
				handler.Unfollow,
				responseBuilder,
				zapLogger,
				conf.Server.JWT.Secret,
			),
			"/user/unfollow",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)
	router.Get(
		"/user/{login}",
		apmmiddleware.Wrap(
			jwtmiddleware.Wrap(
				handler.Get,
				responseBuilder,
				zapLogger,
				conf.Server.JWT.Secret,
			),
			"/user/{login}",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)
	router.Get(
		"/internal/user/{login}",
		apmmiddleware.Wrap(
			jwtmiddleware.WrapInternal(
				handler.Get, // todo use InternalGet?
				responseBuilder,
				zapLogger,
				conf.Server.JWT.InternalSecrets,
			),
			"/internal/user/{login}",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)
	router.Get(
		"/internal/user/{uid}/followers",
		apmmiddleware.Wrap(
			jwtmiddleware.WrapInternal(
				handler.InternalFollowers,
				responseBuilder,
				zapLogger,
				conf.Server.JWT.InternalSecrets,
			),
			"/internal/user/{uid}/followers",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)
	router.Get(
		"/user/{login}/followers",
		apmmiddleware.Wrap(
			jwtmiddleware.Wrap(
				handler.Followers,
				responseBuilder,
				zapLogger,
				conf.Server.JWT.Secret,
			),
			"/user/{login}/followers",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)
	router.Get(
		"/user/{login}/following",
		apmmiddleware.Wrap(
			jwtmiddleware.Wrap(
				handler.Following,
				responseBuilder,
				zapLogger,
				conf.Server.JWT.Secret,
			),
			"/user/{login}/following",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)
	router.Get(
		"/user/{login}/following",
		apmmiddleware.Wrap(
			jwtmiddleware.Wrap(
				handler.Following,
				responseBuilder,
				zapLogger,
				conf.Server.JWT.Secret,
			),
			"/user/{login}/following",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)
	router.Get(
		"/internal/user",
		apmmiddleware.Wrap(
			jwtmiddleware.WrapInternal(
				handler.InternalUsers,
				responseBuilder,
				zapLogger,
				conf.Server.JWT.InternalSecrets,
			),
			"/internal/user",
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
