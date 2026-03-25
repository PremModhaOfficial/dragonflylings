package main

import (
	"context"
	"fmt"
	"testing"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMemoryLimit(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	// Read current value to verify GetMemoryLimit works
	limit, err := GetMemoryLimit(client)
	require.NoError(t, err, "GetMemoryLimit should not error (bug uses wrong config key 'max-memory')")
	assert.Greater(t, limit, int64(0), "maxmemory should be a positive value in this setup")
}

func TestSetAndGetMemoryLimit(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// Save original to restore after test
	original, err := GetMemoryLimit(client)
	require.NoError(t, err)
	t.Cleanup(func() {
		// Restore original -- never set to 0 (causes Dragonfly OOM with rss_oom_deny_ratio)
		client.ConfigSet(context.Background(), "maxmemory", fmt.Sprintf("%d", original))
	})

	// Set a large limit that's different from original
	testLimit := original / 2
	if testLimit < 512*1024*1024 {
		testLimit = 2 * 1024 * 1024 * 1024 // 2GB minimum
	}

	err = SetMemoryLimit(client, testLimit)
	require.NoError(t, err, "SetMemoryLimit should not error (bug uses wrong config key 'max_memory')")

	got, err := GetMemoryLimit(client)
	require.NoError(t, err)
	assert.Equal(t, testLimit, got, "GetMemoryLimit should return the value we just set via SetMemoryLimit")
	_ = ctx
}
