package main

// EXERCISE: dragon1 - Multi-Threaded Throughput
//
// PREDICT: Before writing any code, answer in your head:
//   Dragonfly processes commands on different keys in parallel.
//   If you send 1000 SET commands sequentially (one at a time), how does
//   that differ from sending them concurrently from 10 goroutines?
//   What's the bottleneck in each case?
//
// Dragonfly's multi-threaded architecture means concurrent clients can
// achieve much higher throughput than sequential requests. This exercise
// demonstrates the difference.
//
// TODO: Fix SetConcurrent — it launches goroutines but never waits for them
// to finish (missing wg.Done() and wg.Wait()), so it returns before any
// work is done.
//
// TODO: Fix SetPipelined — it creates a new pipeline per key instead of
// batching all commands in one pipeline. This gives zero pipelining benefit
// (N round trips instead of 1). Move the pipeline outside the loop.

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/redis/go-redis/v9"
)

// SetSequential sets all keys one at a time. Returns count of successful sets.
func SetSequential(ctx context.Context, client *redis.Client, keys []string, value string) int64 {
	var count int64
	for _, key := range keys {
		if err := client.Set(ctx, key, value, 0).Err(); err == nil {
			count++
		}
	}
	return count
}

// SetPipelined sets all keys using Redis pipelining to reduce round trips.
// Pipelining lets you send multiple commands without waiting for each reply —
// the client buffers them and flushes once, the server replies in order.
//
// BUG: creates a new pipeline per key — N pipelines with 1 command each
// means N round trips, giving zero batching benefit over sequential SET.
// Fix: create one pipeline outside the loop, queue all commands, Exec once.
func SetPipelined(ctx context.Context, client *redis.Client, keys []string, value string) (int64, error) {
	var count int64
	for _, key := range keys {
		// BUG: new pipeline per key — defeats pipelining entirely
		pipe := client.Pipeline()
		pipe.Set(ctx, key, value, 0)
		cmds, err := pipe.Exec(ctx)
		if err != nil {
			return count, err
		}
		if cmds[0].Err() == nil {
			count++
		}
	}
	return count, nil
}

// SetConcurrent sets all keys using one goroutine per key.
// BUG 1: goroutines are launched but wg.Done() is never called
// BUG 2: wg.Wait() is missing — function returns before goroutines finish
// This means count will always be 0 (or far less than len(keys)).
func SetConcurrent(ctx context.Context, client *redis.Client, keys []string, value string) int64 {
	var wg sync.WaitGroup
	var count int64

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			// BUG: missing wg.Done() — WaitGroup never decrements
			defer wg.Done()
			if err := client.Set(ctx, k, value, 0).Err(); err == nil {
				atomic.AddInt64(&count, 1)
			}
		}(key)
	}
	// BUG: missing wg.Wait() — returns immediately before goroutines complete
	wg.Wait()
	return count
}
