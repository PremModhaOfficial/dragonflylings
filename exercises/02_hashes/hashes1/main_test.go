package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetAndGetUserField(t *testing.T) {
	client := testutil.NewTestClient(t)
	userKey := fmt.Sprintf("test:hashes1:user:%d", time.Now().UnixNano())

	err := SetUserField(client, userKey, "name", "alice")
	require.NoError(t, err)

	err = SetUserField(client, userKey, "email", "alice@example.com")
	require.NoError(t, err)

	name, err := GetUserField(client, userKey, "name")
	require.NoError(t, err)
	assert.Equal(t, "alice", name)

	email, err := GetUserField(client, userKey, "email")
	require.NoError(t, err)
	assert.Equal(t, "alice@example.com", email)
}

func TestHashIsSingleKey(t *testing.T) {
	client := testutil.NewTestClient(t)
	userKey := fmt.Sprintf("test:hashes1:user:%d", time.Now().UnixNano())

	_ = SetUserField(client, userKey, "name", "bob")
	_ = SetUserField(client, userKey, "email", "bob@example.com")
	_ = SetUserField(client, userKey, "role", "admin")

	// The entire user should be stored under ONE key (the hash)
	keyType, err := client.Type(context.Background(), userKey).Result()
	require.NoError(t, err)
	assert.Equal(t, "hash", keyType,
		"user data should be stored as a hash type, not string; got type: %s", keyType)
}

func TestGetMissingHashField(t *testing.T) {
	client := testutil.NewTestClient(t)
	userKey := fmt.Sprintf("test:hashes1:user:%d", time.Now().UnixNano())

	// Set at least one field so the key exists
	_ = SetUserField(client, userKey, "name", "alice")

	_, err := GetUserField(client, userKey, "nonexistent")
	assert.ErrorIs(t, err, redis.Nil, "missing hash field returns redis.Nil")
}
