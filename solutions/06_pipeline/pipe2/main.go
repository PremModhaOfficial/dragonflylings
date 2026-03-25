package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func CountPipelineFailures(client *redis.Client, ctx context.Context, key string) (int, error) {
	pipe := client.Pipeline()
	pipe.Set(ctx, key+":a", "good1", 0)
	pipe.Incr(ctx, key)
	pipe.Set(ctx, key+":b", "good2", 0)

	cmds, _ := pipe.Exec(ctx)

	failures := 0
	for _, cmd := range cmds {
		if cmd.Err() != nil {
			failures++
		}
	}
	return failures, nil
}

func main() {}
