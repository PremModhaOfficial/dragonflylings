package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func AcquireLock(client *redis.Client, lockKey, ownerID string, ttl time.Duration) (bool, error) {
	ctx := context.Background()
	return client.SetNX(ctx, lockKey, ownerID, ttl).Result()
}

func ReleaseLock(client *redis.Client, lockKey string) error {
	ctx := context.Background()
	return client.Del(ctx, lockKey).Err()
}
