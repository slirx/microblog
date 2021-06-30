# MicroBlog (WIP)

this is my pet-project.

## TODO

- add license

- generate random hash for uploaded image (for example item id + unix time)

- health check / synthetic monitoring in kibana

- unit tests

- graphql

- golangci-lint

- profiling

- load tests

- implement rate limiter

## Technologies used

### vue.js

### yarn, webpack

Install webpack:
```bash
npm install --save-dev webpack
```

Build frontend (dev version):
```bash
yarn build-dev
```

### RabbitMQ

### Envoy

It's used as reverse proxy. It combines all API endpoints to one entry point, so frontend communicates only with that host/port.

### ELK (with APM)

## Design Patterns Used

- Decorator:

    pkg/redis/client_with_apm.go

- Repository:

  internal/post/repository.go


## Migrations

Install migrate command:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Create migration:

```bash
migrate -database "postgres://microblog:microblog@localhost:5432/microblog_user?sslmode=disable" -path . create -seq -ext sql add_name_column
```

Run migrations:

```bash
migrate -database "postgres://microblog:microblog@localhost:5432/microblog_auth?sslmode=disable" -path build/migrations/auth up
migrate -database "postgres://microblog:microblog@localhost:5432/microblog_post?sslmode=disable" -path build/migrations/post up
migrate -database "postgres://microblog:microblog@localhost:5432/microblog_registration?sslmode=disable" -path build/migrations/registration up
migrate -database "postgres://microblog:microblog@localhost:5432/microblog_user?sslmode=disable" -path build/migrations/user up
```

## GraphQL

For GraphQL interactions https://github.com/99designs/gqlgen repository is used.

Command for re-generating Go code from schema:

```bash
gqlgen generate
```

It has to be executed from internal/graphql directory.
