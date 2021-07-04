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

	"gitlab.com/slirx/newproj/internal/registration"
	"gitlab.com/slirx/newproj/pkg/api"
	"gitlab.com/slirx/newproj/pkg/http/apmmiddleware"
	"gitlab.com/slirx/newproj/pkg/logger"
	"gitlab.com/slirx/newproj/pkg/queue/manager"
	"gitlab.com/slirx/newproj/pkg/template"
	"gitlab.com/slirx/newproj/pkg/tracer"
	"gitlab.com/slirx/newproj/pkg/utils"
)

// todo configure distributed tracing: https://www.elastic.co/guide/en/apm/get-started/current/distributed-tracing.html

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

	conf, err := NewConfig("REGISTRATION_HTTPD_")
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

	tg := template.New()

	service := registration.NewService(t, registration.NewRepository(db), m, tg)
	handler := registration.NewHandler(service, zapLogger, responseBuilder)

	apmTracer := apm.DefaultTracer

	recoveryFunc := utils.NewRecoveryFunc(zapLogger, responseBuilder)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   conf.Server.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Post(
		"/registration",
		apmmiddleware.Wrap(
			handler.Register,
			"/registration",
			apmmiddleware.WithTracer(apmTracer),
			apmmiddleware.WithRecovery(recoveryFunc),
		),
	)
	router.Post("/registration/confirm",
		apmmiddleware.Wrap(
			handler.Confirm,
			"/registration/confirm",
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
