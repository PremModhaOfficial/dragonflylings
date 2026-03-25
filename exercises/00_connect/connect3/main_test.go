package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPoolConcurrent(t *testing.T) {
	client := NewPool("localhost:6380")
	require.NotNil(t, client)
	defer client.Close()

	// Send 20 concurrent PINGs
	const goroutines = 20
	errs := make(chan error, goroutines)
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs <- Ping(client)
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		require.NoError(t, err, "all concurrent pings should succeed")
	}

	stats := PoolStats(client)
	// With PoolSize=1, TotalConns will be 1 (bottleneck).
	// With PoolSize>=10, TotalConns should be higher to handle concurrency.
	assert.GreaterOrEqual(t, int(stats.TotalConns), 2,
		"pool should have created multiple connections for concurrent load")
}

func TestPoolHasIdleConns(t *testing.T) {
	client := NewPool("localhost:6380")
	require.NotNil(t, client)
	defer client.Close()

	// Warm the pool with one ping
	require.NoError(t, Ping(client))

	stats := PoolStats(client)
	assert.GreaterOrEqual(t, int(stats.IdleConns), 1,
		"pool should maintain idle connections (MinIdleConns > 0)")
}
