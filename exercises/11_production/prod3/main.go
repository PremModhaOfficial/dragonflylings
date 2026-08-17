package main

// EXERCISE: prod3 - Sliding Window Rate Limiter
//
// PREDICT: Before writing any code, answer in your head:
//   Fixed-window rate limiting: 10 requests per minute, window resets at :00.
//   User sends 10 requests at :59, then 10 more at :01. They sent 20 requests
//   in 2 seconds. Is this within the limit? (Yes, by fixed-window rules.)
//   Why is this a problem?
//
// The sliding window rate limiter uses a sorted set:
//   - Score = timestamp (nanoseconds)
//   - Member = unique ID per request (timestamp works for low traffic)
//   - Remove all members older than window start
//   - Count remaining members
//   - If count < limit: allow request (add new member)
//
// This prevents boundary bursts — the window always looks back exactly
// [window duration] from right now, not from a fixed clock boundary.
//
// TODO: Fix Allow() — it uses INCR (fixed window), not a sorted set.
// Replace with: ZREMRANGEBYSCORE (remove old) + ZADD + ZCARD + EXPIRE

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Allow checks if a request should be allowed under the rate limit.
// key: identifier for the rate limit bucket (e.g. "user:42" or "ip:1.2.3.4")
// limit: maximum requests per window
// window: rolling time window duration
// Returns true if the request is allowed, false if rate limited.
//
// BUG: uses INCR with a fixed window — allows boundary bursts.
// For example: 10 req/min limit, user sends 10 at :59 and 10 at :01.
// Fixed window allows both bursts (different minute buckets).
// Sliding window would correctly reject the second burst.
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

// RateLimitInfo returns the current request count and remaining capacity.
func RateLimitInfo(ctx context.Context, client *redis.Client, key string, window time.Duration) (count int64, err error) {
	now := time.Now()
	windowStart := now.Add(-window).UnixNano()

	count, err = client.ZCount(ctx, key,
		strconv.FormatInt(windowStart, 10), "+inf").Result()
	return
}
