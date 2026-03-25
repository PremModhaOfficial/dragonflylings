package main

import (
	"sync"
	"testing"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSafeTransferSequential(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	from := testutil.UniqueKey("tx2:from")
	to := testutil.UniqueKey("tx2:to")

	require.NoError(t, client.MSet(ctx, from, 500, to, 0).Err())

	err := SafeTransfer(client, ctx, from, to, 200, 5)
	require.NoError(t, err)

	fromVal, _ := client.Get(ctx, from).Int64()
	toVal, _ := client.Get(ctx, to).Int64()
	assert.Equal(t, int64(300), fromVal)
	assert.Equal(t, int64(200), toVal)
}

func TestSafeTransferConcurrent(t *testing.T) {
	client := testutil.NewTestClient(t)
	defer client.Close()

	ctx := t.Context()
	from := testutil.UniqueKey("tx2:cfrom")
	to := testutil.UniqueKey("tx2:cto")

	const goroutines = 20
	const amount = int64(5)
	const initial = int64(goroutines) * amount

	require.NoError(t, client.MSet(ctx, from, initial, to, 0).Err())

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			SafeTransfer(client, ctx, from, to, amount, 20)
		}()
	}
	wg.Wait()

	fromVal, _ := client.Get(ctx, from).Int64()
	toVal, _ := client.Get(ctx, to).Int64()
	total := fromVal + toVal

	assert.Equal(t, initial, total,
		"total balance must be preserved under concurrent transfers (lost %d -- retry loop missing?)", initial-total)
	assert.Equal(t, int64(0), fromVal, "all funds should have transferred")
}
