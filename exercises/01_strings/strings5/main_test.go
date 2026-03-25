package main

import (
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetAndGetPreferences(t *testing.T) {
	client := testutil.NewTestClient(t)
	prefix := fmt.Sprintf("test:strings5:%d", time.Now().UnixNano())

	prefs := map[string]string{
		prefix + ":theme":    "dark",
		prefix + ":lang":     "en",
		prefix + ":timezone": "UTC",
		prefix + ":pagesize": "20",
		prefix + ":notify":   "true",
	}

	err := SetPreferences(client, prefs)
	require.NoError(t, err)

	keys := []string{
		prefix + ":theme",
		prefix + ":lang",
		prefix + ":timezone",
		prefix + ":pagesize",
		prefix + ":notify",
	}

	results, err := GetPreferences(client, keys)
	require.NoError(t, err)
	require.Len(t, results, 5)

	assert.Equal(t, "dark", results[0])
	assert.Equal(t, "en", results[1])
	assert.Equal(t, "UTC", results[2])
	assert.Equal(t, "20", results[3])
	assert.Equal(t, "true", results[4])
}

func TestGetMissingPreference(t *testing.T) {
	client := testutil.NewTestClient(t)
	prefix := fmt.Sprintf("test:strings5:%d", time.Now().UnixNano())

	keys := []string{prefix + ":missing1", prefix + ":missing2"}
	results, err := GetPreferences(client, keys)
	require.NoError(t, err)
	require.Len(t, results, 2)

	// MGET returns nil for missing keys (not an error)
	assert.Nil(t, results[0], "missing key should return nil, not error")
	assert.Nil(t, results[1])
}
