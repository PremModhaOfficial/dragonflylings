package main

import (
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcquireLock(t *testing.T) {
	client := testutil.NewTestClient(t)
	lockKey := fmt.Sprintf("test:strings3:lock:%d", time.Now().UnixNano())

	// First acquisition should succeed
	acquired, err := AcquireLock(client, lockKey, "worker-1", 5*time.Second)
	require.NoError(t, err)
	assert.True(t, acquired, "first lock acquisition should succeed")

	// Second acquisition should fail — lock is already held
	acquired2, err := AcquireLock(client, lockKey, "worker-2", 5*time.Second)
	require.NoError(t, err)
	assert.False(t, acquired2, "second acquisition should fail — lock is already held")
}

func TestReleaseLock(t *testing.T) {
	client := testutil.NewTestClient(t)
	lockKey := fmt.Sprintf("test:strings3:lock:%d", time.Now().UnixNano())

	// Acquire
	acquired, err := AcquireLock(client, lockKey, "worker-1", 5*time.Second)
	require.NoError(t, err)
	require.True(t, acquired)

	// Release
	err = ReleaseLock(client, lockKey)
	require.NoError(t, err)

	// Now another worker can acquire
	acquired2, err := AcquireLock(client, lockKey, "worker-2", 5*time.Second)
	require.NoError(t, err)
	assert.True(t, acquired2, "after release, lock should be available")
}
