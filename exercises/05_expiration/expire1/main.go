package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// EXERCISE: expire1 - The Janitor's Tools
//
// PREDICT: Before writing any code, answer in your head:
//   What does TTL return for a key with no expiry set? For a key that doesn't exist?
//   What is the difference between PERSIST and DEL?
//   If you call EXPIRE with duration 0, what happens?
//
// TODO: Fix the three bugs below. Each function has exactly one bug.

// SetWithExpiry sets key=value and applies a TTL of duration d.
func SetWithExpiry(client *redis.Client, key, value string, d time.Duration) error {
	ctx := context.Background()
	if err := client.Set(ctx, key, value, 0).Err(); err != nil {
		return err
	}
	// BUG: passing 0 as the duration -- should pass d
	return client.Expire(ctx, key, 0).Err()
}

// HasExpiry returns true if the key has an active expiration set.
func HasExpiry(client *redis.Client, key string) (bool, error) {
	ctx := context.Background()
	ttl, err := client.TTL(ctx, key).Result()
	if err != nil {
		return false, err
	}
	// BUG: TTL returns -1 for keys with NO expiry (persistent), not keys WITH expiry.
	// A positive TTL means the key has an active expiry.
	return ttl == -1, nil
}

// MakePersistent removes the expiration from key without deleting the key itself.
func MakePersistent(client *redis.Client, key string) error {
	ctx := context.Background()
	// BUG: Del deletes the key entirely. PERSIST removes expiry but keeps the value.
	return client.Del(ctx, key).Err()
}

func main() {}
