// main package represents executable for sending emails. it's rabbitmq worker which reads one message at a time and
// sends email specified in message.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"go.elastic.co/apm"

	"gitlab.com/slirx/newproj/pkg/queue"
	"gitlab.com/slirx/newproj/pkg/rabbitmq"

	"gitlab.com/slirx/newproj/internal/email"
	"gitlab.com/slirx/newproj/pkg/logger"
	"gitlab.com/slirx/newproj/pkg/queue/worker"
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

	conf, err := NewConfig("EMAIL_SENDER_")
	if err != nil {
		zapLogger.Fatal(err)
	}

	tracer := apm.DefaultTracer
	tracer.Service.Name = "email-worker"

	w := worker.NewWorker(
		"email/send",
		rabbitmq.NewClient(conf.RabbitMQ, zapLogger, queue.JobEmailSend),
		zapLogger,
		tracer,
		email.NewHandler(zapLogger),
		queue.JobEmailSend,
	)
	w.Run(ctx)
}
