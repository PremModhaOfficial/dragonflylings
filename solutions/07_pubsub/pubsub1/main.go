package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func Chat(subClient *redis.Client, pubClient *redis.Client, channel string, messages []string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sub := subClient.Subscribe(ctx, channel)
	defer sub.Close()

	if _, err := sub.Receive(ctx); err != nil {
		return nil, fmt.Errorf("subscribe failed: %w", err)
	}

	go func() {
		for _, msg := range messages {
			pubClient.Publish(context.Background(), channel, msg)
		}
	}()

	received := make([]string, 0, len(messages))
	for i := 0; i < len(messages); i++ {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			return received, nil
		}
		received = append(received, msg.Payload)
	}
	return received, nil
}

func main() {}
