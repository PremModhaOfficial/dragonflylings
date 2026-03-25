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

func TestCircuitBreaker_ClosedNormalOperation(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("prod6")

	cb := NewCircuitBreaker(client, 3, time.Second)
	require.NoError(t, client.Set(ctx, key, "value", time.Minute).Err())

	val, err := cb.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, "value", val)
	assert.Equal(t, StateClosed, cb.GetState())
}

func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	cb := NewCircuitBreaker(client, 3, 5*time.Second)

	// Simulate failures using a non-existent key
	// redis.Nil is not a "failure" — we need actual errors
	// Use a bad context to force errors
	nonExistentKey := testutil.UniqueKey("prod6-nokey")
	for i := 0; i < 3; i++ {
		// Pre-cancel the context so the Get call always gets a real error
		// (not redis.Nil), guaranteeing recordFailure() is triggered.
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel()
		cb.Get(cancelCtx, nonExistentKey) //nolint:errcheck
	}

	// After 3 failures, circuit should be open
	assert.Equal(t, StateOpen, cb.GetState(),
		"circuit should be OPEN after %d failures; got %v — implement recordFailure() to track failures", 3, cb.GetState())
}

func TestCircuitBreaker_OpenRejectsCalls(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("prod6-open")

	cb := NewCircuitBreaker(client, 1, 10*time.Second)
	require.NoError(t, client.Set(ctx, key, "v", time.Minute).Err())

	// Force open by recording a failure directly
	cb.mu.Lock()
	cb.recordFailure()
	cb.mu.Unlock()

	require.Equal(t, StateOpen, cb.GetState())

	// All calls should return ErrCircuitOpen immediately
	_, err := cb.Get(ctx, key)
	assert.ErrorIs(t, err, ErrCircuitOpen,
		"open circuit should return ErrCircuitOpen immediately without calling Redis; got: %v", err)
}

func TestCircuitBreaker_HalfOpenAfterCooldown(t *testing.T) {
	client := testutil.NewTestClient(t)

	cb := NewCircuitBreaker(client, 1, 100*time.Millisecond)

	// Force open
	cb.mu.Lock()
	cb.recordFailure()
	cb.mu.Unlock()

	require.Equal(t, StateOpen, cb.GetState())

	// Wait for cooldown
	time.Sleep(150 * time.Millisecond)

	// Should be half-open now
	assert.Equal(t, StateHalfOpen, cb.GetState(),
		"circuit should transition to HALF-OPEN after cooldown; got %v", cb.GetState())
}

func TestCircuitBreaker_HalfOpenCloseOnSuccess(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("prod6-halfopen")

	cb := NewCircuitBreaker(client, 1, 100*time.Millisecond)
	require.NoError(t, client.Set(ctx, key, "recovered", time.Minute).Err())

	// Force open then wait for half-open
	cb.mu.Lock()
	cb.recordFailure()
	cb.mu.Unlock()

	time.Sleep(150 * time.Millisecond)
	require.Equal(t, StateHalfOpen, cb.GetState())

	// Successful probe call should close the circuit
	val, err := cb.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, "recovered", val)
	assert.Equal(t, StateClosed, cb.GetState(),
		"circuit should close after successful half-open probe; got %v", cb.GetState())
}

func TestCircuitBreaker_RedisNilNotFailure(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	cb := NewCircuitBreaker(client, 3, time.Second)
	nonExistent := testutil.UniqueKey("prod6-nil")

	// redis.Nil (key not found) should NOT count as a circuit breaker failure
	_, err := cb.Get(ctx, nonExistent)
	assert.ErrorIs(t, err, redis.Nil)
	assert.Equal(t, StateClosed, cb.GetState(),
		"redis.Nil (key not found) should not trip the circuit breaker; got state %v", cb.GetState())
}
