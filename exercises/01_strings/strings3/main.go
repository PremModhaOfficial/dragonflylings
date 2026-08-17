package main

// EXERCISE: strings3 - SETNX (Set If Not Exists)
//
// PREDICT: Before fixing anything, answer:
//   - What does SETNX return when the key already exists?
//   - What does SETNX return when the key doesn't exist?
//   - Why is "set only if not exists" the foundation of distributed locking?
//
// The test tries to "acquire a lock" by setting a key only if absent.
// BUG: AcquireLock uses Set (always overwrites) instead of SetNX (conditional).
//
// TODO: Change Set to SetNX and return the boolean result.

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// AcquireLock attempts to acquire a distributed lock.
// Returns true if the lock was acquired, false if already held.
// BUG: Set always overwrites — it never returns false even if lock is held.
func AcquireLock(client *redis.Client, lockKey, ownerID string, ttl time.Duration) (bool, error) {
	ctx := context.Background()
	// TODO: use SetNX instead of Set
	// SetNX returns true if key was set (lock acquired), false if key existed (lock held)
	res, err := client.SetNX(ctx, lockKey, ownerID, ttl).Result()
	if err != nil || errors.Is(err, redis.Nil) {
		return false, err
	}
	return res, nil // BUG: always returns true — should return SetNX result
}

// ReleaseLock deletes the lock key.
func ReleaseLock(client *redis.Client, lockKey string) error {
	ctx := context.Background()
	return client.Del(ctx, lockKey).Err()
}
