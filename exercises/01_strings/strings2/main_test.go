package main

import (
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionExpires(t *testing.T) {
	client := testutil.NewTestClient(t)
	sessionID := fmt.Sprintf("test:strings2:%d", time.Now().UnixNano())

	err := SetSession(client, sessionID, "token-abc123", 2*time.Second)
	require.NoError(t, err)

	// Immediately retrievable
	token, err := GetSession(client, sessionID)
	require.NoError(t, err)
	assert.Equal(t, "token-abc123", token)

	// TTL should be positive (key has expiry set)
	ttl, err := GetTTL(client, sessionID)
	require.NoError(t, err)
	assert.Greater(t, ttl, time.Duration(0),
		"TTL should be > 0; got -1 means no expiry was set (the bug)")
	assert.LessOrEqual(t, ttl, 2*time.Second)
}

func TestSessionTTLSentinels(t *testing.T) {
	client := testutil.NewTestClient(t)
	prefix := fmt.Sprintf("test:strings2:%d", time.Now().UnixNano())

	// No-expiry key: TTL returns -1
	err := SetSession(client, prefix+":permanent", "tok", 0)
	require.NoError(t, err)
	ttl, _ := GetTTL(client, prefix+":permanent")
	// -1 means key exists but has no expiry
	assert.Equal(t, time.Duration(-1), ttl,
		"key with no expiry should return TTL=-1")
}
