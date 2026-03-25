## Hint 1

**Lock bug:** `SetNX` accepts a TTL as the last argument. Currently it's `0` (no expiry). Pass the `ttl` parameter:

```go
return client.SetNX(ctx, lockKey, token, ttl).Result()
```

That's it for Lock — one character change.

## Hint 2

**Unlock bug:** Two-step GET + DEL has a race condition. Use the `unlockScript` Lua script that's already defined at the top of the file. It atomically checks the token and deletes in one operation:

```go
result, err := client.Eval(ctx, unlockScript, []string{lockKey}, token).Int()
```

If `result == 0`, the token didn't match (lock held by someone else) — return nil (no-op).

## Hint 3

Complete fix:

```go
func Lock(ctx context.Context, client *redis.Client, lockKey, token string, ttl time.Duration) (bool, error) {
    return client.SetNX(ctx, lockKey, token, ttl).Result()
}

func Unlock(ctx context.Context, client *redis.Client, lockKey, token string) error {
    _, err := client.Eval(ctx, unlockScript, []string{lockKey}, token).Int()
    if err == redis.Nil {
        return nil
    }
    return err
}
```

Note: result == 0 means token mismatch — this is not an error, just a no-op.
