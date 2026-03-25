package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func SetPreferences(client *redis.Client, prefs map[string]string) error {
	ctx := context.Background()
	return client.MSet(ctx, prefs).Err()
}

func GetPreferences(client *redis.Client, keys []string) ([]interface{}, error) {
	ctx := context.Background()
	return client.MGet(ctx, keys...).Result()
}
