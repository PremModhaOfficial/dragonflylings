package main

// EXERCISE: zsets1 - ZADD/ZRANGE/ZSCORE (Leaderboard)
//
// PREDICT: Before fixing anything, answer:
//   - What's the difference between a Set and a Sorted Set (ZSet)?
//   - What does ZRANGE return by default — ascending or descending?
//   - Can two members have the same score in a Sorted Set?
//
// The test builds a game leaderboard sorted by score.
// BUG: AddScore uses ZAdd correctly, but GetLeaderboard uses ZRANGE with
//      wrong direction — it returns lowest scores first instead of highest first.
//
// TODO: Fix GetLeaderboard to return players sorted by score descending (highest first).

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// AddScore adds or updates a player's score on the leaderboard.
func AddScore(client *redis.Client, leaderboardKey, player string, score float64) error {
	ctx := context.Background()
	return client.ZAdd(ctx, leaderboardKey, redis.Z{
		Score:  score,
		Member: player,
	}).Err()
}

// GetLeaderboard returns the top N players by score (highest first).
// BUG: ZRange returns members in ASCENDING score order (lowest first).
//      A leaderboard should show the highest scores first.
func GetLeaderboard(client *redis.Client, leaderboardKey string, topN int64) ([]string, error) {
	ctx := context.Background()
	// TODO: use ZRevRange (reverse range = highest score first)
	// or ZRangeArgs with Rev: true
	return client.ZRange(ctx, leaderboardKey, 0, topN-1).Result() // BUG: ascending order
}

// GetScore returns a player's current score.
func GetScore(client *redis.Client, leaderboardKey, player string) (float64, error) {
	ctx := context.Background()
	return client.ZScore(ctx, leaderboardKey, player).Result()
}
