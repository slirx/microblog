package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/pq"

	"gitlab.com/slirx/newproj/internal/post"
	"gitlab.com/slirx/newproj/pkg/logger"
	"gitlab.com/slirx/newproj/pkg/queue"
	"gitlab.com/slirx/newproj/pkg/queue/worker"
	"gitlab.com/slirx/newproj/pkg/rabbitmq"
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

	conf, err := NewConfig("POST_WORKER_UNFOLLOW_")
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

	apmTrace := apm.DefaultTracer
	apmTrace.Service.Name = "post-worker-unfollow"

	w := worker.NewWorker(
		"post/worker/unfollow",
		rabbitmq.NewClient(conf.RabbitMQ, zapLogger, queue.JobPostUnfollow),
		zapLogger,
		apmTrace,
		post.NewUnfollowHandler(zapLogger, post.NewRepository(db)),
		queue.JobPostUnfollow,
	)
	w.Run(ctx)
}
