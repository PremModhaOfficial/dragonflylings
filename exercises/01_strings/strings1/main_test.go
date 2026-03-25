package main

import (
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetAndGet(t *testing.T) {
	client := testutil.NewTestClient(t)
	prefix := fmt.Sprintf("test:strings1:%d", time.Now().UnixNano())
	userID := prefix + ":user42"

	err := SetUsername(client, userID, "alice")
	require.NoError(t, err)

	got, err := GetUsername(client, userID)
	require.NoError(t, err)
	assert.Equal(t, "alice", got)
}

func TestGetMissingKey(t *testing.T) {
	client := testutil.NewTestClient(t)
	prefix := fmt.Sprintf("test:strings1:%d", time.Now().UnixNano())
	userID := prefix + ":nonexistent"

	_, err := GetUsername(client, userID)
	// go-redis returns redis.Nil for missing keys — not a regular error
	assert.ErrorIs(t, err, redis.Nil, "missing key should return redis.Nil")
}
