package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"gitlab.com/slirx/newproj/internal/auth"
	"gitlab.com/slirx/newproj/pkg/logger"
	"gitlab.com/slirx/newproj/pkg/queue"
	"gitlab.com/slirx/newproj/pkg/queue/worker"
	"gitlab.com/slirx/newproj/pkg/rabbitmq"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/pq"
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
		log.Fatalln(err)
	}

	conf, err := NewConfig("AUTH_WORKER_UPDATE_USER_ID_")
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

	tracer := apm.DefaultTracer
	tracer.Service.Name = "auth-worker-update-user-id"

	w := worker.NewWorker(
		"auth/worker/update-user-id",
		rabbitmq.NewClient(conf.RabbitMQ, zapLogger, queue.JobAuthUpdateUserIDAuth),
		zapLogger,
		tracer,
		auth.NewUpdateUserIDHandler(zapLogger, auth.NewRepository(db)),
		queue.JobAuthUpdateUserIDAuth,
	)
	w.Run(ctx)
}
