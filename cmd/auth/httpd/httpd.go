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

	"gitlab.com/slirx/newproj/internal/auth"
	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/http/apmmiddleware"
	"gitlab.com/slirx/newproj/pkg/logger"
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

	conf, err := NewConfig("AUTH_HTTPD_")
	if err != nil {
		zapLogger.Fatal(err)
	}

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

	service := auth.NewService(auth.NewRepository(db), conf.ServiceConfig)
	handler := auth.NewHandler(service, zapLogger, responseBuilder)

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
		"/auth/login",
		apmmiddleware.Wrap(
			handler.Login,
			"/auth/login",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)
	router.Post(
		"/internal/auth/login",
		apmmiddleware.Wrap(
			handler.InternalLogin,
			"/internal/auth/login",
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
