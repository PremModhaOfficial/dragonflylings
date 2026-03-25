package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func AddAndTrim(client *redis.Client, ctx context.Context, stream string, fields map[string]interface{}, maxLen int64) error {
	if err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		ID:     "*",
		Values: fields,
	}).Err(); err != nil {
		return err
	}
	// Use exact trim for predictable test behavior.
	// In high-throughput production, XTrimMaxLenApprox is more efficient.
	return client.XTrimMaxLen(ctx, stream, maxLen).Err()
}

func main() {}
