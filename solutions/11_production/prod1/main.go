package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

// Cache implements a read-through cache with Redis backing and thundering herd protection.
type Cache struct {
	redis   *redis.Client
	ttl     time.Duration
	sfGroup singleflight.Group
}

// NewCache creates a Cache backed by the given Redis client with the given TTL.
func NewCache(client *redis.Client, ttl time.Duration) *Cache {
	return &Cache{
		redis: client,
		ttl:   ttl,
	}
}

// Get retrieves a value from cache. On cache miss, uses singleflight to ensure
// only ONE concurrent fetch() call per key, even under high concurrency.
func (c *Cache) Get(ctx context.Context, key string, fetch func() (string, error)) (string, error) {
	// Try cache first
	val, err := c.redis.Get(ctx, key).Result()
	if err == nil {
		return val, nil
	}
	if err != redis.Nil {
		return "", err
	}

	// Cache miss — use singleflight to deduplicate concurrent fetches
	result, err, _ := c.sfGroup.Do(key, func() (interface{}, error) {
		v, fetchErr := fetch()
		if fetchErr != nil {
			return "", fetchErr
		}
		// Populate cache (best-effort — ignore set errors)
		c.redis.Set(ctx, key, v, c.ttl)
		return v, nil
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}
