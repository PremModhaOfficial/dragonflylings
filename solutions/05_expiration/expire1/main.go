package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// SetWithExpiry sets key=value and applies a TTL of duration d.
func SetWithExpiry(client *redis.Client, key, value string, d time.Duration) error {
	ctx := context.Background()
	if err := client.Set(ctx, key, value, 0).Err(); err != nil {
		return err
	}
	return client.Expire(ctx, key, d).Err()
}

// HasExpiry returns true if the key has an active expiration set.
func HasExpiry(client *redis.Client, key string) (bool, error) {
	ctx := context.Background()
	ttl, err := client.TTL(ctx, key).Result()
	if err != nil {
		return false, err
	}
	// Positive TTL = has expiry. -1 = persistent. -2 = doesn't exist.
	return ttl > 0, nil
}

// MakePersistent removes the expiration from key without deleting the key itself.
func MakePersistent(client *redis.Client, key string) error {
	ctx := context.Background()
	return client.Persist(ctx, key).Err()
}

func main() {}
