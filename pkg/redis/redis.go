package redis

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

var (
	ErrNoData = errors.New("no data")
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

type Client interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	SetIntSlice(ctx context.Context, key string, value []int, expiration time.Duration) error
	GetIntSlice(ctx context.Context, key string) ([]int, error)
	HIncrBy(ctx context.Context, key string, field string, incr int64) (int64, error)
	Close() error
}

type client struct {
	RedisClient *redis.Client
}

func (c client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := c.RedisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c client) Get(ctx context.Context, key string) (string, error) {
	result, err := c.RedisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrNoData
		}

		return "", errors.WithStack(err)
	}

	return result, nil
}

func (c client) SetIntSlice(ctx context.Context, key string, value []int, expiration time.Duration) error {
	b := make([]string, len(value))
	for i, v := range value {
		b[i] = strconv.Itoa(v)
	}

	err := c.RedisClient.Set(ctx, key, strings.Join(b, ","), expiration).Err()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c client) GetIntSlice(ctx context.Context, key string) ([]int, error) {
	result, err := c.RedisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNoData
		}

		return nil, errors.WithStack(err)
	}

	l := strings.Split(result, ",")
	response := make([]int, len(l))

	for i, s := range l {
		response[i], _ = strconv.Atoi(s)
	}

	return response, nil
}

func (c client) HIncrBy(ctx context.Context, key string, field string, incr int64) (int64, error) {
	v, err := c.RedisClient.HIncrBy(ctx, key, field, incr).Result()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return v, nil
}

func (c client) Close() error {
	err := c.RedisClient.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func New(ctx context.Context, conf Config) (Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return client{
		RedisClient: redisClient,
	}, nil
}
