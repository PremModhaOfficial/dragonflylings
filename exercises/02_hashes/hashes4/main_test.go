package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreAndGetUser(t *testing.T) {
	client := testutil.NewTestClient(t)
	userKey := fmt.Sprintf("test:hashes4:user:%d", time.Now().UnixNano())

	data := map[string]string{
		"name":  "eve",
		"email": "eve@example.com",
		"role":  "viewer",
		"age":   "31",
	}

	err := StoreUser(client, userKey, data)
	require.NoError(t, err)

	got, err := GetUser(client, userKey)
	require.NoError(t, err)
	assert.Equal(t, "eve", got["name"])
	assert.Equal(t, "eve@example.com", got["email"])
	assert.Equal(t, "viewer", got["role"])
	assert.Equal(t, "31", got["age"])
}

func TestUserIsStoredAsHash(t *testing.T) {
	client := testutil.NewTestClient(t)
	userKey := fmt.Sprintf("test:hashes4:user:%d", time.Now().UnixNano())

	_ = StoreUser(client, userKey, map[string]string{
		"name": "frank",
		"age":  "40",
	})

	// User data should be stored as a single hash, not scattered string keys
	keyType, err := client.Type(context.Background(), userKey).Result()
	require.NoError(t, err)
	assert.Equal(t, "hash", keyType,
		"user should be a hash type; storing as separate string keys is the bug")
}

func TestDeleteUserRemovesOneKey(t *testing.T) {
	client := testutil.NewTestClient(t)
	userKey := fmt.Sprintf("test:hashes4:user:%d", time.Now().UnixNano())

	_ = StoreUser(client, userKey, map[string]string{"name": "grace", "email": "g@example.com"})
	err := DeleteUser(client, userKey)
	require.NoError(t, err)

	// After deletion, userKey should not exist
	exists, err := client.Exists(context.Background(), userKey).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists, "after DeleteUser, the hash key should not exist")
}
