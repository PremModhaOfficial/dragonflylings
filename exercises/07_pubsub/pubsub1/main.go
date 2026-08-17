package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// EXERCISE: pubsub1 - The Chat Room
//
// PREDICT: Before writing any code, answer in your head:
//   If nobody is subscribed when Publish is called, what happens to the message?
//   Why must Subscribe run in a separate goroutine from the publisher?
//   What is the first message you receive after calling Subscribe?
//
// TODO: Fix Chat -- the context is cancelled immediately, preventing any messages from arriving.

// Chat subscribes to channel, receives n messages, and returns them.
// A publisher goroutine sends the messages after subscribing.
// BUG: cancel() is called immediately after creating the context, so the context is
// already Done before Subscribe can receive anything.
func Chat(subClient *redis.Client, pubClient *redis.Client, channel string, messages []string) ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // BUG: immediately cancels -- subscription and receive will fail

	sub := subClient.Subscribe(ctx, channel)
	defer sub.Close()

	// Wait for subscription confirmation
	if _, err := sub.Receive(ctx); err != nil {
		return nil, fmt.Errorf("subscribe failed: %w", err)
	}

	// Publish messages after subscription is confirmed
	go func() {
		for _, msg := range messages {
			pubClient.Publish(context.Background(), channel, msg)
		}
	}()

	received := make([]string, 0, len(messages))
	for range messages {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			return received, nil
		}
		received = append(received, msg.Payload)
	}
	return received, nil
}

func main() {}
