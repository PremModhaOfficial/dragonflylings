package main

// CAPSTONE 2: Rate-Limited API Cache
//
// Time target: 45-60 minutes
// No hints available. You are the expert now.
//
// CHALLENGE: Build a production-grade API response cache that combines:
//   - Cache-aside pattern (read-through, write-on-miss)
//   - Sliding window rate limiter (sorted sets, per-user)
//   - Circuit breaker (protect against Dragonfly degradation)
//   - Pub/Sub cache invalidation (push invalidation across instances)
//
// ARCHITECTURE:
//
//   Request → [Rate Limiter] → [Circuit Breaker] → [Cache Lookup]
//                                                       ↓ miss
//                                                 [Fetch from origin]
//                                                       ↓
//                                                 [Store in cache]
//
//   Write/Update → Publish invalidation → Subscribers delete cache key
//
// BUGS TO FIX (there are 4):
//
// Bug 1 (isRateLimited): Always returns false — rate limiting never triggers.
//   Implement sliding window using sorted sets (same algorithm as prod3).
//   Rate limit key: "ratelimit:" + userID
//
// Bug 2 (canCall): Always returns true — circuit never opens.
//   Track failures: if failures >= threshold, open circuit for cooldown period.
//   Transition: closed→open on threshold, open→half-open after cooldown,
//   half-open→closed on success, half-open→open on failure.
//
// Bug 3 (Invalidate): Publishes to wrong channel format.
//   Channel should be: "cache:invalidate:" + channel (not just channel).
//   The Subscribe function listens on "cache:invalidate:" + channel.
//
// Bug 4 (Subscribe): Not implemented — just returns nil without subscribing.
//   Should subscribe to "cache:invalidate:" + channel,
//   and for each message, delete the cached key from Redis.
//
// Read the tests carefully. They verify each component individually and together.

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrRateLimited = errors.New("rate limit exceeded")
	ErrCircuitOpen = errors.New("circuit breaker is open: Redis calls are being short-circuited")
)

type cbState int

const (
	cbClosed   cbState = iota
	cbOpen
	cbHalfOpen
)

// RateLimitedAPICache combines cache-aside, rate limiting, circuit breaking,
// and pub/sub invalidation in one component.
type RateLimitedAPICache struct {
	client    *redis.Client
	cacheTTL  time.Duration
	rlLimit   int
	rlWindow  time.Duration
	mu        sync.Mutex
	failures  int
	threshold int
	openUntil time.Time
	state     cbState
	cooldown  time.Duration
}

// NewRateLimitedAPICache creates the cache with given configuration.
func NewRateLimitedAPICache(client *redis.Client, cacheTTL time.Duration, rlLimit int, rlWindow time.Duration, cbThreshold int, cbCooldown time.Duration) *RateLimitedAPICache {
	return &RateLimitedAPICache{
		client:    client,
		cacheTTL:  cacheTTL,
		rlLimit:   rlLimit,
		rlWindow:  rlWindow,
		threshold: cbThreshold,
		cooldown:  cbCooldown,
	}
}

// Fetch retrieves an API response from cache or origin, with rate limiting and circuit breaking.
// userID: used for per-user rate limiting
// cacheKey: the Redis key to cache the response under
// fetchFn: called on cache miss to get the fresh value from origin
func (c *RateLimitedAPICache) Fetch(ctx context.Context, userID, cacheKey string, fetchFn func() (string, error)) (string, error) {
	// Check rate limit
	limited, err := c.isRateLimited(ctx, userID)
	if err != nil {
		return "", err
	}
	if limited {
		return "", ErrRateLimited
	}

	// Check circuit breaker
	c.mu.Lock()
	if !c.canCall() {
		c.mu.Unlock()
		return "", ErrCircuitOpen
	}
	c.mu.Unlock()

	// Cache-aside: check cache
	val, err := c.client.Get(ctx, cacheKey).Result()
	if err == nil {
		c.mu.Lock()
		c.recordSuccess()
		c.mu.Unlock()
		return val, nil
	}
	if err != redis.Nil {
		c.mu.Lock()
		c.recordFailure()
		c.mu.Unlock()
		return "", err
	}

	// Cache miss: fetch from origin
	val, err = fetchFn()
	if err != nil {
		return "", err
	}

	// Store in cache
	setErr := c.client.Set(ctx, cacheKey, val, c.cacheTTL).Err()
	if setErr != nil {
		c.mu.Lock()
		c.recordFailure()
		c.mu.Unlock()
		return "", setErr
	}

	c.mu.Lock()
	c.recordSuccess()
	c.mu.Unlock()
	return val, nil
}

// isRateLimited checks if the user has exceeded their rate limit.
// BUG 1: always returns false — implement sliding window with sorted sets.
func (c *RateLimitedAPICache) isRateLimited(ctx context.Context, userID string) (bool, error) {
	// BUG: no rate limiting — always returns false
	return false, nil
}

// Invalidate publishes a cache invalidation message for the given key.
// BUG 3: publishes to wrong channel — should use "cache:invalidate:" + channel prefix.
func (c *RateLimitedAPICache) Invalidate(ctx context.Context, channel, key string) error {
	// BUG: publishes to bare channel name instead of "cache:invalidate:" + channel
	return c.client.Publish(ctx, channel, key).Err()
}

// Subscribe listens for cache invalidation messages and deletes affected keys.
// BUG 4: not implemented — just returns nil without subscribing.
// Should subscribe to "cache:invalidate:" + channel and delete keys on message.
func (c *RateLimitedAPICache) Subscribe(ctx context.Context, channel string) error {
	// BUG: not implemented
	return nil
}

// canCall returns true if the circuit breaker allows a call.
// BUG 2: always returns true — circuit never opens.
// Must be called with c.mu held.
func (c *RateLimitedAPICache) canCall() bool {
	// BUG: ignores circuit state
	return true
}

// recordSuccess resets failures and closes the circuit.
// Must be called with c.mu held.
func (c *RateLimitedAPICache) recordSuccess() {
	c.failures = 0
	c.state = cbClosed
}

// recordFailure increments failures and opens circuit if threshold reached.
// Must be called with c.mu held.
func (c *RateLimitedAPICache) recordFailure() {
	c.failures++
	if c.failures >= c.threshold {
		c.state = cbOpen
		c.openUntil = time.Now().Add(c.cooldown)
	}
}

// GetState returns the current circuit state string (for testing).
func (c *RateLimitedAPICache) GetState() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Check OPEN → HALF-OPEN transition
	if c.state == cbOpen && time.Now().After(c.openUntil) {
		c.state = cbHalfOpen
	}
	switch c.state {
	case cbClosed:
		return "closed"
	case cbOpen:
		return "open"
	case cbHalfOpen:
		return "half-open"
	}
	return "unknown"
}

// rateLimitKey returns the sorted set key for rate limiting.
func rateLimitKey(userID string) string {
	return "ratelimit:" + userID
}

// helper used by the sliding window implementation
func formatNano(t int64) string {
	return strconv.FormatInt(t, 10)
}
