package main

import (
	"testing"

	"dragonflylings/lib/testutil"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddEvent(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	stream := testutil.UniqueKey("streams1")
	testutil.CleanupKeys(t, client, stream)

	id1, err := AddEvent(client, ctx, stream, map[string]interface{}{"action": "login", "user": "alice"})
	require.NoError(t, err, "first AddEvent should succeed")
	assert.NotEmpty(t, id1, "ID should be auto-generated and non-empty")

	id2, err := AddEvent(client, ctx, stream, map[string]interface{}{"action": "logout", "user": "alice"})
	require.NoError(t, err, "second AddEvent should succeed (fails with ID=0 because IDs must increase)")
	assert.NotEmpty(t, id2)

	assert.NotEqual(t, id1, id2, "each entry should have a unique ID")
}

func TestReadAllEvents(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	stream := testutil.UniqueKey("streams1r")
	testutil.CleanupKeys(t, client, stream)

	for i := 0; i < 3; i++ {
		_, err := client.XAdd(ctx, &redis.XAddArgs{
			Stream: stream,
			ID:     "*",
			Values: map[string]interface{}{"seq": i},
		}).Result()
		require.NoError(t, err)
	}

	msgs, err := ReadAllEvents(client, ctx, stream)
	require.NoError(t, err)
	assert.Len(t, msgs, 3, "ReadAllEvents should return all 3 events (fails with '$' which means 'new only')")
}
