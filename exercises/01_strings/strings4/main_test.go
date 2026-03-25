package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncrementBasic(t *testing.T) {
	client := testutil.NewTestClient(t)
	key := fmt.Sprintf("test:strings4:counter:%d", time.Now().UnixNano())

	val, err := IncrementCounter(client, key)
	require.NoError(t, err)
	assert.Equal(t, int64(1), val)

	val, err = IncrementCounter(client, key)
	require.NoError(t, err)
	assert.Equal(t, int64(2), val)
}

func TestIncrementConcurrent(t *testing.T) {
	client := testutil.NewTestClient(t)
	key := fmt.Sprintf("test:strings4:concurrent:%d", time.Now().UnixNano())

	const goroutines = 50
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			IncrementCounter(client, key)
		}()
	}
	wg.Wait()

	final, err := GetCounter(client, key)
	require.NoError(t, err)
	assert.Equal(t, int64(goroutines), final,
		"with atomic INCR, all %d increments must be counted; GET+SET loses updates under concurrency", goroutines)
}
