package main

// EXERCISE: connect1 - Your First PING
//
// PREDICT: Before writing any code, answer in your head:
//   - What network protocol does Redis/Dragonfly use?
//   - What do you think PING returns? Why would that command exist?
//   - What's the difference between creating a client and connecting?
//
// The test expects:
//   1. Connect() returns a client pointing at localhost:6380
//   2. Ping() returns "PONG" with no error
//
// TODO: Fix the two bugs below so the tests pass.

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Connect returns a new Redis client.
// BUG: The port is wrong — Dragonfly runs on 6380, not 6379.
func Connect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6380", // TODO: fix the port
	})
}

// Ping sends a health-check to Dragonfly and returns the response.
// BUG: Echo sends the ECHO command (echoes a string back), not PING.
func Ping(client *redis.Client) (string, error) {
	ctx := context.Background()
	return client.Ping(ctx).Result()
}
