package main

// EXERCISE: hashes2 - HMGET Batch Field Operations
//
// PREDICT: Before fixing anything, answer:
//   - How many round trips does reading 5 hash fields one-by-one take?
//   - What does HMGET return for a field that doesn't exist in the hash?
//   - Can you read fields from a hash that doesn't exist? What do you get?
//
// The test reads multiple user profile fields in one call.
// BUG: GetUserFields loops over HGet calls (N round trips) instead of one HMGet.
//
// TODO: Replace the loop with a single client.HMGet(ctx, userKey, fields...).Result()

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// SetUserProfile stores all profile fields for a user at once.
func SetUserProfile(client *redis.Client, userKey string, fields map[string]interface{}) error {
	ctx := context.Background()
	return client.HSet(ctx, userKey, fields).Err()
}

// GetUserFields retrieves multiple fields from a user hash in one call.
// BUG: Loops with individual HGet calls (N round trips).
func GetUserFields(client *redis.Client, userKey string, fields []string) ([]interface{}, error) {
	ctx := context.Background()
	// TODO: replace this loop with client.HMGet(ctx, userKey, fields...).Result()

	return client.HMGet(ctx, userKey, fields...).Result()

	// results := make([]interface{}, len(fields))
	// for i, field := range fields {
	// 	val, err := client.HGet(ctx, userKey, field).Result()
	// 	if err != nil {
	// 		return nil, err // BUG: treats missing field as error; HMGet returns nil instead
	// 	}
	// 	results[i] = val
	// }
	// return results, nil
}
