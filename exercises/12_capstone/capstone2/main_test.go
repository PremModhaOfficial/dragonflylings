package main

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"dragonflylings/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestCache(t *testing.T) (*RateLimitedAPICache, context.Context) {
	t.Helper()
	client := testutil.NewTestClient(t)
	cache := NewRateLimitedAPICache(
		client,
		30*time.Second, // cacheTTL
		5,              // rate limit: 5 req
		time.Minute,    // per minute
		3,              // circuit: open after 3 failures
		200*time.Millisecond, // circuit: 200ms cooldown
	)
	return cache, context.Background()
}

// ── Cache-aside ────────────────────────────────────────────────────────────

func TestFetch_CacheHit(t *testing.T) {
	cache, ctx := newTestCache(t)
	cacheKey := testutil.UniqueKey("cap2-hit")
	userID := testutil.UniqueKey("user")

	var fetchCalls int64
	fetch := func() (string, error) {
		atomic.AddInt64(&fetchCalls, 1)
		return "origin-value", nil
	}

	// First call: cache miss → fetch from origin
	val, err := cache.Fetch(ctx, userID, cacheKey, fetch)
	require.NoError(t, err)
	assert.Equal(t, "origin-value", val)
	assert.Equal(t, int64(1), atomic.LoadInt64(&fetchCalls))

	// Second call: cache hit → no fetch
	val, err = cache.Fetch(ctx, userID, cacheKey, fetch)
	require.NoError(t, err)
	assert.Equal(t, "origin-value", val)
	assert.Equal(t, int64(1), atomic.LoadInt64(&fetchCalls), "fetch should not be called on cache hit")
}

// ── Rate limiting ──────────────────────────────────────────────────────────

func TestFetch_RateLimiting(t *testing.T) {
	cache, ctx := newTestCache(t)
	userID := testutil.UniqueKey("user-rl")

	// Make requests with different cache keys to avoid cache hits
	allowed := 0
	for i := 0; i < 10; i++ {
		cacheKey := testutil.UniqueKey("cap2-rl")
		_, err := cache.Fetch(ctx, userID, cacheKey, func() (string, error) {
			return "v", nil
		})
		if err == nil {
			allowed++
		} else if err != ErrRateLimited {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	assert.LessOrEqual(t, allowed, 5,
		"should allow at most 5 requests per minute; allowed %d — implement sliding window in isRateLimited()", allowed)
	assert.Greater(t, allowed, 0, "should allow some requests")
}

func TestFetch_RateLimitPerUser(t *testing.T) {
	cache, ctx := newTestCache(t)
	user1 := testutil.UniqueKey("user1")
	user2 := testutil.UniqueKey("user2")

	// Exhaust user1's rate limit
	for i := 0; i < 5; i++ {
		cache.Fetch(ctx, user1, testutil.UniqueKey("cap2"), func() (string, error) { return "v", nil }) //nolint:errcheck
	}

	// user1 should be rate limited
	_, err := cache.Fetch(ctx, user1, testutil.UniqueKey("cap2"), func() (string, error) { return "v", nil })
	assert.ErrorIs(t, err, ErrRateLimited, "user1 should be rate limited")

	// user2 should NOT be rate limited (separate bucket)
	_, err = cache.Fetch(ctx, user2, testutil.UniqueKey("cap2"), func() (string, error) { return "v", nil })
	assert.NoError(t, err, "user2 should not be affected by user1's rate limit")
}

// ── Circuit breaker ────────────────────────────────────────────────────────

func TestCircuitBreaker_OpensAfterFailures(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	cache := NewRateLimitedAPICache(client, time.Minute, 1000, time.Minute, 2, time.Second)

	// Force failures by pre-cancelling the context so Get always errors
	// (not redis.Nil), guaranteeing recordFailure() is triggered.
	for i := 0; i < 2; i++ {
		badCtx, cancel := context.WithCancel(ctx)
		cancel()
		key := testutil.UniqueKey("cap2-cb")
		cache.Fetch(badCtx, "user", key, func() (string, error) { return "", nil }) //nolint:errcheck
	}

	assert.Equal(t, "open", cache.GetState(),
		"circuit should be open after 2 failures — implement canCall() with state tracking")
}

func TestCircuitBreaker_OpenRejectsRequests(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	cache := NewRateLimitedAPICache(client, time.Minute, 1000, time.Minute, 1, 5*time.Second)

	// Force circuit open
	cache.mu.Lock()
	cache.recordFailure()
	cache.mu.Unlock()
	require.Equal(t, "open", cache.GetState())

	_, err := cache.Fetch(ctx, "user", testutil.UniqueKey("cap2"), func() (string, error) { return "v", nil })
	assert.ErrorIs(t, err, ErrCircuitOpen,
		"open circuit should reject requests immediately — fix canCall() to check state")
}

// ── Pub/Sub invalidation ───────────────────────────────────────────────────

func TestInvalidate_DeletesFromCache(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cache := NewRateLimitedAPICache(client, time.Minute, 1000, time.Minute, 10, time.Second)
	channel := testutil.UniqueKey("cap2-chan")
	cacheKey := testutil.UniqueKey("cap2-invalidate")

	// Pre-populate cache
	require.NoError(t, client.Set(ctx, cacheKey, "stale-value", time.Minute).Err())

	// Start subscriber in background
	subscribed := make(chan struct{})
	go func() {
		close(subscribed) // signal ready (simplification for test)
		cache.Subscribe(ctx, channel) //nolint:errcheck
	}()

	<-subscribed
	time.Sleep(50 * time.Millisecond) // let subscriber connect

	// Publish invalidation
	err := cache.Invalidate(ctx, channel, cacheKey)
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond) // let invalidation propagate

	// Key should be deleted
	exists, err := client.Exists(ctx, cacheKey).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists,
		"cache key should be deleted after invalidation message — "+
			"implement Subscribe() to listen and Delete(), fix Invalidate() channel prefix")
}

func TestInvalidate_ChannelFormat(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	cache := NewRateLimitedAPICache(client, time.Minute, 1000, time.Minute, 10, time.Second)
	channel := "my-channel"
	expectedChannel := "cache:invalidate:my-channel"

	// Subscribe to the expected channel
	sub := client.Subscribe(ctx, expectedChannel)
	defer sub.Close()

	// Give subscription time to register
	time.Sleep(20 * time.Millisecond)

	err := cache.Invalidate(ctx, channel, "some-key")
	require.NoError(t, err)

	// Should receive message on "cache:invalidate:my-channel", not "my-channel"
	msgCh := sub.Channel()
	select {
	case msg := <-msgCh:
		assert.Equal(t, expectedChannel, msg.Channel,
			"Invalidate should publish to 'cache:invalidate:'+channel, got %q", msg.Channel)
		assert.Equal(t, "some-key", msg.Payload)
	case <-time.After(500 * time.Millisecond):
		t.Fatal("no message received on 'cache:invalidate:my-channel' — fix Invalidate() to use correct channel prefix")
	}
}
