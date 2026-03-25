package main

import (
	"context"
	"testing"
	"time"

	"dragonflylings/lib/testutil"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLock_Acquire(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	lockKey := testutil.UniqueKey("prod2-lock")

	acquired, err := Lock(ctx, client, lockKey, "token-1", 5*time.Second)
	require.NoError(t, err)
	assert.True(t, acquired, "first lock acquisition should succeed")
}

func TestLock_HasTTL(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	lockKey := testutil.UniqueKey("prod2-ttl")

	acquired, err := Lock(ctx, client, lockKey, "token-1", 5*time.Second)
	require.NoError(t, err)
	require.True(t, acquired)

	// Verify the lock has a TTL set (not -1 which means no expiry)
	ttl, err := client.TTL(ctx, lockKey).Result()
	require.NoError(t, err)
	assert.Greater(t, ttl, time.Duration(0),
		"lock must have a positive TTL to auto-expire if holder crashes; got TTL=%v "+
			"(a TTL of -1 means no expiry — the lock will be held forever if the holder crashes)", ttl)
	assert.LessOrEqual(t, ttl, 5*time.Second, "TTL should be ≤ 5s")
}

func TestLock_Exclusive(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	lockKey := testutil.UniqueKey("prod2-exclusive")

	acquired1, err := Lock(ctx, client, lockKey, "token-1", 10*time.Second)
	require.NoError(t, err)
	require.True(t, acquired1)

	// Second acquisition should fail while lock is held
	acquired2, err := Lock(ctx, client, lockKey, "token-2", 10*time.Second)
	require.NoError(t, err)
	assert.False(t, acquired2, "second lock acquisition should fail while lock is held")
}

func TestUnlock_CorrectToken(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	lockKey := testutil.UniqueKey("prod2-unlock")

	acquired, err := Lock(ctx, client, lockKey, "my-token", 10*time.Second)
	require.NoError(t, err)
	require.True(t, acquired)

	err = Unlock(ctx, client, lockKey, "my-token")
	require.NoError(t, err)

	// Lock should be released — someone else can acquire it now
	acquired2, err := Lock(ctx, client, lockKey, "other-token", 10*time.Second)
	require.NoError(t, err)
	assert.True(t, acquired2, "lock should be re-acquirable after unlock")
}

func TestUnlock_WrongToken(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	lockKey := testutil.UniqueKey("prod2-wrongtoken")

	// Process A acquires the lock
	acquired, err := Lock(ctx, client, lockKey, "token-A", 10*time.Second)
	require.NoError(t, err)
	require.True(t, acquired)

	// Process B tries to unlock with wrong token — should be a no-op
	err = Unlock(ctx, client, lockKey, "token-B")
	require.NoError(t, err)

	// Lock should still be held by A
	val, err := client.Get(ctx, lockKey).Result()
	require.NoError(t, err)
	assert.Equal(t, "token-A", val, "lock should still be held by A after B's failed unlock")
}

func TestLock_ReacquireAfterExpiry(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	lockKey := testutil.UniqueKey("prod2-expire")

	// Acquire with short TTL
	acquired, err := Lock(ctx, client, lockKey, "token-1", 100*time.Millisecond)
	require.NoError(t, err)
	require.True(t, acquired)

	// Wait for TTL to expire
	time.Sleep(200 * time.Millisecond)

	// Verify the lock expired
	exists, err := client.Exists(ctx, lockKey).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists, "lock should auto-expire after TTL")

	// Should be re-acquirable now
	acquired2, err := Lock(ctx, client, lockKey, "token-2", 5*time.Second)
	require.NoError(t, err)
	assert.True(t, acquired2, "lock should be re-acquirable after TTL expiry")
}

func TestUnlock_AlreadyExpired(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	lockKey := testutil.UniqueKey("prod2-expired")

	// Key doesn't exist — unlock should be a no-op, not an error
	err := Unlock(ctx, client, lockKey, "token-1")
	assert.NoError(t, err, "unlocking a non-existent key should not error")
}

// TestUnlock_AtomicCheck verifies the Lua script constant is well-formed
func TestUnlock_LuaScriptExists(t *testing.T) {
	assert.NotEmpty(t, unlockScript, "unlockScript constant should be defined")
	assert.Contains(t, unlockScript, "redis.call", "unlockScript should contain redis.call")
	_ = redis.NewClient // ensure redis package is used
}
