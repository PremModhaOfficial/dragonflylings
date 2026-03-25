package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllow_UnderLimit(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("prod3")

	// Send 5 requests with limit=10 — all should be allowed
	for i := 0; i < 5; i++ {
		allowed, err := Allow(ctx, client, key, 10, time.Minute)
		require.NoError(t, err)
		assert.True(t, allowed, "request %d should be allowed (under limit)", i+1)
	}
}

func TestAllow_AtLimit(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("prod3-limit")

	// Fill up to the limit
	for i := 0; i < 3; i++ {
		allowed, err := Allow(ctx, client, key, 3, time.Minute)
		require.NoError(t, err)
		require.True(t, allowed, "request %d should be allowed", i+1)
	}

	// Next request should be rejected
	allowed, err := Allow(ctx, client, key, 3, time.Minute)
	require.NoError(t, err)
	assert.False(t, allowed, "4th request should be rejected when limit=3")
}

func TestAllow_SlidingWindow_NoBoundaryBurst(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("prod3-sliding")

	// Use a short window for testing
	window := 500 * time.Millisecond
	limit := 3

	// Fill to limit
	for i := 0; i < limit; i++ {
		allowed, err := Allow(ctx, client, key, limit, window)
		require.NoError(t, err)
		require.True(t, allowed)
	}

	// Immediately rejected at limit
	allowed, err := Allow(ctx, client, key, limit, window)
	require.NoError(t, err)
	assert.False(t, allowed, "should be rejected at limit")

	// Wait for window to pass
	time.Sleep(window + 50*time.Millisecond)

	// Now should be allowed again — sliding window has cleared old entries
	allowed, err = Allow(ctx, client, key, limit, window)
	require.NoError(t, err)
	assert.True(t, allowed,
		"should be allowed after window expires — if this fails with a fixed-window "+
			"implementation, the bucket may have already reset; use ZCARD to verify sliding behavior")
}

func TestAllow_SlidingWindow_KeyStructure(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("prod3-structure")

	_, err := Allow(ctx, client, key, 10, time.Minute)
	require.NoError(t, err)

	// The key should be a sorted set (ZSET), not a string
	keyType, err := client.Type(ctx, key).Result()
	require.NoError(t, err)
	assert.Equal(t, "zset", keyType,
		"rate limiter key should be a sorted set (zset) for sliding window; got %q — "+
			"if it's 'string', you're using INCR (fixed window), not ZADD (sliding window)", keyType)
}

func TestRateLimitInfo(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("prod3-info")

	for i := 0; i < 3; i++ {
		_, err := Allow(ctx, client, key, 10, time.Minute)
		require.NoError(t, err)
	}

	count, err := RateLimitInfo(ctx, client, key, time.Minute)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count, fmt.Sprintf("should show 3 requests in window"))
}
