package main

// EXERCISE: hashes1 - HSET and HGET
//
// PREDICT: Before fixing anything, answer:
//   - What's the difference between SET "user:1:name" "alice" and HSET "user:1" "name" "alice"?
//   - How many keys does storing 5 user fields as Strings use? As a Hash?
//   - What does HGET return for a field that doesn't exist?
//
// The test models a user profile as a Redis Hash (one key, many fields).
// BUG: SetUserField uses client.Set (string key) instead of client.HSet (hash field).
//      GetUserField uses client.Get instead of client.HGet.
//
// TODO: Change Set→HSet and Get→HGet to use the Hash data structure.

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// SetUserField stores a single field of a user profile.
// BUG: Uses string Set — stores a new top-level key instead of a hash field.
func SetUserField(client *redis.Client, userKey, field, value string) error {
	ctx := context.Background()
	return client.HSet(ctx, userKey, field, value).Err() // TODO: use HSet
}

// GetUserField retrieves a single field from a user profile hash.
// BUG: Uses string Get — looks up wrong key.
func GetUserField(client *redis.Client, userKey, field string) (string, error) {
	ctx := context.Background()
	return client.HGet(ctx, userKey, field).Result() // TODO: use HGet
}

// DeleteUser removes the entire user hash.
func DeleteUser(client *redis.Client, userKey string) error {
	ctx := context.Background()
	return client.Del(ctx, userKey).Err()
}
