package main

import (
	"context"
	"testing"
	"time"

	"dragonflylings/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet_LocalCacheHit(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("prod5")

	require.NoError(t, client.Set(ctx, key, "hot-value", time.Minute).Err())
	cache := NewHotKeyCache(client, 5*time.Second)

	val1, err := cache.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, "hot-value", val1)

	// Delete key from Redis to prove second call uses local cache
	require.NoError(t, client.Del(ctx, key).Err())

	val2, err := cache.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, "hot-value", val2,
		"second Get should return value from local cache even though Redis key was deleted — "+
			"add a sync.Map local cache to HotKeyCache")
}

func TestGet_LocalCacheExpires(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("prod5-expire")

	require.NoError(t, client.Set(ctx, key, "v1", time.Minute).Err())
	cache := NewHotKeyCache(client, 100*time.Millisecond)

	val, err := cache.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, "v1", val)

	// Update value in Redis
	require.NoError(t, client.Set(ctx, key, "v2", time.Minute).Err())

	// Before TTL expires — local cache should serve old value
	val, err = cache.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, "v1", val, "local cache should serve v1 before TTL expires")

	// Wait for local cache TTL to expire
	time.Sleep(150 * time.Millisecond)

	// After TTL expires — should fetch v2 from Redis
	val, err = cache.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, "v2", val, "after local TTL expires, should fetch fresh value from Redis")
}

func TestSet_InvalidatesLocalCache(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("prod5-invalidate")

	require.NoError(t, client.Set(ctx, key, "original", time.Minute).Err())
	cache := NewHotKeyCache(client, 5*time.Second)

	_, err := cache.Get(ctx, key)
	require.NoError(t, err)

	// Update via cache.Set — should invalidate local cache
	err = cache.Set(ctx, key, "updated", time.Minute)
	require.NoError(t, err)

	// Next Get should return new value (not stale local cache)
	val, err := cache.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, "updated", val, "Set should invalidate local cache entry")
}

func TestGet_MissingKey(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	cache := NewHotKeyCache(client, 5*time.Second)
	key := testutil.UniqueKey("prod5-missing")

	// Key doesn't exist in Redis — should return redis.Nil error
	_, err := cache.Get(ctx, key)
	assert.Error(t, err, "missing key should return an error")
}
