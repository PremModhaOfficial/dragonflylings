## Hint 1

**SetSession TTL bug:** After `HSet`, add one line:
```go
client.Expire(ctx, key, sessionTTL)
```

**GetSession TTL refresh bug:** After `HGetAll`, add:
```go
if len(fields) > 0 {
    client.Expire(ctx, key, sessionTTL)
}
```

Two lines — one in each function.

## Hint 2

**GetSession graceful degradation bug:** When Redis is down, `HGetAll` returns an error. Instead of returning it, check if it looks like a connectivity error and return an empty session:

```go
fields, err := client.HGetAll(ctx, key).Result()
if err != nil {
    if isRedisDown(err) {
        return Session{ID: id, Data: map[string]string{}}, nil
    }
    return Session{}, err
}
```

The `isRedisDown()` helper is already implemented in main.go.

## Hint 3

Complete fixes for both functions:

```go
func SetSession(ctx context.Context, client *redis.Client, s Session) error {
    key := sessionKey(s.ID)
    if err := client.HSet(ctx, key, toFields(s.Data)).Err(); err != nil {
        return err
    }
    return client.Expire(ctx, key, sessionTTL).Err()
}

func GetSession(ctx context.Context, client *redis.Client, id string) (Session, error) {
    key := sessionKey(id)
    fields, err := client.HGetAll(ctx, key).Result()
    if err != nil {
        if isRedisDown(err) {
            return Session{ID: id, Data: map[string]string{}}, nil
        }
        return Session{}, err
    }
    if len(fields) > 0 {
        client.Expire(ctx, key, sessionTTL) // refresh TTL on access
    }
    return Session{ID: id, Data: fields}, nil
}
```
