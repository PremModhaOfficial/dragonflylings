# Hints for hashes3

## Hint 1 — The Concept

`HINCRBY` is like `INCR` but for a field inside a hash. Instead of incrementing a standalone string key, it increments a numeric field within a hash.

This means you can have a hash like:
```
analytics:2024-03-25
  home    → 142
  about   → 31
  contact → 8
```

All in one key. With `INCR`, you'd need a separate key for each page.

`HINCRBY key field increment` — atomically increments the field's integer value by `increment`. If the field doesn't exist, it's initialized to 0 first.

## Hint 2 — The Specific Issue

Two bugs to fix:

1. `client.Incr(ctx, analyticsKey+":"+page)` → `client.HIncrBy(ctx, analyticsKey, page, 1)`
   - This stores the counter as a hash field instead of a separate string key

2. `client.Get(ctx, analyticsKey+":"+page).Int64()` → `client.HGet(ctx, analyticsKey, page).Int64()`
   - Reads the field from the hash

## Hint 3 — Near Solution

```go
func IncrPageView(client *redis.Client, analyticsKey, page string) (int64, error) {
    ctx := context.Background()
    return client.HIncrBy(ctx, analyticsKey, page, 1).Result()
}

func GetPageViews(client *redis.Client, analyticsKey, page string) (int64, error) {
    ctx := context.Background()
    return client.HGet(ctx, analyticsKey, page).Int64()
}
```
