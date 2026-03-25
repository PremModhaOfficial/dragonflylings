# Hints for zsets2

## Hint 1 — The Concept

The sliding window rate limiter works by:
1. Each request is stored as a ZSet member with **score = current timestamp**
2. To check the rate: count members with score in `[now - window, now]`
3. Before counting: **remove members older than the window** with `ZREMRANGEBYSCORE`

Without step 3, old requests never get cleaned up. They accumulate and eventually make the counter permanently high, blocking all future requests — even legitimate ones.

`ZREMRANGEBYSCORE key min max` removes all members with scores between min and max.
Use `"-inf"` as min to remove everything older than the window start.

## Hint 2 — The Specific Issue

The commented-out line in `IsAllowed` needs to be uncommented and added before the ZCount:

```go
client.ZRemRangeByScore(ctx, rateLimitKey, "-inf", fmt.Sprintf("%d", windowStart)).Err()
```

This removes all entries with timestamp < windowStart (older than the window). After cleanup, ZCount only sees current-window entries.

## Hint 3 — Near Solution

```go
func IsAllowed(client *redis.Client, rateLimitKey string, limit int, windowSeconds int64) (bool, error) {
    ctx := context.Background()
    now := time.Now().UnixNano()
    windowStart := now - windowSeconds*int64(time.Second)

    // Clean up old entries
    client.ZRemRangeByScore(ctx, rateLimitKey, "-inf", fmt.Sprintf("%d", windowStart))

    count, err := client.ZCount(ctx, rateLimitKey,
        fmt.Sprintf("%d", windowStart), fmt.Sprintf("%d", now),
    ).Result()
    // ... rest of function
```
