package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Follow(client *redis.Client, followerKey, targetID string) error {
	ctx := context.Background()
	return client.SAdd(ctx, followerKey, targetID).Err()
}

func GetFollows(client *redis.Client, followerKey string) ([]string, error) {
	ctx := context.Background()
	return client.SMembers(ctx, followerKey).Result()
}

func CommonFollows(client *redis.Client, key1, key2 string) ([]string, error) {
	ctx := context.Background()
	return client.SInter(ctx, key1, key2).Result()
}

func AllFollows(client *redis.Client, key1, key2 string) ([]string, error) {
	ctx := context.Background()
	return client.SUnion(ctx, key1, key2).Result()
}

func UniqueFollows(client *redis.Client, key1, key2 string) ([]string, error) {
	ctx := context.Background()
	return client.SDiff(ctx, key1, key2).Result()
}
