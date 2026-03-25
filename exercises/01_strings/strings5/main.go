package main

// EXERCISE: strings5 - MGET/MSET Batch Operations
//
// PREDICT: Before fixing anything, answer:
//   - How many round trips does setting 5 keys individually take?
//   - How many round trips does MSET take for the same 5 keys?
//   - What does MGET return for a key that doesn't exist — nil or error?
//
// The test stores and retrieves multiple user preferences at once.
// BUG: GetPreferences returns an error for missing keys instead of nil.
//      MGET returns nil for missing keys (not an error) — that's the point.
//
// TODO: Replace the loop in GetPreferences with a single MGet call,
//       and replace the loop in SetPreferences with a single MSet call.

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// SetPreferences stores multiple key-value preferences.
// BUG: Uses N individual Set calls (N round trips) instead of one MSet.
func SetPreferences(client *redis.Client, prefs map[string]string) error {
	ctx := context.Background()
	// TODO: replace this loop with client.MSet(ctx, prefs).Err()
	for k, v := range prefs {
		if err := client.Set(ctx, k, v, 0).Err(); err != nil {
			return err
		}
	}
	return nil
}

// GetPreferences retrieves values for the given keys in order.
// Missing keys should return nil in the result slice — NOT an error.
// BUG: This returns an error when a key is missing. MGET never errors on missing keys.
func GetPreferences(client *redis.Client, keys []string) ([]interface{}, error) {
	ctx := context.Background()
	// TODO: replace this with client.MGet(ctx, keys...).Result()
	results := make([]interface{}, len(keys))
	for i, k := range keys {
		val, err := client.Get(ctx, k).Result()
		if err != nil {
			// BUG: treating redis.Nil as an error — should store nil in results instead
			return nil, err
		}
		results[i] = val
	}
	return results, nil
}
