package main

import (
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoundedList(t *testing.T) {
	client := testutil.NewTestClient(t)
	listKey := fmt.Sprintf("test:lists4:events:%d", time.Now().UnixNano())
	const maxItems = 5

	// Push 10 events
	for i := 1; i <= 10; i++ {
		err := AddEvent(client, listKey, fmt.Sprintf("event%d", i), maxItems)
		require.NoError(t, err)
	}

	count, err := EventCount(client, listKey)
	require.NoError(t, err)
	assert.Equal(t, int64(maxItems), count,
		"after pushing 10 events with max=5, list should have exactly 5 items (LTRIM missing?)")

	events, err := GetRecentEvents(client, listKey)
	require.NoError(t, err)
	require.Len(t, events, maxItems)

	// Most recent events should be kept (LPush = most recent at index 0)
	assert.Equal(t, "event10", events[0], "most recent event should be first")
	assert.Equal(t, "event6", events[4], "5th most recent event should be last")
}

func TestGetRecentEvents(t *testing.T) {
	client := testutil.NewTestClient(t)
	listKey := fmt.Sprintf("test:lists4:recent:%d", time.Now().UnixNano())

	require.NoError(t, AddEvent(client, listKey, "login", 3))
	require.NoError(t, AddEvent(client, listKey, "purchase", 3))
	require.NoError(t, AddEvent(client, listKey, "logout", 3))
	require.NoError(t, AddEvent(client, listKey, "login", 3)) // pushes out oldest

	events, err := GetRecentEvents(client, listKey)
	require.NoError(t, err)
	assert.Len(t, events, 3)
	assert.Equal(t, "login", events[0])    // most recent
	assert.Equal(t, "logout", events[1])
	assert.Equal(t, "purchase", events[2]) // purchase pushed out the first login
}
