package main

import (
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockingPop(t *testing.T) {
	client := testutil.NewTestClient(t)
	queueKey := fmt.Sprintf("test:lists3:jobs:%d", time.Now().UnixNano())

	// Push a job after 200ms delay (simulating a producer)
	go func() {
		time.Sleep(200 * time.Millisecond)
		_ = EnqueueJob(client, queueKey, "build-report")
	}()

	// WaitForJob should block until the job arrives (up to 2s)
	start := time.Now()
	job, err := WaitForJob(client, queueKey, 2*time.Second)
	elapsed := time.Since(start)

	require.NoError(t, err, "should receive the job without error")
	assert.Equal(t, "build-report", job, "should receive the job pushed by producer")
	assert.GreaterOrEqual(t, elapsed, 100*time.Millisecond,
		"blocking pop should wait for the job, not return immediately (non-blocking returns ~0ms)")
}

func TestBlockingPopTimeout(t *testing.T) {
	client := testutil.NewTestClient(t)
	queueKey := fmt.Sprintf("test:lists3:timeout:%d", time.Now().UnixNano())

	// Queue is empty — BLPop should return error after timeout
	start := time.Now()
	_, err := WaitForJob(client, queueKey, 300*time.Millisecond)
	elapsed := time.Since(start)

	require.Error(t, err, "empty queue with timeout should return error")
	assert.GreaterOrEqual(t, elapsed, 200*time.Millisecond,
		"should wait at least the timeout duration before erroring")
}
