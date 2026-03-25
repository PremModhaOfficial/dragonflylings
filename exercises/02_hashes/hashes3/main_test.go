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

func TestIncrPageView(t *testing.T) {
	client := testutil.NewTestClient(t)
	analyticsKey := fmt.Sprintf("test:hashes3:analytics:%d", time.Now().UnixNano())

	count, err := IncrPageView(client, analyticsKey, "home")
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	count, err = IncrPageView(client, analyticsKey, "home")
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	count, err = IncrPageView(client, analyticsKey, "about")
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	homeViews, err := GetPageViews(client, analyticsKey, "home")
	require.NoError(t, err)
	assert.Equal(t, int64(2), homeViews)
}

func TestAllViewsInOneHash(t *testing.T) {
	client := testutil.NewTestClient(t)
	analyticsKey := fmt.Sprintf("test:hashes3:analytics:%d", time.Now().UnixNano())

	_, _ = IncrPageView(client, analyticsKey, "home")
	_, _ = IncrPageView(client, analyticsKey, "about")
	_, _ = IncrPageView(client, analyticsKey, "contact")

	// All page views should be fields of ONE hash key
	keyType, err := client.Type(context.Background(), analyticsKey).Result()
	require.NoError(t, err)
	assert.Equal(t, "hash", keyType,
		"analytics data should be stored in a single hash, not separate string keys")

	all, err := GetAllPageViews(client, analyticsKey)
	require.NoError(t, err)
	assert.Len(t, all, 3, "hash should have one field per page")
}
