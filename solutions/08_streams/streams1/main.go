package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func AddEvent(client *redis.Client, ctx context.Context, stream string, fields map[string]interface{}) (string, error) {
	return client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		ID:     "*",
		Values: fields,
	}).Result()
}

func ReadAllEvents(client *redis.Client, ctx context.Context, stream string) ([]redis.XMessage, error) {
	results, err := client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{stream, "0"},
		Count:   100,
		Block:   -1, // non-blocking
	}).Result()
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0].Messages, nil
}

func main() {}
