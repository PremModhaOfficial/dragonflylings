package main

import (
	"context"
	"testing"

	"dragonflylings/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAtomicTransfer_Basic(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	fromKey := testutil.UniqueKey("lua2-from")
	toKey := testutil.UniqueKey("lua2-to")

	// Initialize balances
	require.NoError(t, client.Set(ctx, fromKey, "100", 0).Err())
	require.NoError(t, client.Set(ctx, toKey, "50", 0).Err())

	// Transfer 30 from → to
	err := AtomicTransfer(ctx, client, fromKey, toKey, 30)
	require.NoError(t, err)

	fromVal, err := client.Get(ctx, fromKey).Int64()
	require.NoError(t, err)
	assert.Equal(t, int64(70), fromVal, "from balance should decrease by amount")

	toVal, err := client.Get(ctx, toKey).Int64()
	require.NoError(t, err)
	assert.Equal(t, int64(80), toVal, "to balance should increase by amount")
}

func TestAtomicTransfer_InsufficientBalance(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	fromKey := testutil.UniqueKey("lua2-poor")
	toKey := testutil.UniqueKey("lua2-rich")

	require.NoError(t, client.Set(ctx, fromKey, "10", 0).Err())
	require.NoError(t, client.Set(ctx, toKey, "0", 0).Err())

	// Attempt to transfer more than available
	err := AtomicTransfer(ctx, client, fromKey, toKey, 50)
	assert.Error(t, err, "should error on insufficient balance")

	// Both values should be unchanged
	fromVal, _ := client.Get(ctx, fromKey).Int64()
	toVal, _ := client.Get(ctx, toKey).Int64()
	assert.Equal(t, int64(10), fromVal, "from balance should be unchanged after failed transfer")
	assert.Equal(t, int64(0), toVal, "to balance should be unchanged after failed transfer")
}

func TestAtomicTransfer_ZeroBalance(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	fromKey := testutil.UniqueKey("lua2-zero")
	toKey := testutil.UniqueKey("lua2-zero2")

	// from key does not exist (treated as 0)
	err := AtomicTransfer(ctx, client, fromKey, toKey, 1)
	assert.Error(t, err, "transferring from non-existent key should fail")
}
