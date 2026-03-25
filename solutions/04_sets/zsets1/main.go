package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func AddScore(client *redis.Client, leaderboardKey, player string, score float64) error {
	ctx := context.Background()
	return client.ZAdd(ctx, leaderboardKey, redis.Z{
		Score:  score,
		Member: player,
	}).Err()
}

func GetLeaderboard(client *redis.Client, leaderboardKey string, topN int64) ([]string, error) {
	ctx := context.Background()
	return client.ZRevRange(ctx, leaderboardKey, 0, topN-1).Result()
}

func GetScore(client *redis.Client, leaderboardKey, player string) (float64, error) {
	ctx := context.Background()
	return client.ZScore(ctx, leaderboardKey, player).Result()
}
