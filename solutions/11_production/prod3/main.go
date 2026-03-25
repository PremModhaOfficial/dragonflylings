package main

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Allow checks if a request should be allowed under the sliding window rate limit.
// Uses a sorted set where score = timestamp to implement a true rolling window.
func Allow(ctx context.Context, client *redis.Client, key string, limit int, window time.Duration) (bool, error) {
	now := time.Now()
	windowStart := now.Add(-window).UnixNano()
	member := strconv.FormatInt(now.UnixNano(), 10)

	pipe := client.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now.UnixNano()), Member: member})
	cardCmd := pipe.ZCard(ctx, key)
	pipe.Expire(ctx, key, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	return cardCmd.Val() <= int64(limit), nil
}

// RateLimitInfo returns the current request count in the window.
func RateLimitInfo(ctx context.Context, client *redis.Client, key string, window time.Duration) (count int64, err error) {
	now := time.Now()
	windowStart := now.Add(-window).UnixNano()
	count, err = client.ZCount(ctx, key,
		strconv.FormatInt(windowStart, 10), "+inf").Result()
	return
}
