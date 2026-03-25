// Package testutil provides shared test helpers for dragonflylings exercises.
// All exercises connect to Dragonfly on localhost:6380.
package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

const addr = "localhost:6380"

// NewTestClient creates a go-redis client connected to the local Dragonfly instance.
// It verifies connectivity with PING and calls t.Fatal if the server is unreachable.
func NewTestClient(t *testing.T) *redis.Client {
	t.Helper()
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("Cannot connect to Dragonfly at %s: %v\nIs 'docker compose up -d' running?", addr, err)
	}
	return client
}

// UniqueKey returns a key prefixed with the given prefix and a nanosecond timestamp,
// ensuring uniqueness across parallel test runs.
func UniqueKey(prefix string) string {
	return fmt.Sprintf("%s:%d", prefix, time.Now().UnixNano())
}

// CleanupKeys deletes all keys matching the given pattern after the test completes.
// Pattern uses Redis glob syntax (e.g. "test:*").
func CleanupKeys(t *testing.T, client *redis.Client, pattern string) {
	t.Helper()
	t.Cleanup(func() {
		ctx := context.Background()
		var cursor uint64
		for {
			keys, next, err := client.Scan(ctx, cursor, pattern, 100).Result()
			if err != nil {
				return
			}
			if len(keys) > 0 {
				client.Del(ctx, keys...)
			}
			cursor = next
			if cursor == 0 {
				break
			}
		}
	})
}
