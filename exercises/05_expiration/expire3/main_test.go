package main

import (
	"context"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpiryNotification(t *testing.T) {
	client := testutil.NewTestClient(t)
	subClient := testutil.NewTestClient(t)
	defer client.Close()
	defer subClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Subscribe to expiry notifications using the function under test
	sub, err := WatchExpiredKeys(subClient, ctx)
	require.NoError(t, err, "WatchExpiredKeys failed -- check the channel name (db index and keyspace vs keyevent)")
	defer sub.Close()

	// Set a key with a very short TTL
	key := testutil.UniqueKey("expire3")
	err = client.Set(ctx, key, "will-expire", 150*time.Millisecond).Err()
	require.NoError(t, err)

	// Wait for expiry notification
	received := make(chan string, 1)
	go func() {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			return
		}
		received <- msg.Payload
	}()

	select {
	case payload := <-received:
		assert.Equal(t, key, payload, "notification payload should be the expired key name")
	case <-ctx.Done():
		t.Fatal("timeout: no expiry notification received (check ExpiryChannel name: keyevent not keyspace, db 0 not db 1)")
	}
}

func TestExpiryChannelName(t *testing.T) {
	channel := ExpiryChannel()
	assert.Equal(t, "__keyevent@0__:expired", channel,
		"ExpiryChannel should return the keyevent channel for db 0")
}
