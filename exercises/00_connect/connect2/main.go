package main

// EXERCISE: connect2 - Handle Unreachable Dragonfly
//
// PREDICT: Before fixing anything, answer:
//   - When you call redis.NewClient(), does it connect immediately?
//   - If Dragonfly is down, when does your code find out?
//   - What should a function return when it can't connect?
//
// The test expects:
//   1. Connect("localhost:6380") succeeds and returns a non-nil client
//   2. Connect("localhost:19999") returns nil client AND a non-nil error
//
// TODO: Fix the bug — the function never checks if Dragonfly is actually reachable.

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Connect attempts to connect to Dragonfly at addr.
// It should return an error if Dragonfly is unreachable.
// BUG: This function always returns (client, nil) — even when Dragonfly is down.
func Connect(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		DialTimeout: 500 * time.Millisecond,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, err
	}
	// TODO: verify the connection works — hint: use Ping with a context timeout
	// Right now we return the client even if Dragonfly is completely unreachable
	return client, nil
}
