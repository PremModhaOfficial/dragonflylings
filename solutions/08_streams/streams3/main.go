package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func CreateGroup(client *redis.Client, ctx context.Context, stream, group string) error {
	return client.XGroupCreateMkStream(ctx, stream, group, "0").Err()
}

func ReadGroup(client *redis.Client, ctx context.Context, stream, group, consumer string, count int64) ([]redis.XMessage, error) {
	results, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{stream, ">"},
		Count:    count,
		Block:    0,
	}).Result()
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0].Messages, nil
}

func AckMessages(client *redis.Client, ctx context.Context, stream, group string, ids ...string) error {
	return client.XAck(ctx, stream, group, ids...).Err()
}

func main() {}
