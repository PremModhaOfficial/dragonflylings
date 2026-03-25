package main

import (
	"context"
	"testing"

	"dragonflylings/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeAccountKeys_UseHashtag(t *testing.T) {
	balanceKey, reservedKey := MakeAccountKeys("42")

	// Both keys must contain a hashtag to ensure same-shard placement
	// Valid formats: "{42}:balance", "{account:42}:balance", etc.
	assert.Contains(t, balanceKey, "{", "balance key must use hashtag notation")
	assert.Contains(t, balanceKey, "}", "balance key must use hashtag notation")
	assert.Contains(t, reservedKey, "{", "reserved key must use hashtag notation")
	assert.Contains(t, reservedKey, "}", "reserved key must use hashtag notation")
}

func TestMakeAccountKeys_SameHashtag(t *testing.T) {
	balanceKey, reservedKey := MakeAccountKeys("99")

	// Extract the hashtag from each key — they must match
	balanceTag := extractHashtag(balanceKey)
	reservedTag := extractHashtag(reservedKey)

	assert.NotEmpty(t, balanceTag, "balance key must have a hashtag")
	assert.NotEmpty(t, reservedTag, "reserved key must have a hashtag")
	assert.Equal(t, balanceTag, reservedTag,
		"both keys must share the same hashtag so they land on the same shard")
}

func TestReserve_Sufficient(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	accountID := testutil.UniqueKey("lua4-acct")

	balanceKey, reservedKey := MakeAccountKeys(accountID)
	require.NoError(t, client.Set(ctx, balanceKey, "100", 0).Err())
	require.NoError(t, client.Set(ctx, reservedKey, "0", 0).Err())

	ok, err := Reserve(ctx, client, accountID, 30)
	require.NoError(t, err)
	assert.True(t, ok, "reservation should succeed with sufficient balance")

	balance, _ := client.Get(ctx, balanceKey).Int64()
	reserved, _ := client.Get(ctx, reservedKey).Int64()
	assert.Equal(t, int64(70), balance)
	assert.Equal(t, int64(30), reserved)
}

func TestReserve_Insufficient(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	accountID := testutil.UniqueKey("lua4-poor")

	balanceKey, reservedKey := MakeAccountKeys(accountID)
	require.NoError(t, client.Set(ctx, balanceKey, "20", 0).Err())
	require.NoError(t, client.Set(ctx, reservedKey, "0", 0).Err())

	ok, err := Reserve(ctx, client, accountID, 50)
	require.NoError(t, err)
	assert.False(t, ok, "reservation should fail with insufficient balance")

	// Neither key should have changed
	balance, _ := client.Get(ctx, balanceKey).Int64()
	reserved, _ := client.Get(ctx, reservedKey).Int64()
	assert.Equal(t, int64(20), balance, "balance unchanged after failed reserve")
	assert.Equal(t, int64(0), reserved, "reserved unchanged after failed reserve")
}

// extractHashtag extracts the content between { and } in a key name.
// Returns empty string if no hashtag found.
func extractHashtag(key string) string {
	start := -1
	for i, c := range key {
		if c == '{' {
			start = i
		} else if c == '}' && start >= 0 {
			return key[start : i+1]
		}
	}
	return ""
}
