package main

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

// EXERCISE: tx2 - Optimistic Locking with WATCH
//
// PREDICT: Before writing any code, answer in your head:
//   What problem does WATCH solve that MULTI/EXEC alone cannot?
//   What error does go-redis return when a WATCHed key changes before EXEC?
//   Why is this called "optimistic" locking?
//
// TODO: Add a retry loop. When TxFailedErr occurs, retry up to maxRetries times.

// ErrMaxRetries is returned when all retry attempts are exhausted.
var ErrMaxRetries = errors.New("transaction failed after max retries")

// SafeTransfer transfers amount from->to using WATCH for conflict detection.
// It should retry up to maxRetries times when concurrent modification is detected.
// BUG: Watch is used correctly, but TxFailedErr is returned immediately without retrying.
func SafeTransfer(client *redis.Client, ctx context.Context, from, to string, amount int64, maxRetries int) error {
	// BUG: no retry loop -- returns TxFailedErr immediately on conflict instead of retrying
	for range maxRetries {
		err := client.Watch(ctx, func(tx *redis.Tx) error {
			fromVal, err := tx.Get(ctx, from).Int64()
			if err != nil {
				return err
			}
			toVal, err := tx.Get(ctx, to).Int64()
			if err != nil {
				return err
			}

			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Set(ctx, from, fromVal-amount, 0)
				pipe.Set(ctx, to, toVal+amount, 0)
				return nil
			})
			return err
		}, from, to)
		if err == nil {
			return nil
		}
		if !errors.Is(err, redis.TxFailedErr) {
			return err
		}
	}

	return ErrMaxRetries
}

func main() {}
