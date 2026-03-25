package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateLimiterAllows(t *testing.T) {
	client := testutil.NewTestClient(t)
	key := fmt.Sprintf("test:zsets2:rate:%d", time.Now().UnixNano())

	// 3 requests within limit of 5 should all be allowed
	for i := 0; i < 3; i++ {
		allowed, err := IsAllowed(client, key, 5, 60)
		require.NoError(t, err)
		assert.True(t, allowed, "request %d should be allowed (under limit)", i+1)
	}
}

func TestRateLimiterBlocks(t *testing.T) {
	client := testutil.NewTestClient(t)
	key := fmt.Sprintf("test:zsets2:block:%d", time.Now().UnixNano())

	// Use up all 3 slots
	for i := 0; i < 3; i++ {
		_, _ = IsAllowed(client, key, 3, 60)
	}

	// 4th request should be blocked
	allowed, err := IsAllowed(client, key, 3, 60)
	require.NoError(t, err)
	assert.False(t, allowed, "4th request should be blocked (over limit)")
}

func TestRateLimiterCleansOldEntries(t *testing.T) {
	client := testutil.NewTestClient(t)
	key := fmt.Sprintf("test:zsets2:cleanup:%d", time.Now().UnixNano())
	ctx := context.Background()

	// Manually add 3 entries from 10 seconds ago (outside the 2s window)
	oldTime := float64(time.Now().Add(-10 * time.Second).UnixNano())
	for i := 0; i < 3; i++ {
		client.ZAdd(ctx, key, redis.Z{
			Score:  oldTime + float64(i),
			Member: fmt.Sprintf("old-%d", i),
		})
	}

	// With limit=3 and window=2s:
	// Broken code (no cleanup): counts 3 old entries → blocked
	// Fixed code (with ZREMRANGEBYSCORE): old entries removed → allowed
	allowed, err := IsAllowed(client, key, 3, 2)
	require.NoError(t, err)
	assert.True(t, allowed,
		"old entries (outside window) should be cleaned up; without ZREMRANGEBYSCORE they block new requests")
}
