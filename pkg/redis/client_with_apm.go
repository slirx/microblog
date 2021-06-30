package redis

import (
	"context"
	"fmt"
	"time"

	"go.elastic.co/apm"
)

// clientWithAPM adds APM to Client interface. It implements Decorator design pattern.
// Contexts for it's methods should contain APM transaction.
type clientWithAPM struct {
	Client Client
}

func (c clientWithAPM) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	tx := apm.TransactionFromContext(ctx)
	span := tx.StartSpan("redis.Set", "redis", nil)
	defer span.End()

	span.Action = "Set"
	span.Outcome = "success"
	span.Context.SetDatabase(apm.DatabaseSpanContext{
		Statement: fmt.Sprintf("key: %s; value: %v; expiration: %v", key, value, expiration),
	})

	err := c.Client.Set(ctx, key, value, expiration)
	if err != nil {
		span.Outcome = "error"
		return err
	}

	return nil
}

func (c clientWithAPM) Get(ctx context.Context, key string) (string, error) {
	tx := apm.TransactionFromContext(ctx)
	span := tx.StartSpan("redis.Get", "redis", nil)
	defer span.End()

	span.Action = "Get"
	span.Outcome = "success"
	span.Context.SetDatabase(apm.DatabaseSpanContext{
		Statement: fmt.Sprintf("key: %s", key),
	})

	response, err := c.Client.Get(ctx, key)
	if err != nil {
		span.Outcome = "error"
		return response, err
	}

	return response, nil
}

func (c clientWithAPM) SetIntSlice(ctx context.Context, key string, value []int, expiration time.Duration) error {
	tx := apm.TransactionFromContext(ctx)
	span := tx.StartSpan("redis.SetIntSlice", "redis", nil)
	defer span.End()

	span.Action = "SetIntSlice"
	span.Outcome = "success"
	span.Context.SetDatabase(apm.DatabaseSpanContext{
		Statement: fmt.Sprintf("key: %s; value: %v; expiration: %v", key, value, expiration),
	})

	err := c.Client.SetIntSlice(ctx, key, value, expiration)
	if err != nil {
		span.Outcome = "error"
		return err
	}

	return nil
}

func (c clientWithAPM) GetIntSlice(ctx context.Context, key string) ([]int, error) {
	tx := apm.TransactionFromContext(ctx)
	span := tx.StartSpan("redis.GetIntSlice", "redis", nil)
	defer span.End()

	span.Action = "GetIntSlice"
	span.Outcome = "success"
	span.Context.SetDatabase(apm.DatabaseSpanContext{
		Statement: fmt.Sprintf("key: %s", key),
	})

	response, err := c.Client.GetIntSlice(ctx, key)
	if err != nil {
		span.Outcome = "error"
		return response, err
	}

	return response, nil
}

func (c clientWithAPM) HIncrBy(ctx context.Context, key string, field string, incr int64) (int64, error) {
	tx := apm.TransactionFromContext(ctx)
	span := tx.StartSpan("redis.HIncrBy", "redis", nil)
	defer span.End()

	span.Action = "HIncrBy"
	span.Outcome = "success"
	span.Context.SetDatabase(apm.DatabaseSpanContext{
		Statement: fmt.Sprintf("key: %s; field: %s; incr: %d", key, field, incr),
	})

	response, err := c.Client.HIncrBy(ctx, key, field, incr)
	if err != nil {
		span.Outcome = "error"
		return response, err
	}

	return response, nil
}

func (c clientWithAPM) Close() error {
	return c.Client.Close()
}

func NewClientWithApm(c Client) Client {
	return clientWithAPM{
		Client: c,
	}
}
