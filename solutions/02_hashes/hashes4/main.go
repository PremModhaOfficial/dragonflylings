package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func StoreUser(client *redis.Client, userKey string, data map[string]string) error {
	ctx := context.Background()
	return client.HSet(ctx, userKey, data).Err()
}

func GetUser(client *redis.Client, userKey string) (map[string]string, error) {
	ctx := context.Background()
	return client.HGetAll(ctx, userKey).Result()
}

func DeleteUser(client *redis.Client, userKey string) error {
	ctx := context.Background()
	return client.Del(ctx, userKey).Err()
}
