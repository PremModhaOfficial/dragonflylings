package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func AddTag(client *redis.Client, itemKey, tag string) error {
	ctx := context.Background()
	return client.SAdd(ctx, itemKey+":tags", tag).Err()
}

func GetTags(client *redis.Client, itemKey string) ([]string, error) {
	ctx := context.Background()
	return client.SMembers(ctx, itemKey+":tags").Result()
}

func HasTag(client *redis.Client, itemKey, tag string) (bool, error) {
	ctx := context.Background()
	return client.SIsMember(ctx, itemKey+":tags", tag).Result()
}
