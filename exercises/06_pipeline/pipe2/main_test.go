package main

import (
	"testing"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountPipelineFailures(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	key := testutil.UniqueKey("pipe2")
	testutil.CleanupKeys(t, client, key+"*")

	// Pre-set key to a non-integer so INCR fails
	require.NoError(t, client.Set(ctx, key, "not-a-number", 0).Err())

	failures, err := CountPipelineFailures(client, ctx, key)
	require.NoError(t, err)
	assert.Equal(t, 1, failures, "exactly 1 of the 3 pipeline commands should fail (the INCR on a non-integer)")

	// Verify the 2 successful commands actually ran (pipeline is NOT a transaction)
	assert.Equal(t, "good1", client.Get(ctx, key+":a").Val(), "successful pipeline commands should still execute")
	assert.Equal(t, "good2", client.Get(ctx, key+":b").Val(), "successful pipeline commands should still execute")
}
