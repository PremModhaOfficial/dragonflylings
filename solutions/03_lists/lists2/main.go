package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func AddActivity(client *redis.Client, listKey, activity string) error {
	ctx := context.Background()
	return client.LPush(ctx, listKey, activity).Err()
}

func GetPage(client *redis.Client, listKey string, page, pageSize int64) ([]string, error) {
	ctx := context.Background()
	start := (page - 1) * pageSize
	stop := start + pageSize - 1
	return client.LRange(ctx, listKey, start, stop).Result()
}

func GetAll(client *redis.Client, listKey string) ([]string, error) {
	ctx := context.Background()
	return client.LRange(ctx, listKey, 0, -1).Result()
}
