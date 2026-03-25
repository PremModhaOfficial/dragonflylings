# Hints for strings5

## Hint 1 — The Concept

Every command you send to Redis requires a network round trip: your app → Dragonfly → your app. With a typical 1ms round trip, 100 individual SETs take 100ms. With MSET, those same 100 SETs take 1ms.

`MGET` returns values in the **same order** as the keys you requested — including `nil` for missing keys. This is crucial: it never errors on missing keys. Missing keys simply return `nil` in the result slice.

MSET/MGET are "width" operations: they trade message size for round-trip count. For bulk reads/writes with many keys, this is almost always a win.

## Hint 2 — The Specific Issue

The broken `GetPreferences` treats `redis.Nil` (missing key) as an error and returns early. But `MGET` never errors on missing keys — it returns `nil` in the slice. The function should match that behavior.

Fix `GetPreferences` to use `client.MGet(ctx, keys...).Result()` — it returns `([]interface{}, error)` where missing keys have `nil` entries.

Fix `SetPreferences` to use `client.MSet(ctx, prefs).Err()` — it accepts `map[string]string` directly.

## Hint 3 — Near Solution

```go
func SetPreferences(client *redis.Client, prefs map[string]string) error {
    ctx := context.Background()
    return client.MSet(ctx, prefs).Err()
}

func GetPreferences(client *redis.Client, keys []string) ([]interface{}, error) {
    ctx := context.Background()
    return client.MGet(ctx, keys...).Result()
}
```
