package main

import (
	"testing"

	"dragonflylings/lib/testutil"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCrashRecovery(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	stream := testutil.UniqueKey("streams4")
	group := testutil.UniqueKey("grp4")
	crashedConsumer := "worker-crashed"
	recoveryConsumer := "worker-recovery"
	testutil.CleanupKeys(t, client, stream)

	require.NoError(t, client.XGroupCreateMkStream(ctx, stream, group, "0").Err())

	for i := 0; i < 3; i++ {
		client.XAdd(ctx, &redis.XAddArgs{
			Stream: stream,
			ID:     "*",
			Values: map[string]interface{}{"job": i},
		})
	}

	// Simulate crashed consumer: read but never ACK
	client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: crashedConsumer,
		Streams:  []string{stream, ">"},
		Count:    10,
	})

	pendingIDs, err := FindPendingMessages(client, ctx, stream, group)
	require.NoError(t, err)
	require.Len(t, pendingIDs, 3, "FindPendingMessages should return 3 pending IDs (bug uses XPending summary which has no IDs)")

	claimed, err := ClaimMessages(client, ctx, stream, group, recoveryConsumer, pendingIDs)
	require.NoError(t, err)
	assert.Len(t, claimed, 3, "ClaimMessages should claim all 3 messages (bug uses 1hr minIdle -- too long)")

	for _, msg := range claimed {
		assert.NotEmpty(t, msg.ID)
	}
}
