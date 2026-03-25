package main

import (
	"fmt"
	"testing"
	"time"

	"dragonflylings/lib/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueueFIFO(t *testing.T) {
	client := testutil.NewTestClient(t)
	queueKey := fmt.Sprintf("test:lists1:queue:%d", time.Now().UnixNano())

	// Enqueue in order
	require.NoError(t, EnqueueTask(client, queueKey, "task1"))
	require.NoError(t, EnqueueTask(client, queueKey, "task2"))
	require.NoError(t, EnqueueTask(client, queueKey, "task3"))

	length, err := QueueLength(client, queueKey)
	require.NoError(t, err)
	assert.Equal(t, int64(3), length)

	// Dequeue should give us FIFO order: task1, task2, task3
	task, err := DequeueTask(client, queueKey)
	require.NoError(t, err)
	assert.Equal(t, "task1", task, "first task submitted should be first dequeued (FIFO)")

	task, err = DequeueTask(client, queueKey)
	require.NoError(t, err)
	assert.Equal(t, "task2", task)

	task, err = DequeueTask(client, queueKey)
	require.NoError(t, err)
	assert.Equal(t, "task3", task)
}

func TestQueueEmptyReturnsError(t *testing.T) {
	client := testutil.NewTestClient(t)
	queueKey := fmt.Sprintf("test:lists1:empty:%d", time.Now().UnixNano())

	_, err := DequeueTask(client, queueKey)
	require.Error(t, err, "dequeuing from empty queue should return error")
}
