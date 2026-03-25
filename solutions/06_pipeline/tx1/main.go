package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Transfer(client *redis.Client, ctx context.Context, from, to string, amount int64) error {
	_, err := client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.DecrBy(ctx, from, amount)
		pipe.IncrBy(ctx, to, amount)
		return nil
	})
	return err
}

func main() {}
