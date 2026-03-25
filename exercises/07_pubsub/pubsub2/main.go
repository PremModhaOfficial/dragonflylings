package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// EXERCISE: pubsub2 - Pattern Subscriptions
//
// PREDICT: Before writing any code, answer in your head:
//   What glob patterns does Redis support for PSUBSCRIBE? (*, ?, [...])
//   If you PSubscribe to "news.*", will "news.sports.local" match?
//   What type does Receive return for a pattern subscription message?
//
// TODO: Fix the two bugs below.

// SubscribePattern subscribes to all channels matching the given glob pattern.
// BUG: Uses Subscribe (exact match) instead of PSubscribe (pattern match).
func SubscribePattern(client *redis.Client, ctx context.Context, pattern string) *redis.PubSub {
	return client.Subscribe(ctx, pattern) // BUG: should be PSubscribe
}

// ReceivePatternMessage receives one message from a pattern subscription.
// Returns the channel name, pattern matched, and payload.
// BUG: Uses ReceiveMessage which returns *redis.Message (regular subscribe).
// For pattern subscriptions, use Receive and check the Pattern field.
func ReceivePatternMessage(sub *redis.PubSub, ctx context.Context) (channel, pattern, payload string, err error) {
	// BUG: For PSubscribe, messages have a non-empty Pattern field.
	// ReceiveMessage works, but it ignores the Pattern -- use Receive instead.
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		return "", "", "", err
	}
	return msg.Channel, "", msg.Payload, nil // BUG: Pattern always empty here
}

func main() {}
