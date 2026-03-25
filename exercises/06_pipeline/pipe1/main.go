package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// EXERCISE: pipe1 - The Assembly Line
//
// PREDICT: Before writing any code, answer in your head:
//   How many network round-trips does SetManyIndividual make for 100 keys?
//   How many round-trips does a pipeline make for 100 keys?
//   Does pipelining guarantee that all commands succeed together?
//
// TODO: Fix SetManyPipelined -- it builds a pipeline but never sends it.

// SetManyIndividual sets n keys one at a time. Reference implementation (already correct).
func SetManyIndividual(client *redis.Client, ctx context.Context, prefix string, n int) error {
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("%s:%d", prefix, i)
		if err := client.Set(ctx, key, i, 0).Err(); err != nil {
			return err
		}
	}
	return nil
}

// SetManyPipelined sets n keys using a single pipeline flush.
// BUG: The pipeline is built but Exec is never called -- keys are never sent to Dragonfly.
func SetManyPipelined(client *redis.Client, ctx context.Context, prefix string, n int) error {
	pipe := client.Pipeline()
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("%s:%d", prefix, i)
		pipe.Set(ctx, key, i, 0)
	}
	// BUG: Missing pipe.Exec(ctx)
	return nil
}

func main() {}
