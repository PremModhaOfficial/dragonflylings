## Hint 1

Add a `localCache sync.Map` field to `HotKeyCache`. `sync.Map` is safe for concurrent access without a mutex.

In `Get()`, first check the local cache:
```go
if raw, ok := c.localCache.Load(key); ok {
    entry := raw.(localEntry)
    if time.Now().Before(entry.expiresAt) {
        return entry.value, nil
    }
    c.localCache.Delete(key) // expired
}
```

## Hint 2

On Redis hit (local miss), store in local cache:
```go
val, err := c.redis.Get(ctx, key).Result()
if err != nil {
    return "", err
}
c.localCache.Store(key, localEntry{
    value:     val,
    expiresAt: time.Now().Add(c.localTTL),
})
return val, nil
```

## Hint 3

In `Set()`, invalidate the local cache entry:
```go
func (c *HotKeyCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
    c.localCache.Delete(key) // invalidate local cache
    return c.redis.Set(ctx, key, value, ttl).Err()
}
```

And update the struct:
```go
type HotKeyCache struct {
    redis      *redis.Client
    localTTL   time.Duration
    localCache sync.Map
}
```
