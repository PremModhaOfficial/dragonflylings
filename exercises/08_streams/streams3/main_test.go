package main

import (
	"testing"

	"dragonflylings/lib/testutil"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsumerGroup(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	stream := testutil.UniqueKey("streams3")
	group := testutil.UniqueKey("grp")
	consumer := "worker-1"
	testutil.CleanupKeys(t, client, stream)

	for i := 0; i < 5; i++ {
		_, err := client.XAdd(ctx, &redis.XAddArgs{
			Stream: stream,
			ID:     "*",
			Values: map[string]interface{}{"seq": i},
		}).Result()
		require.NoError(t, err)
	}

	err := CreateGroup(client, ctx, stream, group)
	require.NoError(t, err)

	msgs, err := ReadGroup(client, ctx, stream, group, consumer, 10)
	require.NoError(t, err)
	require.Len(t, msgs, 5, "ReadGroup should deliver all 5 pre-existing messages (fix CreateGroup to use '0' and ReadGroup to use '>')")

	ids := make([]string, len(msgs))
	for i, m := range msgs {
		ids[i] = m.ID
	}
	err = AckMessages(client, ctx, stream, group, ids...)
	require.NoError(t, err)

	pending, err := client.XPending(ctx, stream, group).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(0), pending.Count, "PEL should be empty after acknowledging all messages")
}
