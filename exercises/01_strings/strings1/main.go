package main

// EXERCISE: strings1 - SET and GET Basics
//
// PREDICT: Before fixing anything, answer:
//   - What does SET return? What does GET return for a missing key?
//   - Are Redis keys case-sensitive? Is "User:1" the same as "user:1"?
//   - What type does go-redis return from GET?
//
// The test stores a username and retrieves it.
// BUG: The Get call uses the wrong key — it looks up "user:name" but Set stored "user:1:name".
//
// TODO: Fix the key used in GetUsername so it matches the key used in SetUsername.

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// SetUsername stores a username in Dragonfly under a user-specific key.
func SetUsername(client *redis.Client, userID, username string) error {
	ctx := context.Background()
	return client.Set(ctx, "user:"+userID+":name", username, 0).Err()
}

// GetUsername retrieves a username from Dragonfly.
// BUG: The key is wrong — it doesn't match what SetUsername stored.
func GetUsername(client *redis.Client, userID string) (string, error) {
	ctx := context.Background()
	return client.Get(ctx, "user:name").Result() // TODO: fix the key
}
