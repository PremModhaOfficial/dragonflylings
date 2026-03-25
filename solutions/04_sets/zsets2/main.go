package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func IsAllowed(client *redis.Client, rateLimitKey string, limit int, windowSeconds int64) (bool, error) {
	ctx := context.Background()
	now := time.Now().UnixNano()
	windowStart := now - windowSeconds*int64(time.Second)

	// Remove entries outside the sliding window
	if err := client.ZRemRangeByScore(ctx, rateLimitKey, "-inf",
		fmt.Sprintf("%d", windowStart)).Err(); err != nil {
		return false, err
	}

	// Count requests in the current window
	count, err := client.ZCount(ctx, rateLimitKey,
		fmt.Sprintf("%d", windowStart),
		fmt.Sprintf("%d", now),
	).Result()
	if err != nil {
		return false, err
	}

	if count >= int64(limit) {
		return false, nil
	}

	// Record this request with current timestamp as score
	return true, client.ZAdd(ctx, rateLimitKey, redis.Z{
		Score:  float64(now),
		Member: now,
	}).Err()
}
