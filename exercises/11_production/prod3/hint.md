## Hint 1

The sliding window algorithm uses a sorted set where each member is a request timestamp (score = nanosecond timestamp). Three steps each request:

1. Remove old entries: `ZREMRANGEBYSCORE key 0 <windowStart>`
2. Add this request: `ZADD key <now> <now>` (score = member = timestamp)
3. Count entries: `ZCARD key`

If count ≤ limit: allow. Also set `EXPIRE` so the key eventually cleans up.

## Hint 2

`windowStart` is `time.Now().Add(-window).UnixNano()`. Use a pipeline for all 4 operations:

```go
now := time.Now().UnixNano()
windowStart := time.Now().Add(-window).UnixNano()
member := strconv.FormatInt(now, 10)

pipe := client.Pipeline()
pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))
pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: member})
cardCmd := pipe.ZCard(ctx, key)
pipe.Expire(ctx, key, window)
```

## Hint 3

Complete sliding window implementation:

```go
func Allow(ctx context.Context, client *redis.Client, key string, limit int, window time.Duration) (bool, error) {
    now := time.Now()
    windowStart := now.Add(-window).UnixNano()
    member := strconv.FormatInt(now.UnixNano(), 10)

    pipe := client.Pipeline()
    pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))
    pipe.ZAdd(ctx, key, redis.Z{Score: float64(now.UnixNano()), Member: member})
    cardCmd := pipe.ZCard(ctx, key)
    pipe.Expire(ctx, key, window)

    _, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }
    return cardCmd.Val() <= int64(limit), nil
}
```
