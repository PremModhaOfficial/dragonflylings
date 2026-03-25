package main

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/redis/go-redis/v9"
)

// SetPipelined sets all keys using a single pipeline — all commands are
// batched in one Exec(), reducing N round trips to 1.
func SetPipelined(ctx context.Context, client *redis.Client, keys []string, value string) (int64, error) {
	pipe := client.Pipeline()
	for _, key := range keys {
		pipe.Set(ctx, key, value, 0)
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	var count int64
	for _, cmd := range cmds {
		if cmd.Err() == nil {
			count++
		}
	}
	return count, nil
}

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

// SetConcurrent sets all keys using one goroutine per key.
// Demonstrates Dragonfly's parallel throughput vs sequential execution.
func SetConcurrent(ctx context.Context, client *redis.Client, keys []string, value string) int64 {
	var wg sync.WaitGroup
	var count int64

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			if err := client.Set(ctx, k, value, 0).Err(); err == nil {
				atomic.AddInt64(&count, 1)
			}
		}(key)
	}
	wg.Wait()
	return count
}
