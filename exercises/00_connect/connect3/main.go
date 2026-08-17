package main

// EXERCISE: connect3 - Connection Pooling
//
// PREDICT: Before fixing anything, answer:
//   - Why would you need multiple connections to Dragonfly?
//   - What happens when 20 goroutines share 1 connection?
//   - What's the difference between PoolSize and MinIdleConns?
//
// The test sends 20 concurrent PINGs and checks pool stats.
// With PoolSize=1, goroutines must wait — they won't all run in parallel.
// With PoolSize=10 and MinIdleConns=5, connections are pre-warmed and ready.
//
// TODO: Fix the pool configuration so concurrent requests don't serialize.

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewPool creates a connection-pooled Redis client.
// BUG: PoolSize=1 means only one goroutine can use Dragonfly at a time.
// BUG: MinIdleConns=0 means no pre-warmed connections — cold start on each burst.
func NewPool(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         addr,
		PoolSize:     10, // TODO: increase to 10
		MinIdleConns: 5,  // TODO: set to 5
		PoolTimeout:  2 * time.Second,
	})
}

// Ping sends a PING and returns the latency.
func Ping(client *redis.Client) error {
	ctx := context.Background()
	return client.Ping(ctx).Err()
}

// PoolStats returns current pool statistics.
func PoolStats(client *redis.Client) *redis.PoolStats {
	return client.PoolStats()
}
