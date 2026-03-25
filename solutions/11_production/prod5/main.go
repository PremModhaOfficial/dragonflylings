package main

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type localEntry struct {
	value     string
	expiresAt time.Time
}

// HotKeyCache wraps a Redis client with a local in-process cache to mitigate
// hot key saturation on a single Dragonfly shard.
type HotKeyCache struct {
	redis      *redis.Client
	localTTL   time.Duration
	localCache sync.Map
}

func NewHotKeyCache(client *redis.Client, localTTL time.Duration) *HotKeyCache {
	return &HotKeyCache{
		redis:    client,
		localTTL: localTTL,
	}
}

// Get checks the local cache first. On miss, fetches from Redis and caches locally.
func (c *HotKeyCache) Get(ctx context.Context, key string) (string, error) {
	// Check local cache
	if raw, ok := c.localCache.Load(key); ok {
		entry := raw.(localEntry)
		if time.Now().Before(entry.expiresAt) {
			return entry.value, nil
		}
		c.localCache.Delete(key) // expired
	}

	// Local miss — fetch from Redis
	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	// Store in local cache
	c.localCache.Store(key, localEntry{
		value:     val,
		expiresAt: time.Now().Add(c.localTTL),
	})
	return val, nil
}

// Set stores in Redis and invalidates the local cache entry.
func (c *HotKeyCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	c.localCache.Delete(key)
	return c.redis.Set(ctx, key, value, ttl).Err()
}

var redisHitCount int64
