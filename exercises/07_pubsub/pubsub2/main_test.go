package main

import (
	"context"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscribePattern(t *testing.T) {
	subClient := testutil.NewTestClient(t)
	pubClient := testutil.NewTestClient(t)
	defer subClient.Close()
	defer pubClient.Close()

	prefix := testutil.UniqueKey("pubsub2")
	pattern := prefix + ".*"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sub := SubscribePattern(subClient, ctx, pattern)
	defer sub.Close()

	// Wait for subscription confirmation
	_, err := sub.Receive(ctx)
	require.NoError(t, err, "pattern subscription confirmation failed")

	// Publish to a channel that matches the pattern
	targetChannel := prefix + ".sports"
	go func() {
		pubClient.Publish(context.Background(), targetChannel, "goal!")
	}()

	ch, pat, payload, err := ReceivePatternMessage(sub, ctx)
	require.NoError(t, err, "ReceivePatternMessage failed (hint: pattern messages need Receive, not ReceiveMessage)")
	assert.Equal(t, targetChannel, ch)
	assert.Equal(t, pattern, pat, "pattern field should contain the matched glob pattern")
	assert.Equal(t, "goal!", payload)
}

func TestPatternMatchesMultipleChannels(t *testing.T) {
	subClient := testutil.NewTestClient(t)
	pubClient := testutil.NewTestClient(t)
	defer subClient.Close()
	defer pubClient.Close()

	prefix := testutil.UniqueKey("pubsub2m")
	pattern := prefix + ".*"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sub := SubscribePattern(subClient, ctx, pattern)
	defer sub.Close()

	_, err := sub.Receive(ctx)
	require.NoError(t, err)

	channels := []string{prefix + ".a", prefix + ".b", prefix + ".c"}
	go func() {
		for _, ch := range channels {
			pubClient.Publish(context.Background(), ch, "msg-"+ch)
		}
	}()

	received := make([]string, 0, len(channels))
	for i := 0; i < len(channels); i++ {
		_, _, payload, err := ReceivePatternMessage(sub, ctx)
		if err != nil {
			break
		}
		received = append(received, payload)
	}
	assert.Len(t, received, len(channels), "pattern subscription should receive messages from all matching channels")
}
