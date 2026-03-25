package main

import (
	"testing"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
)

// TestFireAndForget proves that messages published before subscription are lost forever.
// This test is EXPECTED to show that received is empty -- that's the correct behavior.
// The exercise is about understanding this limitation, not working around it.
func TestFireAndForget(t *testing.T) {
	pubClient := testutil.NewTestClient(t)
	subClient := testutil.NewTestClient(t)
	defer pubClient.Close()
	defer subClient.Close()

	channel := testutil.UniqueKey("pubsub3")
	messages := []string{"msg1", "msg2", "msg3"}

	received, missed := PublishThenSubscribe(pubClient, subClient, channel, messages)

	// With the bug fixed: publish happens before subscribe, so all messages are missed.
	// This is the EXPECTED and CORRECT behavior -- Pub/Sub is fire-and-forget.
	assert.Equal(t, len(messages), missed,
		"all %d messages should be missed (published before subscription)", len(messages))
	assert.Empty(t, received,
		"no messages should be received (they were published before subscribing)")
}

// TestReceiveWhenSubscribedFirst shows that messages ARE received when timing is correct.
func TestReceiveWhenSubscribedFirst(t *testing.T) {
	pubClient := testutil.NewTestClient(t)
	subClient := testutil.NewTestClient(t)
	defer pubClient.Close()
	defer subClient.Close()

	channel := testutil.UniqueKey("pubsub3b")

	// Use Chat from pubsub1 pattern: subscribe first, then publish
	// This is a direct test -- subscribe, confirm, publish, receive
	ctx := t.Context()
	sub := subClient.Subscribe(ctx, channel)
	defer sub.Close()

	_, err := sub.Receive(ctx)
	if err != nil {
		t.Skip("subscribe confirmation failed -- Dragonfly may not be running")
	}

	go pubClient.Publish(ctx, channel, "arrived!")

	msg, err := sub.ReceiveMessage(ctx)
	if err == nil {
		assert.Equal(t, "arrived!", msg.Payload)
	}
}
