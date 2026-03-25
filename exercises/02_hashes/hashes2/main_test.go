package main

import (
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMultipleFields(t *testing.T) {
	client := testutil.NewTestClient(t)
	userKey := fmt.Sprintf("test:hashes2:user:%d", time.Now().UnixNano())

	err := SetUserProfile(client, userKey, map[string]interface{}{
		"name":  "charlie",
		"email": "charlie@example.com",
		"role":  "editor",
		"age":   "29",
	})
	require.NoError(t, err)

	results, err := GetUserFields(client, userKey, []string{"name", "email", "role"})
	require.NoError(t, err)
	require.Len(t, results, 3)

	assert.Equal(t, "charlie", results[0])
	assert.Equal(t, "charlie@example.com", results[1])
	assert.Equal(t, "editor", results[2])
}

func TestGetMissingFieldReturnsNil(t *testing.T) {
	client := testutil.NewTestClient(t)
	userKey := fmt.Sprintf("test:hashes2:user:%d", time.Now().UnixNano())

	err := SetUserProfile(client, userKey, map[string]interface{}{
		"name": "dave",
	})
	require.NoError(t, err)

	// "avatar" field doesn't exist — HMGET returns nil, not an error
	results, err := GetUserFields(client, userKey, []string{"name", "avatar"})
	require.NoError(t, err, "HMGet should not error on missing fields")
	require.Len(t, results, 2)
	assert.Equal(t, "dave", results[0])
	assert.Nil(t, results[1], "missing hash field should be nil in results")
}
