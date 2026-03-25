package main

import (
	"context"
	"testing"
	"time"

	"dragonflylings/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetSession_Basic(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	s := Session{
		ID:   testutil.UniqueKey("sess"),
		Data: map[string]string{"user_id": "42", "role": "admin"},
	}

	err := SetSession(ctx, client, s)
	require.NoError(t, err)

	// Verify data stored correctly
	role, err := client.HGet(ctx, sessionKey(s.ID), "role").Result()
	require.NoError(t, err)
	assert.Equal(t, "admin", role)
}

func TestSetSession_HasTTL(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	s := Session{
		ID:   testutil.UniqueKey("sess-ttl"),
		Data: map[string]string{"user_id": "1"},
	}
	require.NoError(t, SetSession(ctx, client, s))

	ttl, err := client.TTL(ctx, sessionKey(s.ID)).Result()
	require.NoError(t, err)
	assert.Greater(t, ttl, time.Duration(0),
		"session must have a TTL; got %v — add client.Expire() after HSET, "+
			"otherwise sessions live forever and fill up memory", ttl)
}

func TestGetSession_Basic(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	s := Session{
		ID:   testutil.UniqueKey("sess-get"),
		Data: map[string]string{"user_id": "99", "email": "test@example.com"},
	}
	require.NoError(t, SetSession(ctx, client, s))

	got, err := GetSession(ctx, client, s.ID)
	require.NoError(t, err)
	assert.Equal(t, "99", got.Data["user_id"])
	assert.Equal(t, "test@example.com", got.Data["email"])
}

func TestGetSession_RefreshesTTL(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	s := Session{
		ID:   testutil.UniqueKey("sess-refresh"),
		Data: map[string]string{"user_id": "5"},
	}
	require.NoError(t, SetSession(ctx, client, s))

	// Manually set a short TTL to simulate a session near expiry
	require.NoError(t, client.Expire(ctx, sessionKey(s.ID), 2*time.Second).Err())

	ttlBefore, err := client.TTL(ctx, sessionKey(s.ID)).Result()
	require.NoError(t, err)
	require.Less(t, ttlBefore, 5*time.Second)

	// GetSession should refresh the TTL
	_, err = GetSession(ctx, client, s.ID)
	require.NoError(t, err)

	ttlAfter, err := client.TTL(ctx, sessionKey(s.ID)).Result()
	require.NoError(t, err)
	assert.Greater(t, ttlAfter, ttlBefore,
		"GetSession should refresh TTL on access; before=%v after=%v — "+
			"add client.Expire() in GetSession to implement sliding expiry", ttlBefore, ttlAfter)
}

func TestGetSession_MissingSession(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	// Session doesn't exist — should return empty session, not error
	s, err := GetSession(ctx, client, "nonexistent-session-id")
	require.NoError(t, err)
	assert.Empty(t, s.Data, "missing session should return empty data, not error")
}

func TestUpdateField_Basic(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	s := Session{
		ID:   testutil.UniqueKey("sess-update"),
		Data: map[string]string{"theme": "dark", "lang": "en"},
	}
	require.NoError(t, SetSession(ctx, client, s))

	updated, err := UpdateField(ctx, client, s.ID, "theme", "light")
	require.NoError(t, err)
	assert.True(t, updated)

	got, err := GetSession(ctx, client, s.ID)
	require.NoError(t, err)
	assert.Equal(t, "light", got.Data["theme"])
	assert.Equal(t, "en", got.Data["lang"], "other fields should be unchanged")
}

func TestDeleteSession(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()

	s := Session{ID: testutil.UniqueKey("sess-del"), Data: map[string]string{"x": "y"}}
	require.NoError(t, SetSession(ctx, client, s))

	require.NoError(t, DeleteSession(ctx, client, s.ID))

	exists, err := client.Exists(ctx, sessionKey(s.ID)).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists)
}
