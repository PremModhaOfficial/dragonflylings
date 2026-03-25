package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func SetUsername(client *redis.Client, userID, username string) error {
	ctx := context.Background()
	return client.Set(ctx, "user:"+userID+":name", username, 0).Err()
}

func GetUsername(client *redis.Client, userID string) (string, error) {
	ctx := context.Background()
	return client.Get(ctx, "user:"+userID+":name").Result()
}
