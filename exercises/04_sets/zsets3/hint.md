# Hints for zsets3

## Hint 1 — The Concept

Two easily confused operations:

- `ZSCORE key member` → returns the **score value** (e.g., `1500.0`)
- `ZRANK key member` → returns the **position index** in ascending order (e.g., `2`)
- `ZREVRANK key member` → returns position in **descending** order (rank 0 = highest score)

And for updating scores:

- `ZADD key score member` → **replaces** the score with a new value
- `ZINCRBY key increment member` → **adds** increment to the existing score

For leaderboards, you almost always want `ZREVRANK` (position from the top) and `ZINCRBY` (add points earned this round).

## Hint 2 — The Specific Issue

Two bugs to fix:

1. `GetRank` uses `ZScore` → returns score value. Replace with `client.ZRevRank(ctx, leaderboardKey, player).Result()` to get the rank position (0 = best).

2. `IncrScore` uses `ZAdd` with the increment as a new score → replaces the score. Replace with `client.ZIncrBy(ctx, leaderboardKey, points, player).Result()` to add to the existing score.

## Hint 3 — Near Solution

```go
func GetRank(client *redis.Client, leaderboardKey, player string) (int64, error) {
    ctx := context.Background()
    return client.ZRevRank(ctx, leaderboardKey, player).Result()
}

func IncrScore(client *redis.Client, leaderboardKey, player string, points float64) (float64, error) {
    ctx := context.Background()
    return client.ZIncrBy(ctx, leaderboardKey, points, player).Result()
}
```
