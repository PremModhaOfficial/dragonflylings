package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func SetManyIndividual(client *redis.Client, ctx context.Context, prefix string, n int) error {
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("%s:%d", prefix, i)
		if err := client.Set(ctx, key, i, 0).Err(); err != nil {
			return err
		}
	}
	return nil
}

func SetManyPipelined(client *redis.Client, ctx context.Context, prefix string, n int) error {
	pipe := client.Pipeline()
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("%s:%d", prefix, i)
		pipe.Set(ctx, key, i, 0)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func main() {}
