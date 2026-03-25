package main

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

func (c *RateLimitedAPICache) Fetch(ctx context.Context, userID, cacheKey string, fetchFn func() (string, error)) (string, error) {
	// Rate limiting is infrastructure-level — use background context so a
	// cancelled caller context doesn't bypass rate limit tracking.
	limited, err := c.isRateLimited(context.Background(), userID)
	if err != nil {
		return "", err
	}
	if limited {
		return "", ErrRateLimited
	}

	c.mu.Lock()
	if !c.canCall() {
		c.mu.Unlock()
		return "", ErrCircuitOpen
	}
	c.mu.Unlock()

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

	val, err = fetchFn()
	if err != nil {
		return "", err
	}

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

// isRateLimited checks the sliding window rate limit for the user.
func (c *RateLimitedAPICache) isRateLimited(ctx context.Context, userID string) (bool, error) {
	now := time.Now()
	windowStart := now.Add(-c.rlWindow).UnixNano()
	member := strconv.FormatInt(now.UnixNano(), 10)
	key := rateLimitKey(userID)

	pipe := c.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", formatNano(windowStart))
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now.UnixNano()), Member: member})
	cardCmd := pipe.ZCard(ctx, key)
	pipe.Expire(ctx, key, c.rlWindow)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	return cardCmd.Val() > int64(c.rlLimit), nil
}

// Invalidate publishes a cache invalidation message.
func (c *RateLimitedAPICache) Invalidate(ctx context.Context, channel, key string) error {
	return c.client.Publish(ctx, "cache:invalidate:"+channel, key).Err()
}

// Subscribe listens for cache invalidation messages and deletes affected keys.
func (c *RateLimitedAPICache) Subscribe(ctx context.Context, channel string) error {
	sub := c.client.Subscribe(ctx, "cache:invalidate:"+channel)
	defer sub.Close()

	ch := sub.Channel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			// Delete the invalidated cache key
			c.client.Del(ctx, msg.Payload)
		}
	}
}

func (c *RateLimitedAPICache) canCall() bool {
	if c.state == cbOpen && time.Now().After(c.openUntil) {
		c.state = cbHalfOpen
	}
	return c.state == cbClosed || c.state == cbHalfOpen
}

func (c *RateLimitedAPICache) recordSuccess() {
	c.failures = 0
	c.state = cbClosed
}

func (c *RateLimitedAPICache) recordFailure() {
	c.failures++
	if c.failures >= c.threshold {
		c.state = cbOpen
		c.openUntil = time.Now().Add(c.cooldown)
	}
}

func (c *RateLimitedAPICache) GetState() string {
	c.mu.Lock()
	defer c.mu.Unlock()
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

func rateLimitKey(userID string) string {
	return "ratelimit:" + userID
}

func formatNano(t int64) string {
	return strconv.FormatInt(t, 10)
}
