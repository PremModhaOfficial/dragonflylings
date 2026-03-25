package main

import (
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPagination(t *testing.T) {
	client := testutil.NewTestClient(t)
	listKey := fmt.Sprintf("test:lists2:activity:%d", time.Now().UnixNano())

	// Add 5 activities (LPush = most recent first in list)
	// After these pushes, list order is: event5, event4, event3, event2, event1
	for i := 1; i <= 5; i++ {
		require.NoError(t, AddActivity(client, listKey, fmt.Sprintf("event%d", i)))
	}

	// Page 1: indices 0-1 (2 per page), should be event5, event4 (most recent)
	page1, err := GetPage(client, listKey, 1, 2)
	require.NoError(t, err)
	require.Len(t, page1, 2)
	assert.Equal(t, "event5", page1[0], "page 1 should start at index 0 (0-based)")
	assert.Equal(t, "event4", page1[1])

	// Page 2: indices 2-3, should be event3, event2
	page2, err := GetPage(client, listKey, 2, 2)
	require.NoError(t, err)
	require.Len(t, page2, 2)
	assert.Equal(t, "event3", page2[0])
	assert.Equal(t, "event2", page2[1])
}

func TestGetAllWithNegativeIndex(t *testing.T) {
	client := testutil.NewTestClient(t)
	listKey := fmt.Sprintf("test:lists2:all:%d", time.Now().UnixNano())

	require.NoError(t, AddActivity(client, listKey, "a"))
	require.NoError(t, AddActivity(client, listKey, "b"))
	require.NoError(t, AddActivity(client, listKey, "c"))

	all, err := GetAll(client, listKey)
	require.NoError(t, err)
	// LPush means c, b, a order in list
	assert.Equal(t, []string{"c", "b", "a"}, all, "LRANGE 0 -1 should return all items")
}
