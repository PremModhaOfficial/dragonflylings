# Hints for hashes2

## Hint 1 — The Concept

`HMGET` is to hashes what `MGET` is to strings: batch retrieval in one round trip. You specify the key and a list of fields, and you get back a list of values in the same order.

Crucially, `HMGET` **never errors on missing fields** — it returns `nil` for fields that don't exist. This is different from `HGET`, which returns `redis.Nil` (an error). This distinction matters: your code that calls `HMGET` shouldn't treat missing fields as failures.

## Hint 2 — The Specific Issue

The broken `GetUserFields` loops with `HGet`, which returns `redis.Nil` for missing fields and the loop treats that as an error.

Replace the loop with `client.HMGet(ctx, userKey, fields...).Result()`.

`HMGet` returns `([]interface{}, error)`. Missing fields appear as `nil` in the slice. The overall call only errors if something fundamental goes wrong (network, wrong type), not if fields are missing.

## Hint 3 — Near Solution

```go
func GetUserFields(client *redis.Client, userKey string, fields []string) ([]interface{}, error) {
    ctx := context.Background()
    return client.HMGet(ctx, userKey, fields...).Result()
}
```
