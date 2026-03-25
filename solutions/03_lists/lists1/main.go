package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func EnqueueTask(client *redis.Client, queueKey, task string) error {
	ctx := context.Background()
	return client.RPush(ctx, queueKey, task).Err()
}

func DequeueTask(client *redis.Client, queueKey string) (string, error) {
	ctx := context.Background()
	return client.LPop(ctx, queueKey).Result()
}

func QueueLength(client *redis.Client, queueKey string) (int64, error) {
	ctx := context.Background()
	return client.LLen(ctx, queueKey).Result()
}
