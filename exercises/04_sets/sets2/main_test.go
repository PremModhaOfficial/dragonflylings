package main

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupFollows(t *testing.T, client *redis.Client, prefix string) (string, string) {
	t.Helper()
	aliceKey := prefix + ":alice:follows"
	bobKey := prefix + ":bob:follows"

	// alice follows: bob, carol, dave
	_ = Follow(client, aliceKey, "bob")
	_ = Follow(client, aliceKey, "carol")
	_ = Follow(client, aliceKey, "dave")

	// bob follows: alice, carol, eve
	_ = Follow(client, bobKey, "alice")
	_ = Follow(client, bobKey, "carol")
	_ = Follow(client, bobKey, "eve")

	return aliceKey, bobKey
}

func TestCommonFollows(t *testing.T) {
	client := testutil.NewTestClient(t)
	prefix := fmt.Sprintf("test:sets2:%d", time.Now().UnixNano())
	aliceKey, bobKey := setupFollows(t, client, prefix)

	common, err := CommonFollows(client, aliceKey, bobKey)
	require.NoError(t, err)
	sort.Strings(common)
	assert.Equal(t, []string{"carol"}, common,
		"only 'carol' is followed by both alice and bob (SINTER, not SUNION)")
}

func TestAllFollows(t *testing.T) {
	client := testutil.NewTestClient(t)
	prefix := fmt.Sprintf("test:sets2:%d", time.Now().UnixNano())
	aliceKey, bobKey := setupFollows(t, client, prefix)

	all, err := AllFollows(client, aliceKey, bobKey)
	require.NoError(t, err)
	sort.Strings(all)
	assert.Equal(t, []string{"alice", "bob", "carol", "dave", "eve"}, all)
}

func TestUniqueFollows(t *testing.T) {
	client := testutil.NewTestClient(t)
	prefix := fmt.Sprintf("test:sets2:%d", time.Now().UnixNano())
	aliceKey, bobKey := setupFollows(t, client, prefix)

	// alice follows that bob doesn't
	unique, err := UniqueFollows(client, aliceKey, bobKey)
	require.NoError(t, err)
	sort.Strings(unique)
	assert.Equal(t, []string{"bob", "dave"}, unique,
		"bob and dave are followed by alice but not by bob")
}
