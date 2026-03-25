## Hint 1

`singleflight.Group` deduplicates concurrent calls with the same key. If 50 goroutines call `group.Do("same-key", fn)` at the same time, only ONE executes `fn` — the other 49 wait and get the same result.

Import: `"golang.org/x/sync/singleflight"`

Add a field to `Cache`:
```go
type Cache struct {
    redis   *redis.Client
    ttl     time.Duration
    sfGroup singleflight.Group
}
```

## Hint 2

In `Get()`, wrap the fetch-and-cache logic with `c.sfGroup.Do(key, ...)`:

```go
result, err, _ := c.sfGroup.Do(key, func() (interface{}, error) {
    val, err := fetch()
    if err != nil {
        return "", err
    }
    c.redis.Set(ctx, key, val, c.ttl)
    return val, nil
})
```

The third return value is `shared bool` — true if the result was shared across multiple callers.

## Hint 3

Complete `Get()` with singleflight:

```go
func (c *Cache) Get(ctx context.Context, key string, fetch func() (string, error)) (string, error) {
    val, err := c.redis.Get(ctx, key).Result()
    if err == nil {
        return val, nil
    }
    if err != redis.Nil {
        return "", err
    }

    result, err, _ := c.sfGroup.Do(key, func() (interface{}, error) {
        v, fetchErr := fetch()
        if fetchErr != nil {
            return "", fetchErr
        }
        c.redis.Set(ctx, key, v, c.ttl)
        return v, nil
    })
    if err != nil {
        return "", err
    }
    return result.(string), nil
}
```
