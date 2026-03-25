package main

// EXERCISE: strings2 - Expiration with TTL
//
// PREDICT: Before fixing anything, answer:
//   - What does SET EX do? How is it different from SET?
//   - What does TTL return for a key with no expiry? For a missing key?
//   - After a key expires, what does GET return?
//
// The test stores a session token that should expire in 2 seconds.
// BUG: Set is called with 0 duration (no expiry). The session lives forever.
//
// TODO: Fix SetSession to store the key with the provided TTL.

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// SetSession stores a session token with an expiration time.
// BUG: The expiration is hardcoded to 0 — keys never expire.
func SetSession(client *redis.Client, sessionID, token string, ttl time.Duration) error {
	ctx := context.Background()
	return client.Set(ctx, "session:"+sessionID, token, 0).Err() // TODO: use ttl instead of 0
}

// GetSession retrieves a session token.
func GetSession(client *redis.Client, sessionID string) (string, error) {
	ctx := context.Background()
	return client.Get(ctx, "session:"+sessionID).Result()
}

// GetTTL returns the remaining time-to-live for a session.
func GetTTL(client *redis.Client, sessionID string) (time.Duration, error) {
	ctx := context.Background()
	return client.TTL(ctx, "session:"+sessionID).Result()
}
