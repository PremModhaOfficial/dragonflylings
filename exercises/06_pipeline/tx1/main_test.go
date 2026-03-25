package main

import (
	"testing"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	from := testutil.UniqueKey("tx1:alice")
	to := testutil.UniqueKey("tx1:bob")
	testutil.CleanupKeys(t, client, "tx1:*")

	require.NoError(t, client.MSet(ctx, from, 1000, to, 0).Err())

	err := Transfer(client, ctx, from, to, 300)
	require.NoError(t, err)

	fromVal, err := client.Get(ctx, from).Int64()
	require.NoError(t, err)
	toVal, err := client.Get(ctx, to).Int64()
	require.NoError(t, err)

	assert.Equal(t, int64(700), fromVal, "sender should have 1000-300=700 after transfer")
	assert.Equal(t, int64(300), toVal, "receiver should have 0+300=300 after transfer")
	assert.Equal(t, int64(1000), fromVal+toVal, "total balance must be conserved")
}
