package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"dragonflylings/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache_Hit(t *testing.T) {
	client := testutil.NewTestClient(t)
	cache := NewCache(client, 30*time.Second)
	ctx := context.Background()
	key := testutil.UniqueKey("prod1-hit")

	// Pre-populate cache
	require.NoError(t, client.Set(ctx, key, "cached-value", 30*time.Second).Err())

	calls := 0
	val, err := cache.Get(ctx, key, func() (string, error) {
		calls++
		return "db-value", nil
	})
	require.NoError(t, err)
	assert.Equal(t, "cached-value", val)
	assert.Equal(t, 0, calls, "fetch should not be called on cache hit")
}

func TestCache_Miss(t *testing.T) {
	client := testutil.NewTestClient(t)
	cache := NewCache(client, 30*time.Second)
	ctx := context.Background()
	key := testutil.UniqueKey("prod1-miss")

	calls := 0
	val, err := cache.Get(ctx, key, func() (string, error) {
		calls++
		return "db-value", nil
	})
	require.NoError(t, err)
	assert.Equal(t, "db-value", val)
	assert.Equal(t, 1, calls, "fetch should be called once on cache miss")

	// Second call should hit cache
	val, err = cache.Get(ctx, key, func() (string, error) {
		calls++
		return "db-value-2", nil
	})
	require.NoError(t, err)
	assert.Equal(t, "db-value", val, "second call should return cached value")
	assert.Equal(t, 1, calls, "fetch should NOT be called again — value is now cached")
}

func TestCache_ThunderingHerd(t *testing.T) {
	client := testutil.NewTestClient(t)
	cache := NewCache(client, 30*time.Second)
	ctx := context.Background()
	key := testutil.UniqueKey("prod1-herd")

	var fetchCount int64
	var wg sync.WaitGroup
	const concurrency = 50

	// Simulate 50 goroutines all requesting the same missing key simultaneously
	start := make(chan struct{})
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start // wait for all goroutines to be ready
			_, err := cache.Get(ctx, key, func() (string, error) {
				atomic.AddInt64(&fetchCount, 1)
				time.Sleep(10 * time.Millisecond) // simulate DB latency
				return fmt.Sprintf("value-%d", time.Now().UnixNano()), nil
			})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	close(start) // release all goroutines at once
	wg.Wait()

	assert.Equal(t, int64(1), fetchCount,
		"singleflight should deduplicate: only 1 fetch() call expected for %d concurrent misses, got %d — "+
			"add a singleflight.Group to Cache and use it in Get()", concurrency, fetchCount)
}
