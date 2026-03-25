package main

import (
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLeaderboard(t *testing.T) {
	client := testutil.NewTestClient(t)
	lbKey := fmt.Sprintf("test:zsets1:lb:%d", time.Now().UnixNano())

	require.NoError(t, AddScore(client, lbKey, "alice", 1500))
	require.NoError(t, AddScore(client, lbKey, "bob", 2300))
	require.NoError(t, AddScore(client, lbKey, "charlie", 1800))
	require.NoError(t, AddScore(client, lbKey, "diana", 2800))
	require.NoError(t, AddScore(client, lbKey, "eve", 950))

	top3, err := GetLeaderboard(client, lbKey, 3)
	require.NoError(t, err)
	require.Len(t, top3, 3)

	// Highest scores first: diana(2800), bob(2300), charlie(1800)
	assert.Equal(t, "diana", top3[0], "highest score should be first (use ZRevRange)")
	assert.Equal(t, "bob", top3[1])
	assert.Equal(t, "charlie", top3[2])
}

func TestGetScore(t *testing.T) {
	client := testutil.NewTestClient(t)
	lbKey := fmt.Sprintf("test:zsets1:score:%d", time.Now().UnixNano())

	require.NoError(t, AddScore(client, lbKey, "alice", 1500))

	score, err := GetScore(client, lbKey, "alice")
	require.NoError(t, err)
	assert.Equal(t, float64(1500), score)
}

func TestUpdateScore(t *testing.T) {
	client := testutil.NewTestClient(t)
	lbKey := fmt.Sprintf("test:zsets1:update:%d", time.Now().UnixNano())

	require.NoError(t, AddScore(client, lbKey, "alice", 100))
	require.NoError(t, AddScore(client, lbKey, "alice", 500)) // update score

	score, err := GetScore(client, lbKey, "alice")
	require.NoError(t, err)
	assert.Equal(t, float64(500), score, "ZADD should update existing score")
}
