package main

import (
	"testing"

	"dragonflylings/lib/testutil"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountEvents(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	stream := testutil.UniqueKey("streams2")
	testutil.CleanupKeys(t, client, stream)

	for i := 0; i < 5; i++ {
		client.XAdd(ctx, &redis.XAddArgs{Stream: stream, ID: "*", Values: map[string]interface{}{"i": i}})
	}

	count, err := CountEvents(client, ctx, stream)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count, "CountEvents should return 5 (bug appends :count suffix to key name)")
}

func TestQueryRange(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	stream := testutil.UniqueKey("streams2r")
	testutil.CleanupKeys(t, client, stream)

	for i := 0; i < 5; i++ {
		client.XAdd(ctx, &redis.XAddArgs{Stream: stream, ID: "*", Values: map[string]interface{}{"i": i}})
	}

	// Query all entries using - (min) to + (max)
	msgs, err := QueryRange(client, ctx, stream, "-", "+")
	require.NoError(t, err)
	assert.Len(t, msgs, 5, "QueryRange with start='-' stop='+' should return all 5 entries (reversed bounds return nothing)")

	for i := 1; i < len(msgs); i++ {
		assert.Greater(t, msgs[i].ID, msgs[i-1].ID, "entries should be in ascending ID order")
	}
}
