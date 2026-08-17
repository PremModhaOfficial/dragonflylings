package main

// EXERCISE: lists3 - BLPOP: Blocking Worker Pattern
//
// PREDICT: Before fixing anything, answer:
//   - What does LPOP return when the list is empty?
//   - What does BLPOP do differently when the list is empty?
//   - Why would you want a blocking pop in a worker process?
//
// The test starts a goroutine that pushes a job 200ms later, then calls WaitForJob.
// With non-blocking LPop: the queue is empty when called, returns error immediately — job missed.
// With blocking BLPop: waits up to the timeout, receives the job when it arrives.
//
// BUG: WaitForJob uses LPop (non-blocking) — it returns immediately when queue is empty.
// TODO: Replace LPop with BLPop using the provided timeout duration.

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// WaitForJob blocks until a job is available or timeout expires.
// Returns the job string, or an error if the queue is empty when timeout expires.
// BUG: Uses LPop (non-blocking) — returns redis.Nil immediately when queue is empty.
func WaitForJob(client *redis.Client, queueKey string, timeout time.Duration) (string, error) {
	ctx := context.Background()
	// TODO: use client.BLPop(ctx, timeout, queueKey).Result()
	// BLPop returns []string{key, value} or error after timeout
	val, err := client.BLPop(ctx, timeout, queueKey).Result() // BUG: non-blocking
	if err != nil {
		return "", err
	}
	return val[1], nil
}

// EnqueueJob pushes a job to the right end of the queue.
func EnqueueJob(client *redis.Client, queueKey, job string) error {
	ctx := context.Background()
	return client.RPush(ctx, queueKey, job).Err()
}
