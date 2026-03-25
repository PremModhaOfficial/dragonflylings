package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func SetUserProfile(client *redis.Client, userKey string, fields map[string]interface{}) error {
	ctx := context.Background()
	return client.HSet(ctx, userKey, fields).Err()
}

func GetUserFields(client *redis.Client, userKey string, fields []string) ([]interface{}, error) {
	ctx := context.Background()
	return client.HMGet(ctx, userKey, fields...).Result()
}
