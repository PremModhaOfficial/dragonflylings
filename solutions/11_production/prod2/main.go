package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// unlockScript atomically checks the token and deletes the lock only if
// the current value matches the caller's token.
// KEYS[1] = lock key, ARGV[1] = expected token
// Returns: 1 if unlocked, 0 if token mismatch
const unlockScript = `
if redis.call('GET', KEYS[1]) == ARGV[1] then
  return redis.call('DEL', KEYS[1])
end
return 0
`

// Lock acquires a distributed lock. Returns true if acquired.
// The lock auto-expires after ttl if the holder crashes.
func Lock(ctx context.Context, client *redis.Client, lockKey, token string, ttl time.Duration) (bool, error) {
	return client.SetNX(ctx, lockKey, token, ttl).Result()
}

// Unlock releases the lock atomically — only if the caller's token matches.
// Uses a Lua script to prevent the check-then-delete race condition.
func Unlock(ctx context.Context, client *redis.Client, lockKey, token string) error {
	_, err := client.Eval(ctx, unlockScript, []string{lockKey}, token).Int()
	if err == redis.Nil {
		return nil
	}
	return err
}
