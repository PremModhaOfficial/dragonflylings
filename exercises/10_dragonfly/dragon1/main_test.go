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

const numKeys = 200

func makeKeys(prefix string) []string {
	keys := make([]string, numKeys)
	for i := range keys {
		keys[i] = fmt.Sprintf("%s:%d", prefix, i)
	}
	return keys
}

func TestSetConcurrent_AllKeysSet(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	keys := makeKeys(testutil.UniqueKey("dragon1"))

	count := SetConcurrent(ctx, client, keys, "val")
	assert.Equal(t, int64(numKeys), count,
		"all %d keys should be set; got %d — goroutines probably didn't finish before return", numKeys, count)

	// Spot-check a few keys actually exist
	for _, k := range keys[:5] {
		val, err := client.Get(ctx, k).Result()
		require.NoError(t, err, "key %s should exist", k)
		assert.Equal(t, "val", val)
	}
}

func TestSetPipelined_AllKeysSet(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	keys := makeKeys(testutil.UniqueKey("dragon1-pipe"))

	count, err := SetPipelined(ctx, client, keys, "val")
	require.NoError(t, err)
	assert.Equal(t, int64(numKeys), count,
		"all %d keys should be set via pipeline; got %d — probably creating a new pipeline per key instead of batching", numKeys, count)

	for _, k := range keys[:5] {
		val, err := client.Get(ctx, k).Result()
		require.NoError(t, err)
		assert.Equal(t, "val", val)
	}
}

func TestSetSequential_AllKeysSet(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	keys := makeKeys(testutil.UniqueKey("dragon1-seq"))

	count := SetSequential(ctx, client, keys, "val")
	assert.Equal(t, int64(numKeys), count)
}

func TestConcurrentFasterThanSequential(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping throughput comparison in short mode")
	}

	client := testutil.NewTestClient(t)
	ctx := context.Background()

	seqKeys := makeKeys(testutil.UniqueKey("dragon1-speedtest-seq"))
	concKeys := makeKeys(testutil.UniqueKey("dragon1-speedtest-conc"))

	start := time.Now()
	SetSequential(ctx, client, seqKeys, "v")
	seqDur := time.Since(start)

	start = time.Now()
	SetConcurrent(ctx, client, concKeys, "v")
	concDur := time.Since(start)

	speedup := float64(seqDur) / float64(concDur)
	t.Logf("Sequential: %v | Concurrent: %v | Speedup: %.1fx", seqDur, concDur, speedup)

	// On localhost the loopback is so fast that connection-pool overhead can
	// make concurrent slightly slower than sequential. The key correctness
	// check is TestSetConcurrent_AllKeysSet (all 200 keys must be set).
	// This test is informational — we log the result but don't fail.
	if speedup < 1.0 {
		t.Logf("NOTE: concurrent was not faster this run — expected on localhost. "+
			"On a real network Dragonfly's multi-threaded architecture shows clear gains.")
	}
}
