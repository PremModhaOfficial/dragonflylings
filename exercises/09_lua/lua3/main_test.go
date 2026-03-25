package main

import (
	"context"
	"testing"

	"dragonflylings/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadScript_ReturnsSHA(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	sha, err := LoadScript(ctx, client)
	require.NoError(t, err)
	// SHA1 hashes are 40 hex characters
	assert.Len(t, sha, 40, "LoadScript should return a 40-character SHA1 hash")
	assert.NotEmpty(t, sha, "SHA1 should not be empty")
}

func TestDecrIfPositive_WithSHA(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("lua3")

	// Load script first
	sha, err := LoadScript(ctx, client)
	require.NoError(t, err)
	require.Len(t, sha, 40, "need valid SHA to test EvalSha")

	// Set counter to 2
	require.NoError(t, client.Set(ctx, key, "2", 0).Err())

	// First decrement: 2 → 1
	decremented, err := DecrIfPositive(ctx, client, sha, key)
	require.NoError(t, err)
	assert.True(t, decremented, "should decrement when counter > 0")

	val, _ := client.Get(ctx, key).Int64()
	assert.Equal(t, int64(1), val)

	// Second decrement: 1 → 0
	decremented, err = DecrIfPositive(ctx, client, sha, key)
	require.NoError(t, err)
	assert.True(t, decremented, "should decrement when counter is 1")

	val, _ = client.Get(ctx, key).Int64()
	assert.Equal(t, int64(0), val)

	// Third attempt: at zero, should NOT decrement
	decremented, err = DecrIfPositive(ctx, client, sha, key)
	require.NoError(t, err)
	assert.False(t, decremented, "should NOT decrement when counter is already 0")

	val, _ = client.Get(ctx, key).Int64()
	assert.Equal(t, int64(0), val, "counter should remain at 0")
}

func TestDecrIfPositive_MissingKey(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("lua3-missing")

	sha, err := LoadScript(ctx, client)
	require.NoError(t, err)

	// Key doesn't exist — treated as 0, should not decrement
	decremented, err := DecrIfPositive(ctx, client, sha, key)
	require.NoError(t, err)
	assert.False(t, decremented, "should not decrement non-existent key (treated as 0)")
}

func TestDecrIfPositive_ScriptCached(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	// Load script twice — should get same SHA (deterministic)
	sha1, err := LoadScript(ctx, client)
	require.NoError(t, err)

	sha2, err := LoadScript(ctx, client)
	require.NoError(t, err)

	assert.Equal(t, sha1, sha2, "same script should always produce same SHA1")
}
