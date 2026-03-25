# Hints for zsets1

## Hint 1 — The Concept

A Sorted Set (ZSet) is like a Set, but every member has a floating-point **score**. Members are always stored in ascending score order. This makes sorted sets perfect for leaderboards, priority queues, and anything that needs ranking.

`ZADD key score member` — adds or updates a member with a score
`ZRANGE key start stop` — returns members in ASCENDING score order (lowest first)
`ZREVRANGE key start stop` — returns members in DESCENDING score order (highest first)
`ZSCORE key member` — returns a member's score

A leaderboard shows the HIGHEST scores first — that's `ZREVRANGE`, not `ZRANGE`.

## Hint 2 — The Specific Issue

`GetLeaderboard` uses `ZRange` which returns members sorted ascending (lowest score first). For a leaderboard, you want descending (highest score first).

Replace `ZRange` with `ZRevRange`:

```go
client.ZRevRange(ctx, leaderboardKey, 0, topN-1).Result()
```

`ZRevRange 0 N-1` returns the top N members by score, highest first.

## Hint 3 — Near Solution

```go
func GetLeaderboard(client *redis.Client, leaderboardKey string, topN int64) ([]string, error) {
    ctx := context.Background()
    return client.ZRevRange(ctx, leaderboardKey, 0, topN-1).Result()
}
```
