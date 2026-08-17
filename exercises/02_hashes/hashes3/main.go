package main

// EXERCISE: hashes3 - HINCRBY for Hash Counters
//
// PREDICT: Before fixing anything, answer:
//   - What's the difference between INCR (on a string key) and HINCRBY (on a hash field)?
//   - Can you have a counter field inside a hash alongside other string fields?
//   - What does HINCRBY return?
//
// The test tracks page views per page inside a single "analytics" hash.
// BUG: IncrPageView uses client.Incr (increments a top-level string key)
//      instead of client.HIncrBy (increments a field inside a hash).
//
// TODO: Change Incr → HIncrBy to store the counter as a hash field.

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// IncrPageView increments the view count for a specific page.
// BUG: Uses Incr on a string key — creates separate top-level keys instead of hash fields.
func IncrPageView(client *redis.Client, analyticsKey, page string) (int64, error) {
	ctx := context.Background()
	return client.HIncrBy(ctx, analyticsKey, page, 1).Result() // TODO: use HIncrBy
}

// GetPageViews retrieves the view count for a specific page.
// BUG: Uses Get instead of HGet.
func GetPageViews(client *redis.Client, analyticsKey, page string) (int64, error) {
	ctx := context.Background()
	return client.HGet(ctx, analyticsKey, page).Int64() // TODO: use HGet(...).Int64()
}

// GetAllPageViews retrieves all page view counts as a map.
func GetAllPageViews(client *redis.Client, analyticsKey string) (map[string]string, error) {
	ctx := context.Background()
	return client.HGetAll(ctx, analyticsKey).Result()
}
