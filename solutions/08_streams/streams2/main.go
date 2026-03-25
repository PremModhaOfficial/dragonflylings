package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func CountEvents(client *redis.Client, ctx context.Context, stream string) (int64, error) {
	return client.XLen(ctx, stream).Result()
}

func QueryRange(client *redis.Client, ctx context.Context, stream, start, stop string) ([]redis.XMessage, error) {
	return client.XRange(ctx, stream, start, stop).Result()
}

func main() {}
