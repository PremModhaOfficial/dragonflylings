package main

import (
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRank(t *testing.T) {
	client := testutil.NewTestClient(t)
	lbKey := fmt.Sprintf("test:zsets3:lb:%d", time.Now().UnixNano())

	require.NoError(t, AddPlayer(client, lbKey, "alice", 1000))
	require.NoError(t, AddPlayer(client, lbKey, "bob", 2000))
	require.NoError(t, AddPlayer(client, lbKey, "charlie", 1500))

	// bob(2000) is rank 0, charlie(1500) is rank 1, alice(1000) is rank 2
	bobRank, err := GetRank(client, lbKey, "bob")
	require.NoError(t, err)
	assert.Equal(t, int64(0), bobRank, "bob has highest score, rank should be 0 (use ZRevRank)")

	aliceRank, err := GetRank(client, lbKey, "alice")
	require.NoError(t, err)
	assert.Equal(t, int64(2), aliceRank, "alice has lowest score, rank should be 2")
}

func TestIncrScore(t *testing.T) {
	client := testutil.NewTestClient(t)
	lbKey := fmt.Sprintf("test:zsets3:incr:%d", time.Now().UnixNano())

	require.NoError(t, AddPlayer(client, lbKey, "alice", 1000))

	// Increment by 500 — new score should be 1500
	newScore, err := IncrScore(client, lbKey, "alice", 500)
	require.NoError(t, err)
	assert.Equal(t, float64(1500), newScore,
		"IncrScore should ADD 500 to existing 1000 = 1500 (use ZIncrBy, not ZAdd)")
}

func TestRankChangesAfterIncr(t *testing.T) {
	client := testutil.NewTestClient(t)
	lbKey := fmt.Sprintf("test:zsets3:rankchange:%d", time.Now().UnixNano())

	require.NoError(t, AddPlayer(client, lbKey, "alice", 100))
	require.NoError(t, AddPlayer(client, lbKey, "bob", 200))

	aliceRank, _ := GetRank(client, lbKey, "alice")
	assert.Equal(t, int64(1), aliceRank, "alice starts at rank 1")

	// Give alice enough points to overtake bob
	_, err := IncrScore(client, lbKey, "alice", 200) // 100+200=300 > bob's 200
	require.NoError(t, err)

	aliceRankAfter, _ := GetRank(client, lbKey, "alice")
	assert.Equal(t, int64(0), aliceRankAfter, "after getting 200 points, alice should be rank 0")
}
