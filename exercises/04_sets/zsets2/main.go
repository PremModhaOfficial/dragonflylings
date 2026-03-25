package main

// EXERCISE: zsets2 - Sliding Window Rate Limiter
//
// PREDICT: Before fixing anything, answer:
//   - How can a Sorted Set represent "requests in the last N seconds"?
//   - What score would you give each request? Why a timestamp?
//   - What does ZREMRANGEBYSCORE do, and why is it needed here?
//
// A sliding window rate limiter stores each request's timestamp as a ZSet score.
// To check the rate: count members with score in [now-window, now].
// To clean up: remove members with score < now-window BEFORE counting.
//
// BUG: IsAllowed counts members but never removes old entries.
//      Old entries keep accumulating and eventually block all new requests.
//
// TODO: Add ZREMRANGEBYSCORE before counting to remove entries outside the window.

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// IsAllowed returns true if the identifier has made fewer than limit
// requests in the last windowSeconds seconds.
// BUG: Uses ZCard (counts ALL members) instead of ZCount with window range,
//      AND never removes old entries. Old entries outside the window keep
//      accumulating and eventually block all new requests.
func IsAllowed(client *redis.Client, rateLimitKey string, limit int, windowSeconds int64) (bool, error) {
	ctx := context.Background()
	now := time.Now().UnixNano()

	// TODO: remove old entries before counting:
	// windowStart := now - windowSeconds*int64(time.Second)
	// client.ZRemRangeByScore(ctx, rateLimitKey, "-inf", fmt.Sprintf("%d", windowStart)).Err()

	// BUG: ZCard counts ALL members (including old ones outside the window).
	// Should use ZCount with [windowStart, now] range instead.
	count, err := client.ZCard(ctx, rateLimitKey).Result()
	if err != nil {
		return false, err
	}

	if count >= int64(limit) {
		return false, nil
	}

	// Record this request
	return true, client.ZAdd(ctx, rateLimitKey, redis.Z{
		Score:  float64(now),
		Member: now, // unique member per request
	}).Err()
}
