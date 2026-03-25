package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func AddPlayer(client *redis.Client, leaderboardKey, player string, score float64) error {
	ctx := context.Background()
	return client.ZAdd(ctx, leaderboardKey, redis.Z{Score: score, Member: player}).Err()
}

func GetRank(client *redis.Client, leaderboardKey, player string) (int64, error) {
	ctx := context.Background()
	return client.ZRevRank(ctx, leaderboardKey, player).Result()
}

func IncrScore(client *redis.Client, leaderboardKey, player string, points float64) (float64, error) {
	ctx := context.Background()
	return client.ZIncrBy(ctx, leaderboardKey, points, player).Result()
}
