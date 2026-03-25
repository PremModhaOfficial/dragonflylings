package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Connect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6380",
	})
}

func Ping(client *redis.Client) (string, error) {
	ctx := context.Background()
	return client.Ping(ctx).Result()
}
