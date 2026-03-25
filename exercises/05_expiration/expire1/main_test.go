package main

import (
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetWithExpiry(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	key := testutil.UniqueKey("expire1:set")
	err := SetWithExpiry(client, key, "hello", 10*time.Second)
	require.NoError(t, err)

	ttl := client.TTL(t.Context(), key)
	require.NoError(t, ttl.Err())
	assert.Greater(t, ttl.Val(), time.Duration(0), "key should have a positive TTL after SetWithExpiry")
}

func TestHasExpiry(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	key := testutil.UniqueKey("expire1:has")

	// Key with expiry
	client.Set(ctx, key, "value", 30*time.Second)
	got, err := HasExpiry(client, key)
	require.NoError(t, err)
	assert.True(t, got, "HasExpiry should return true for a key with TTL set")

	// Key without expiry
	persistent := testutil.UniqueKey("expire1:persist")
	client.Set(ctx, persistent, "value", 0)
	got2, err := HasExpiry(client, persistent)
	require.NoError(t, err)
	assert.False(t, got2, "HasExpiry should return false for a persistent key")
}

func TestMakePersistent(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	key := testutil.UniqueKey("expire1:mk")

	// Set key with expiry
	client.Set(ctx, key, "important", 5*time.Second)

	err := MakePersistent(client, key)
	require.NoError(t, err)

	// Key must still exist
	val, err := client.Get(ctx, key).Result()
	require.NoError(t, err, "key should still exist after MakePersistent (not be deleted!)")
	assert.Equal(t, "important", val)

	// TTL must now be -1 (no expiry)
	ttl := client.TTL(ctx, key).Val()
	assert.Equal(t, time.Duration(-1), ttl, "TTL should be -1 (persistent) after MakePersistent")
}
