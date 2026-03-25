package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func SubscribePattern(client *redis.Client, ctx context.Context, pattern string) *redis.PubSub {
	return client.PSubscribe(ctx, pattern)
}

func ReceivePatternMessage(sub *redis.PubSub, ctx context.Context) (channel, pattern, payload string, err error) {
	iface, err := sub.Receive(ctx)
	if err != nil {
		return "", "", "", err
	}
	switch msg := iface.(type) {
	case *redis.Message:
		// Pattern subscription messages have a non-empty Pattern field
		return msg.Channel, msg.Pattern, msg.Payload, nil
	case *redis.Subscription:
		// Skip subscription confirmations, try again
		return ReceivePatternMessage(sub, ctx)
	default:
		return "", "", "", fmt.Errorf("unexpected message type: %T", iface)
	}
}

func main() {}
