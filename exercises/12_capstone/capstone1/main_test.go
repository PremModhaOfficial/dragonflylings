package main

import (
	"context"
	"testing"

	"dragonflylings/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventPipeline_PublishAndProcess(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	gameID := testutil.UniqueKey("capstone1")
	p := NewEventPipeline(client, gameID)

	// Publish events BEFORE creating group (tests Bug 1: start ID must be "0")
	_, err := p.PublishScore(ctx, "alice", 100)
	require.NoError(t, err)
	_, err = p.PublishScore(ctx, "bob", 75)
	require.NoError(t, err)
	_, err = p.PublishScore(ctx, "alice", 50)
	require.NoError(t, err)

	// Create consumer group
	require.NoError(t, p.CreateGroup(ctx))

	// Process all events
	count, err := p.ProcessPending(ctx)
	require.NoError(t, err)
	assert.Equal(t, 3, count, "should process all 3 events")
}

func TestEventPipeline_LeaderboardCorrect(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	gameID := testutil.UniqueKey("capstone1-lb")
	p := NewEventPipeline(client, gameID)

	// Publish score events
	require.NoError(t, func() error {
		_, err := p.PublishScore(ctx, "alice", 100)
		if err != nil {
			return err
		}
		_, err = p.PublishScore(ctx, "bob", 200)
		if err != nil {
			return err
		}
		_, err = p.PublishScore(ctx, "alice", 50)
		return err
	}())

	require.NoError(t, p.CreateGroup(ctx))

	_, err := p.ProcessPending(ctx)
	require.NoError(t, err)

	top, err := p.GetTopPlayers(ctx, 3)
	require.NoError(t, err)
	require.Len(t, top, 2, "should have 2 players in leaderboard")

	assert.Equal(t, "bob", top[0].Member, "bob should be #1 (200 points)")
	assert.Equal(t, float64(200), top[0].Score)

	assert.Equal(t, "alice", top[1].Member, "alice should be #2 (150 points total: 100+50)")
	assert.Equal(t, float64(150), top[1].Score,
		"alice's scores should be summed: 100+50=150 — Lua ZINCRBY adds to existing score")
}

func TestEventPipeline_MessagesAcknowledged(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	gameID := testutil.UniqueKey("capstone1-ack")
	p := NewEventPipeline(client, gameID)

	_, err := p.PublishScore(ctx, "carol", 50)
	require.NoError(t, err)

	require.NoError(t, p.CreateGroup(ctx))

	_, err = p.ProcessPending(ctx)
	require.NoError(t, err)

	// After processing, PEL should be empty (all messages acknowledged)
	pending, err := p.GetPendingCount(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), pending,
		"PEL should be empty after processing — add XACK after each processed message; "+
			"unacknowledged messages accumulate in PEL and get redelivered on restart")
}

func TestEventPipeline_IdempotentGroupCreate(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	gameID := testutil.UniqueKey("capstone1-idem")
	p := NewEventPipeline(client, gameID)

	// Publishing first to ensure stream exists
	_, _ = p.PublishScore(ctx, "player", 1)

	// Creating group twice should not error
	require.NoError(t, p.CreateGroup(ctx))
	require.NoError(t, p.CreateGroup(ctx), "creating group twice should be idempotent")
}

func TestEventPipeline_LuaAtomicity(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	gameID := testutil.UniqueKey("capstone1-lua")
	p := NewEventPipeline(client, gameID)

	// Publish multiple events for same player
	for i := 0; i < 5; i++ {
		_, err := p.PublishScore(ctx, "dave", 10)
		require.NoError(t, err)
	}

	require.NoError(t, p.CreateGroup(ctx))
	count, err := p.ProcessPending(ctx)
	require.NoError(t, err)
	assert.Equal(t, 5, count)

	top, err := p.GetTopPlayers(ctx, 1)
	require.NoError(t, err)
	require.Len(t, top, 1)
	assert.Equal(t, float64(50), top[0].Score,
		"5 × 10 points = 50 total — Lua ZINCRBY must add, not set")
}
