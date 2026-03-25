package main

import (
	"context"
	"testing"
	"time"

	"dragonflylings/lib/testutil"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUsedMemory_ReturnsBytes(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	mem, err := GetUsedMemory(ctx, client)
	require.NoError(t, err)

	// used_memory should be at least 1 MB (Dragonfly's baseline overhead)
	assert.Greater(t, mem, int64(1024*1024),
		"used_memory should be >1MB (baseline overhead); got %d bytes — "+
			"if this is a small number like 1-100, you're probably returning DBSIZE (key count) instead of bytes", mem)

	// And less than 10 GB (sanity check)
	assert.Less(t, mem, int64(10*1024*1024*1024),
		"used_memory should be <10GB; something seems wrong")
}

func TestGetUsedMemory_IncreasesAfterWrites(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	before, err := GetUsedMemory(ctx, client)
	require.NoError(t, err)

	// Write 1000 keys with 100-byte values
	prefix := testutil.UniqueKey("dragon2")
	testutil.CleanupKeys(t, client, prefix+":*")
	pipe := client.Pipeline()
	for i := 0; i < 1000; i++ {
		key := testutil.UniqueKey(prefix)
		pipe.Set(ctx, key, string(make([]byte, 100)), 0)
	}
	_, err = pipe.Exec(ctx)
	require.NoError(t, err)

	after, err := GetUsedMemory(ctx, client)
	require.NoError(t, err)

	assert.Greater(t, after, before,
		"memory should increase after writing data; before=%d after=%d", before, after)
}

func TestWaitForSnapshot_Completes(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := TriggerSnapshot(ctx, client)
	// Dragonfly may return "ERR" if a snapshot is already running — that's OK
	if err != nil && err != redis.Nil {
		// Check if it's a "already running" error — acceptable
		t.Logf("TriggerSnapshot returned: %v (may be OK if snapshot already running)", err)
	}

	// WaitForSnapshot should return without error once complete
	err = WaitForSnapshot(ctx, client)
	require.NoError(t, err)

	// Verify no snapshot is in progress after Wait returns.
	// Dragonfly uses "saving" as its primary indicator; fall back to
	// "rdb_bgsave_in_progress" for Redis compat deployments.
	info, err := client.Info(ctx, "persistence").Result()
	require.NoError(t, err)
	if saving, ok := parseInfoField(info, "saving"); ok {
		assert.Equal(t, "0", saving,
			"saving should be 0 after WaitForSnapshot returns")
	} else {
		val, _ := parseInfoField(info, "rdb_bgsave_in_progress")
		assert.Equal(t, "0", val,
			"rdb_bgsave_in_progress should be 0 after WaitForSnapshot returns")
	}
}

func TestParseInfoField(t *testing.T) {
	info := "# Memory\r\nused_memory:2048000\r\nused_memory_human:1.95M\r\n"
	val, ok := parseInfoField(info, "used_memory")
	assert.True(t, ok)
	assert.Equal(t, "2048000", val)
}
