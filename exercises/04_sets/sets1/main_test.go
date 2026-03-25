package main

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddAndGetTags(t *testing.T) {
	client := testutil.NewTestClient(t)
	itemKey := fmt.Sprintf("test:sets1:item:%d", time.Now().UnixNano())

	require.NoError(t, AddTag(client, itemKey, "go"))
	require.NoError(t, AddTag(client, itemKey, "redis"))
	require.NoError(t, AddTag(client, itemKey, "database"))

	tags, err := GetTags(client, itemKey)
	require.NoError(t, err)
	sort.Strings(tags)
	assert.Equal(t, []string{"database", "go", "redis"}, tags)
}

func TestNoDuplicateTags(t *testing.T) {
	client := testutil.NewTestClient(t)
	itemKey := fmt.Sprintf("test:sets1:dedup:%d", time.Now().UnixNano())

	// Add "go" tag three times — should only appear once
	require.NoError(t, AddTag(client, itemKey, "go"))
	require.NoError(t, AddTag(client, itemKey, "go"))
	require.NoError(t, AddTag(client, itemKey, "go"))

	tags, err := GetTags(client, itemKey)
	require.NoError(t, err)
	assert.Len(t, tags, 1, "duplicate tags should be deduplicated (use SADD, not LPUSH)")
}

func TestHasTag(t *testing.T) {
	client := testutil.NewTestClient(t)
	itemKey := fmt.Sprintf("test:sets1:has:%d", time.Now().UnixNano())

	require.NoError(t, AddTag(client, itemKey, "redis"))

	has, err := HasTag(client, itemKey, "redis")
	require.NoError(t, err)
	assert.True(t, has)

	missing, err := HasTag(client, itemKey, "mysql")
	require.NoError(t, err)
	assert.False(t, missing)
}
