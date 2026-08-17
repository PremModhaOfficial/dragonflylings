package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// EXERCISE: pubsub3 - Fire and Forget
//
// PREDICT: Before writing any code, answer in your head:
//   What happens to a message published when zero subscribers are connected?
//   If a subscriber disconnects and reconnects, does it receive missed messages?
//   How is Pub/Sub fundamentally different from Streams in terms of durability?
//
// TODO: Fix PublishThenSubscribe to publish BEFORE subscribing.
// The test proves fire-and-forget: when you publish before subscribing, messages are lost.
// BUG: The current code subscribes FIRST, then publishes -- so messages ARE received.
// The test expects zero received messages. Fix: move Publish to happen before Subscribe.

// PublishThenSubscribe should publish messages first (before subscribing),
// then subscribe and try to collect them -- proving they are lost forever.
// BUG: subscribes before publishing, so messages arrive and test fails.
func PublishThenSubscribe(pubClient, subClient *redis.Client, channel string, messages []string) (received []string, missed int) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// BUG: subscribing BEFORE publishing -- messages will be received
	// Fix: move these lines to AFTER the Publish loop below

	// Publish messages (currently happens AFTER subscription -- messages arrive)
	for _, msg := range messages {
		pubClient.Publish(context.Background(), channel, msg)
	}
	sub := subClient.Subscribe(ctx, channel)
	defer sub.Close()
	sub.Receive(ctx) //nolint:errcheck -- wait for subscription confirmation

	// Collect messages
	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			break
		}
		received = append(received, msg.Payload)
	}

	missed = len(messages) - len(received)
	return received, missed
}

func main() {}
