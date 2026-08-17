package main

// EXERCISE: zsets3 - ZRANK and ZINCRBY
//
// PREDICT: Before fixing anything, answer:
//   - What's the difference between ZSCORE and ZRANK?
//   - What does ZRANK return for the member with the highest score?
//   - What does ZINCRBY do that ZADD NX doesn't?
//
// The test tracks a player's rank on a leaderboard and increases their score.
// BUG 1: GetRank uses ZScore (returns the score value) instead of ZRank (returns position).
// BUG 2: IncrScore uses ZAdd (overwrites score) instead of ZIncrBy (adds to current score).
//
// TODO: Fix GetRank to use ZRank/ZRevRank and IncrScore to use ZIncrBy.

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// AddPlayer adds a player with an initial score.
func AddPlayer(client *redis.Client, leaderboardKey, player string, score float64) error {
	ctx := context.Background()
	return client.ZAdd(ctx, leaderboardKey, redis.Z{Score: score, Member: player}).Err()
}

// GetRank returns the player's rank (0 = highest score, 1 = second highest, etc.).
// BUG: Uses ZScore — returns the score value, not the rank position.
func GetRank(client *redis.Client, leaderboardKey, player string) (int64, error) {
	ctx := context.Background()
	// TODO: use client.ZRevRank(ctx, leaderboardKey, player).Result()
	// ZRevRank gives rank from highest rank (rank 0 = best player)
	rank, err := client.ZRevRank(ctx, leaderboardKey, player).Result() // BUG: returns score, not rank
	return int64(rank), err
}

// IncrScore adds points to a player's existing score.
// BUG: Uses ZAdd which REPLACES the score instead of incrementing it.
func IncrScore(client *redis.Client, leaderboardKey, player string, points float64) (float64, error) {
	ctx := context.Background()
	// TODO: use client.ZIncrBy(ctx, leaderboardKey, points, player).Result()
	points, err := client.ZIncrBy(ctx, leaderboardKey, points, player).Result() // BUG: replaces score
	return points, err
}
