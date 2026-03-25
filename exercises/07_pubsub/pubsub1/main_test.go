package main

import (
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChat(t *testing.T) {
	subClient := testutil.NewTestClient(t)
	pubClient := testutil.NewTestClient(t)
	defer subClient.Close()
	defer pubClient.Close()

	channel := testutil.UniqueKey("pubsub1")
	messages := []string{"hello", "world", "goodbye"}

	done := make(chan []string, 1)
	errCh := make(chan error, 1)

	go func() {
		received, err := Chat(subClient, pubClient, channel, messages)
		if err != nil {
			errCh <- err
			return
		}
		done <- received
	}()

	select {
	case received := <-done:
		require.Len(t, received, len(messages), "should receive all %d messages", len(messages))
		assert.Equal(t, messages, received)
	case err := <-errCh:
		t.Fatalf("Chat returned error: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout: Chat did not complete in 5s (context cancelled too early?)")
	}
}
