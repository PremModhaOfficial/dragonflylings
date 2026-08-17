package main

// EXERCISE: prod1 - Cache-Aside with Thundering Herd Protection
//
// PREDICT: Before writing any code, answer in your head:
//   Your cache key expires. At that exact moment, 500 goroutines all
//   try to read it. What happens? How many database calls result?
//   What if each DB call takes 100ms and the DB has a connection limit of 20?
//
// The thundering herd problem: when a popular cache key expires, every
// concurrent reader simultaneously misses the cache and hits the database.
// This can overwhelm the DB and cause cascading failures.
//
// Solution: singleflight — deduplicate concurrent fetches for the same key.
// Only ONE goroutine fetches from DB; all others wait and share the result.
//
// TODO: Add a singleflight.Group field to Cache and use it in Get() so that
// concurrent cache misses for the same key result in only ONE database call.

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

// Cache implements a read-through cache with Redis backing.
type Cache struct {
	redis *redis.Client
	ttl   time.Duration
	// TODO: add a singleflight.Group field here
	sfGroup singleflight.Group
}

// NewCache creates a Cache backed by the given Redis client with the given TTL.
func NewCache(client *redis.Client, ttl time.Duration) *Cache {
	return &Cache{
		redis: client,
		ttl:   ttl,
	}
}

// Get retrieves a value from cache. On cache miss, calls fetch() to load from
// the source of truth and populates the cache.
//
// BUG: No thundering herd protection. If 100 goroutines call Get() for the
// same missing key simultaneously, all 100 will call fetch() at the same time.
// This can overwhelm your database on cache misses at high traffic.
func (c *Cache) Get(ctx context.Context, key string, fetch func() (string, error)) (string, error) {
	// Try cache first
	val, err := c.redis.Get(ctx, key).Result()
	if err == nil {
		return val, nil
	}
	if err != redis.Nil {
		return "", err
	}

	result, err, _ := c.sfGroup.Do(key, func() (any, error) {
		v, fetchErr := fetch()
		if fetchErr != nil {
			return "", fetchErr
		}
		c.redis.Set(ctx, key, v, c.ttl)
		return v, nil
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}
