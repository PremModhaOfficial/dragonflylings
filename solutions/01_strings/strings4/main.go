package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func IncrementCounter(client *redis.Client, counterKey string) (int64, error) {
	ctx := context.Background()
	return client.Incr(ctx, counterKey).Result()
}

func GetCounter(client *redis.Client, counterKey string) (int64, error) {
	ctx := context.Background()
	return client.Get(ctx, counterKey).Int64()
}
