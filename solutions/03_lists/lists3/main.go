package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func WaitForJob(client *redis.Client, queueKey string, timeout time.Duration) (string, error) {
	ctx := context.Background()
	result, err := client.BLPop(ctx, timeout, queueKey).Result()
	if err != nil {
		return "", err
	}
	// result[0] = key name, result[1] = value
	return result[1], nil
}

func EnqueueJob(client *redis.Client, queueKey, job string) error {
	ctx := context.Background()
	return client.RPush(ctx, queueKey, job).Err()
}
