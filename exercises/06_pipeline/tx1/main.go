package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// EXERCISE: tx1 - The Atomic Transfer
//
// PREDICT: Before writing any code, answer in your head:
//   What does MULTI/EXEC guarantee that a plain pipeline does not?
//   Can a MULTI/EXEC transaction "roll back" if one command fails mid-execution?
//   If two goroutines both call TxPipelined at the same time, can they interleave?
//
// TODO: Fix the two swapped keys in Transfer.

// Transfer moves amount from the "from" balance key to the "to" balance key atomically.
// Both keys hold integer values representing balances.
// BUG: The DecrBy and IncrBy targets are swapped -- from is incremented and to is decremented.
func Transfer(client *redis.Client, ctx context.Context, from, to string, amount int64) error {
	_, err := client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.DecrBy(ctx, to, amount)   // BUG: should decrement "from"
		pipe.IncrBy(ctx, from, amount) // BUG: should increment "to"
		return nil
	})
	return err
}

func main() {}
