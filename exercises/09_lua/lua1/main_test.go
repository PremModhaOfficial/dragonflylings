package main

import (
	"context"
	"testing"

	"dragonflylings/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompareAndSwap_Success(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("lua1")

	// Set initial value
	require.NoError(t, client.Set(ctx, key, "hello", 0).Err())

	// Swap should succeed: current "hello" matches expected "hello"
	swapped, err := CompareAndSwap(ctx, client, key, "hello", "world")
	require.NoError(t, err)
	assert.True(t, swapped, "CAS should succeed when expected value matches current")

	// Verify the value actually changed
	val, err := client.Get(ctx, key).Result()
	require.NoError(t, err)
	assert.Equal(t, "world", val, "key should hold the new value after successful CAS")
}

func TestCompareAndSwap_Failure(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("lua1-fail")

	require.NoError(t, client.Set(ctx, key, "world", 0).Err())

	// Swap should fail: current is "world", expected "hello" does NOT match
	swapped, err := CompareAndSwap(ctx, client, key, "hello", "universe")
	require.NoError(t, err)
	assert.False(t, swapped, "CAS should fail when expected value doesn't match current")

	// Value should be unchanged
	val, err := client.Get(ctx, key).Result()
	require.NoError(t, err)
	assert.Equal(t, "world", val, "value should be unchanged after failed CAS")
}

func TestCompareAndSwap_EmptyKey(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("lua1-empty")

	// Key does not exist — treat as empty string ""
	swapped, err := CompareAndSwap(ctx, client, key, "", "first")
	require.NoError(t, err)
	assert.True(t, swapped, "CAS on missing key with empty expected should succeed")

	val, err := client.Get(ctx, key).Result()
	require.NoError(t, err)
	assert.Equal(t, "first", val)
}

func TestCompareAndSwap_Idempotent(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("lua1-idempotent")

	require.NoError(t, client.Set(ctx, key, "v1", 0).Err())

	// First swap: v1 → v2
	swapped, err := CompareAndSwap(ctx, client, key, "v1", "v2")
	require.NoError(t, err)
	require.True(t, swapped)

	// Retry same swap: should fail (value is now v2, not v1)
	swapped, err = CompareAndSwap(ctx, client, key, "v1", "v2")
	require.NoError(t, err)
	assert.False(t, swapped, "second CAS with same expected should fail (value already changed)")
}
