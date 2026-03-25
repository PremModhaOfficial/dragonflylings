package main

import (
	"testing"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddAndTrim(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	stream := testutil.UniqueKey("streams5")
	testutil.CleanupKeys(t, client, stream)

	const maxLen = int64(5)
	const total = 10

	for i := 0; i < total; i++ {
		err := AddAndTrim(client, ctx, stream, map[string]interface{}{"seq": i}, maxLen)
		require.NoError(t, err)
	}

	count, err := client.XLen(ctx, stream).Result()
	require.NoError(t, err)

	assert.Equal(t, maxLen, count,
		"stream should have exactly %d entries after adding %d with maxLen=%d (bug uses 0 which deletes everything)",
		maxLen, total, maxLen)
}

func TestStreamStaysBounded(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	stream := testutil.UniqueKey("streams5b")
	testutil.CleanupKeys(t, client, stream)

	const maxLen = int64(3)

	for i := 0; i < 20; i++ {
		require.NoError(t, AddAndTrim(client, ctx, stream, map[string]interface{}{"i": i}, maxLen))
	}

	count, _ := client.XLen(ctx, stream).Result()
	assert.LessOrEqual(t, count, maxLen, "stream should never exceed maxLen entries")
}
