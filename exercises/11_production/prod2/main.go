package main

// EXERCISE: prod2 - Distributed Lock
//
// PREDICT: Before writing any code, answer in your head:
//   SETNX sets a key only if it doesn't exist — it's the foundation of
//   a distributed lock. But what happens if the lock holder crashes before
//   calling Unlock? The lock is held forever. How do you prevent this?
//
//   Second question: if Unlock is GET + DEL (two commands), what happens
//   if another process acquires the lock between your GET and your DEL?
//   You'd delete someone else's lock. How do you prevent this?
//
// A production-grade distributed lock needs:
//   1. Expiry (TTL) so crashed holders don't block forever
//   2. A unique token per holder so only the holder can unlock
//   3. Atomic unlock (Lua script) to prevent check-then-delete races
//
// TODO: Fix TWO bugs:
//   Lock: add TTL so the lock auto-expires if the holder crashes
//   Unlock: make it atomic using a Lua check-and-delete script

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// unlockScript atomically checks the token and deletes the lock only if
// the current value matches the caller's token.
// KEYS[1] = lock key, ARGV[1] = expected token
// Returns: 1 if unlocked, 0 if token mismatch (someone else owns the lock)
const unlockScript = `
if redis.call('GET', KEYS[1]) == ARGV[1] then
  return redis.call('DEL', KEYS[1])
end
return 0
`

// Lock acquires a distributed lock using SET NX with a unique token.
// Returns true if the lock was acquired, false if it's already held.
//
// BUG: lock has no expiry (TTL=0). If the holder crashes, the lock
// is held forever — no other process can acquire it.
func Lock(ctx context.Context, client *redis.Client, lockKey, token string, ttl time.Duration) (bool, error) {
	// BUG: passes 0 as TTL — lock never expires
	return client.SetNX(ctx, lockKey, token, 0).Result()
}

// Unlock releases the lock only if the caller's token matches.
//
// BUG: not atomic — uses GET then DEL as two separate commands.
// Race condition: another process could acquire the lock between
// our GET and DEL, and we'd delete their lock.
func Unlock(ctx context.Context, client *redis.Client, lockKey, token string) error {
	val, err := client.Get(ctx, lockKey).Result()
	if err == redis.Nil {
		return nil // already expired — that's OK
	}
	if err != nil {
		return err
	}
	if val == token {
		// BUG: another process could SET lockKey between the GET above and this DEL
		return client.Del(ctx, lockKey).Err()
	}
	return nil
}
