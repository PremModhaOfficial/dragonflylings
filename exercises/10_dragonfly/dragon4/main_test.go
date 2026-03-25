package main

import (
	"context"
	"testing"

	"dragonflylings/lib/testutil"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWaitForReplication_StandaloneReturnsZero(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	// In Dragonfly standalone mode, WAIT returns 0 immediately.
	// This is correct behavior — there are no replicas.
	// The function should NOT treat this as an error.
	replicas, err := WaitForReplication(ctx, client, 1, 100)
	require.NoError(t, err,
		"WaitForReplication should not error when Dragonfly returns 0 replicas — "+
			"0 is the correct response for standalone mode, not a failure")
	assert.Equal(t, int64(0), replicas,
		"standalone Dragonfly has no replicas, so WAIT returns 0")
}

func TestWaitForReplication_ZeroReplicas(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	// Asking for 0 replicas should always succeed immediately
	replicas, err := WaitForReplication(ctx, client, 0, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(0), replicas)
}

func TestGetEncoding_String(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("dragon4-str")

	require.NoError(t, client.Set(ctx, key, "hello", 0).Err())

	// Should not panic — even if Dragonfly returns an encoding Redis doesn't use
	info, err := GetEncoding(ctx, client, key)
	require.NoError(t, err)
	assert.NotEmpty(t, info.Encoding, "encoding should not be empty")
	t.Logf("String encoding on Dragonfly: %q (IsKnown=%v)", info.Encoding, info.IsKnown)
}

func TestGetEncoding_Integer(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("dragon4-int")

	require.NoError(t, client.Set(ctx, key, "42", 0).Err())

	info, err := GetEncoding(ctx, client, key)
	require.NoError(t, err)
	assert.NotEmpty(t, info.Encoding)
	t.Logf("Integer encoding on Dragonfly: %q (IsKnown=%v)", info.Encoding, info.IsKnown)
}

func TestGetEncoding_Hash(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("dragon4-hash")

	require.NoError(t, client.HSet(ctx, key, "field", "value").Err())

	info, err := GetEncoding(ctx, client, key)
	require.NoError(t, err)
	assert.NotEmpty(t, info.Encoding)
	t.Logf("Hash encoding on Dragonfly: %q (IsKnown=%v)", info.Encoding, info.IsKnown)
}

func TestGetEncoding_UnknownDoesNotPanic(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("dragon4-unknown")

	// Create a sorted set — Dragonfly may return a non-standard encoding
	require.NoError(t, client.ZAdd(ctx, key, redis.Z{Score: 1.0, Member: "member"}).Err())

	// Must not panic — should return gracefully with IsKnown=false if encoding is non-standard
	require.NotPanics(t, func() {
		info, err := GetEncoding(ctx, client, key)
		require.NoError(t, err)
		assert.NotEmpty(t, info.Encoding)
		t.Logf("ZSet encoding on Dragonfly: %q (IsKnown=%v)", info.Encoding, info.IsKnown)
	})
}
