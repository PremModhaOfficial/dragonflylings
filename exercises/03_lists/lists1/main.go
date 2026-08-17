package main

// EXERCISE: lists1 - LPUSH/RPUSH/LPOP/RPOP
//
// PREDICT: Before fixing anything, answer:
//   - What's the difference between LPUSH and RPUSH?
//   - What does LPUSH + LPOP give you? What does LPUSH + RPOP give you?
//   - Which combination implements a FIFO queue? Which is a LIFO stack?
//
// The test submits task1, task2, task3 and expects them processed in that order (FIFO).
// BUG: EnqueueTask uses LPush (prepends to the LEFT). DequeueTask uses LPop (pops from LEFT).
//      LPush+LPop = LIFO stack: task3 is dequeued first, not task1.
//
// TODO: Fix so tasks are dequeued FIFO (first submitted = first processed).
//       Option A: change EnqueueTask to RPush (append to right), keep DequeueTask as LPop.
//       Option B: keep EnqueueTask as LPush, change DequeueTask to RPop.

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// EnqueueTask adds a task to the queue.
// BUG: LPush prepends — each new task goes to the front (LIFO if popped from same end).
func EnqueueTask(client *redis.Client, queueKey, task string) error {
	ctx := context.Background()
	return client.LPush(ctx, queueKey, task).Err() // TODO: use RPush for FIFO
}

// DequeueTask removes and returns the next task from the queue.
func DequeueTask(client *redis.Client, queueKey string) (string, error) {
	ctx := context.Background()
	return client.RPop(ctx, queueKey).Result()
}

// QueueLength returns the number of tasks waiting.
func QueueLength(client *redis.Client, queueKey string) (int64, error) {
	ctx := context.Background()
	return client.LLen(ctx, queueKey).Result()
}
