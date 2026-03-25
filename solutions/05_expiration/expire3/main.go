package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// ExpiryChannel returns the pub/sub channel for keyevent expiry notifications in db 0.
func ExpiryChannel() string {
	return "__keyevent@0__:expired"
}

// WatchExpiredKeys subscribes to expiry notifications for db 0 and returns the subscription.
func WatchExpiredKeys(subClient *redis.Client, ctx context.Context) (*redis.PubSub, error) {
	channel := ExpiryChannel()
	sub := subClient.Subscribe(ctx, channel)
	if _, err := sub.Receive(ctx); err != nil {
		sub.Close()
		return nil, err
	}
	return sub, nil
}

func main() {}
