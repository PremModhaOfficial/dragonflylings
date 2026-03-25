# Hints for strings2

## Hint 1 — The Concept

In Redis, every key can have an optional expiration time (TTL — Time To Live). When the TTL reaches zero, Redis automatically deletes the key. This is how session stores, caches, and rate limiters implement "data that expires."

`TTL key` returns:
- A positive duration: the key exists and has this much time left
- `-1`: the key exists but has **no** expiration (lives forever)
- `-2`: the key does not exist

The value `-1` from TTL is a warning sign in production: it means someone forgot to set an expiry.

## Hint 2 — The Specific Issue

`client.Set(ctx, key, value, 0)` — the fourth argument is the expiration duration. `0` means "no expiration." The key lives forever.

To set a key with a TTL, pass the duration as the fourth argument:
- `client.Set(ctx, key, value, 2*time.Second)` — expires in 2 seconds
- `client.Set(ctx, key, value, 24*time.Hour)` — expires in 24 hours

The broken code ignores the `ttl` parameter. Fix: pass `ttl` instead of `0`.

## Hint 3 — Near Solution

```go
func SetSession(client *redis.Client, sessionID, token string, ttl time.Duration) error {
    ctx := context.Background()
    return client.Set(ctx, "session:"+sessionID, token, ttl).Err()
}
```
