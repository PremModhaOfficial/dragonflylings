package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func SetUserField(client *redis.Client, userKey, field, value string) error {
	ctx := context.Background()
	return client.HSet(ctx, userKey, field, value).Err()
}

func GetUserField(client *redis.Client, userKey, field string) (string, error) {
	ctx := context.Background()
	return client.HGet(ctx, userKey, field).Result()
}

func DeleteUser(client *redis.Client, userKey string) error {
	ctx := context.Background()
	return client.Del(ctx, userKey).Err()
}
