package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func IncrPageView(client *redis.Client, analyticsKey, page string) (int64, error) {
	ctx := context.Background()
	return client.HIncrBy(ctx, analyticsKey, page, 1).Result()
}

func GetPageViews(client *redis.Client, analyticsKey, page string) (int64, error) {
	ctx := context.Background()
	return client.HGet(ctx, analyticsKey, page).Int64()
}

func GetAllPageViews(client *redis.Client, analyticsKey string) (map[string]string, error) {
	ctx := context.Background()
	return client.HGetAll(ctx, analyticsKey).Result()
}
