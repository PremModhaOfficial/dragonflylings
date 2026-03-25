package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func FindPendingMessages(client *redis.Client, ctx context.Context, stream, group string) ([]string, error) {
	entries, err := client.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: stream,
		Group:  group,
		Start:  "-",
		End:    "+",
		Count:  100,
	}).Result()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(entries))
	for _, e := range entries {
		ids = append(ids, e.ID)
	}
	return ids, nil
}

func ClaimMessages(client *redis.Client, ctx context.Context, stream, group, newConsumer string, msgIDs []string) ([]redis.XMessage, error) {
	if len(msgIDs) == 0 {
		return nil, nil
	}
	return client.XClaim(ctx, &redis.XClaimArgs{
		Stream:   stream,
		Group:    group,
		Consumer: newConsumer,
		MinIdle:  0,
		Messages: msgIDs,
	}).Result()
}

func main() {}
