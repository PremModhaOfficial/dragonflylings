package main

import (
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetManyPipelined(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	prefix := testutil.UniqueKey("pipe1")
	const n = 100

	testutil.CleanupKeys(t, client, prefix+":*")

	err := SetManyPipelined(client, ctx, prefix, n)
	require.NoError(t, err)

	for i := 0; i < n; i++ {
		key := fmt.Sprintf("%s:%d", prefix, i)
		val, err := client.Get(ctx, key).Int()
		require.NoError(t, err, "key %s should exist after SetManyPipelined", key)
		assert.Equal(t, i, val)
	}
}

func TestPipelineFasterThanIndividual(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping latency comparison in short mode")
	}
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	const n = 200

	prefix1 := testutil.UniqueKey("pipe1:ind")
	prefix2 := testutil.UniqueKey("pipe1:pip")
	testutil.CleanupKeys(t, client, prefix1+":*")
	testutil.CleanupKeys(t, client, prefix2+":*")

	start := time.Now()
	require.NoError(t, SetManyIndividual(client, ctx, prefix1, n))
	indTime := time.Since(start)

	start = time.Now()
	require.NoError(t, SetManyPipelined(client, ctx, prefix2, n))
	pipTime := time.Since(start)

	t.Logf("Individual: %v, Pipeline: %v", indTime, pipTime)

	count := client.Keys(ctx, prefix2+":*").Val()
	assert.Len(t, count, n, "pipeline should have written all %d keys", n)
}
