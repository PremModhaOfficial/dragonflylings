package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func SetSession(client *redis.Client, sessionID, token string, ttl time.Duration) error {
	ctx := context.Background()
	return client.Set(ctx, "session:"+sessionID, token, ttl).Err()
}

func GetSession(client *redis.Client, sessionID string) (string, error) {
	ctx := context.Background()
	return client.Get(ctx, "session:"+sessionID).Result()
}

func GetTTL(client *redis.Client, sessionID string) (time.Duration, error) {
	ctx := context.Background()
	return client.TTL(ctx, "session:"+sessionID).Result()
}
