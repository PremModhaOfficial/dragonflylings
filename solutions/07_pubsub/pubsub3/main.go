package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// PublishThenSubscribe publishes messages BEFORE subscribing to prove fire-and-forget.
// All messages are lost because they were published before any subscriber existed.
func PublishThenSubscribe(pubClient, subClient *redis.Client, channel string, messages []string) (received []string, missed int) {
	// Publish all messages first (no subscriber exists yet)
	for _, msg := range messages {
		pubClient.Publish(context.Background(), channel, msg)
	}

	// Now subscribe -- all messages already gone forever
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	sub := subClient.Subscribe(ctx, channel)
	defer sub.Close()

	// Try to receive -- nothing will arrive
	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			break // timeout: no messages came (they were already published)
		}
		received = append(received, msg.Payload)
	}

	missed = len(messages) - len(received)
	return received, missed
}

func main() {}
