package main

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/pq"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/slirx/newproj/pkg/logger"
)

func main() {
	// todo add 10m new posts

	//pass, _ := bcrypt.GenerateFromPassword([]byte("user-password"), bcrypt.DefaultCost)
	//fmt.Println(string(pass))
	//err := createNewPosts()
	//if err != nil {
	//	fmt.Println(err)
	//}

	send("some msg")

	fmt.Println("done")
}

func createNewUsers() error {
	zapLogger, err := logger.NewZapLogger()
	if err != nil {
		log.Fatalf("can not initialize logger: %v", err)
	}

	dbDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		"microblog", "microblog", "localhost", 5432, "microblog_user",
	)
	dbUser, err := apmsql.Open("postgres", dbDSN)
	if err != nil {
		zapLogger.Fatal(err)
	}

	defer func() {
		if err := dbUser.Close(); err != nil {
			zapLogger.Error(err)
		}
	}()

	dbDSN = fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		"microblog", "microblog", "localhost", 5432, "microblog_auth",
	)
	dbAuth, err := apmsql.Open("postgres", dbDSN)
	if err != nil {
		zapLogger.Fatal(err)
	}

	defer func() {
		if err := dbAuth.Close(); err != nil {
			zapLogger.Error(err)
		}
	}()

	pass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	authSTMT, err := dbAuth.Prepare(`INSERT INTO "auth" (user_id, email, login, password) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return errors.WithStack(err)
	}

	var id int

	for i := 2000; i < 10000; i++ {
		err = dbUser.QueryRow(
			`INSERT INTO "user" (email, login, bio, name) VALUES ($1, $2, $3, $4) RETURNING id`,
			fmt.Sprintf("test-usr%d@test.com", i),
			fmt.Sprintf("test-usr%d", i),
			fmt.Sprintf("bio for user %d..", i),
			fmt.Sprintf("Name%d", i),
		).Scan(&id)
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = authSTMT.Exec(
			id,
			fmt.Sprintf("test-usr%d@test.com", i),
			fmt.Sprintf("test-usr%d", i),
			string(pass),
		)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func createNewPosts() error {
	zapLogger, err := logger.NewZapLogger()
	if err != nil {
		log.Fatalf("can not initialize logger: %v", err)
	}

	dbDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		"microblog", "microblog", "localhost", 5432, "microblog_post",
	)
	dbPost, err := apmsql.Open("postgres", dbDSN)
	if err != nil {
		zapLogger.Fatal(err)
	}

	defer func() {
		if err := dbPost.Close(); err != nil {
			zapLogger.Error(err)
		}
	}()

	var id int

	for i := 10; i < 33; i++ {
		err = dbPost.QueryRow(
			`INSERT INTO "post" (user_id, text) VALUES (2, $1) RETURNING id`,
			fmt.Sprintf("test message #%d", i),
		).Scan(&id)
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = dbPost.Exec(
			`INSERT INTO "feed" (user_id, post_id) VALUES (2, $1)`,
			id,
		)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

type Notifier interface {
	Send(msg string)
}

type Email interface {
	Send(msg string)
}

type SMS interface {
	Send(msg string)
}

type sms struct {
	Notifier Notifier
}

func (s sms) Send(msg string) {
	fmt.Println("send sms")
	s.Notifier.Send(msg)
}

type email struct {
	//Notifier Notifier
}

func (s email) Send(msg string) {
	fmt.Println("send email")
	//s.Notifier.Send(msg)
}

func send(msg string) {
	e := email{}
	s := sms{
		Notifier: e,
	}
	s.Send(msg)
}
