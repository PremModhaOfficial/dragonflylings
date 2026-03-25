package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func AddEvent(client *redis.Client, listKey, event string, maxItems int64) error {
	ctx := context.Background()
	if err := client.LPush(ctx, listKey, event).Err(); err != nil {
		return err
	}
	return client.LTrim(ctx, listKey, 0, maxItems-1).Err()
}

func GetRecentEvents(client *redis.Client, listKey string) ([]string, error) {
	ctx := context.Background()
	return client.LRange(ctx, listKey, 0, -1).Result()
}

func EventCount(client *redis.Client, listKey string) (int64, error) {
	ctx := context.Background()
	return client.LLen(ctx, listKey).Result()
}
