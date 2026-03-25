package main

// EXERCISE: prod5 - Hot Key Mitigation
//
// PREDICT: Before writing any code, answer in your head:
//   A "hot key" is a Redis key that receives a disproportionate number of
//   requests. For example, the "featured products" key is read 100,000 times/sec.
//   All those reads go to the same Dragonfly shard.
//   What's the bottleneck? Can throwing more Dragonfly nodes help?
//
// Hot keys are a shard-level bottleneck — more nodes don't help because
// all traffic hits the same shard. Solution: add a local in-process cache
// (sync.Map) with a short TTL. 99% of reads are served from memory.
// Only cache misses (or expired entries) hit Dragonfly.
//
// TODO: Add a local cache to HotKeyCache. When Get() is called:
//   1. Check local cache first — return if found and not expired
//   2. On local miss: fetch from Redis and store in local cache
//   3. Track how many Redis hits vs local hits for the test

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// localEntry caches a value with an expiry time.
type localEntry struct {
	value     string
	expiresAt time.Time
}

// HotKeyCache wraps a Redis client with a local in-process cache.
type HotKeyCache struct {
	redis    *redis.Client
	localTTL time.Duration
	// TODO: add localCache sync.Map for local caching
}

// NewHotKeyCache creates a HotKeyCache with the given local cache TTL.
func NewHotKeyCache(client *redis.Client, localTTL time.Duration) *HotKeyCache {
	return &HotKeyCache{
		redis:    client,
		localTTL: localTTL,
	}
}

// Get retrieves a value, checking the local cache before hitting Redis.
// BUG: no local cache — every call hits Redis.
// On a hot key with 100,000 req/sec, this saturates the Dragonfly shard.
func (c *HotKeyCache) Get(ctx context.Context, key string) (string, error) {
	// BUG: no local cache check — always hits Redis
	return c.redis.Get(ctx, key).Result()
}

// Set stores a value in Redis and invalidates the local cache entry.
func (c *HotKeyCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	// TODO: also invalidate local cache on Set
	return c.redis.Set(ctx, key, value, ttl).Err()
}

// RedisHitCount returns the number of times Redis was accessed.
// Used by tests to verify local cache is working.
// (In production you'd use metrics; here we use a simple counter.)
var redisHitCount int64
