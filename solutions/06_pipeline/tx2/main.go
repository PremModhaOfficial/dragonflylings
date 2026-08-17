package main

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ErrMaxRetries = errors.New("transaction failed after max retries")

func SafeTransfer(client *redis.Client, ctx context.Context, from, to string, amount int64, maxRetries int) error {
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
