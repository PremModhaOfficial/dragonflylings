# Hints for hashes4

## Hint 1 — The Concept

This exercise makes the contrast explicit: storing user data as multiple string keys vs. one hash. The hash version is almost always better because:

1. **One key to delete** — `DEL user:42` removes everything
2. **One call to read all** — `HGETALL user:42` returns all fields
3. **Less memory overhead** — one key's metadata vs. N keys' metadata
4. **Atomic multi-field reads** — `HMGET` sees a consistent snapshot

The string version (many keys) only wins when fields need individual TTLs or when you access fields across different users in batch with MGET.

## Hint 2 — The Specific Issue

Two functions need fixing:

1. `StoreUser`: Replace the loop with `client.HSet(ctx, userKey, data).Err()`
   - `HSet` accepts `map[string]string` directly

2. `GetUser`: Replace the loop with `client.HGetAll(ctx, userKey).Result()`
   - `HGetAll` returns `map[string]string` — all fields at once

3. `DeleteUser`: Replace the loop with `client.Del(ctx, userKey).Err()`
   - One `DEL` removes the entire hash

## Hint 3 — Near Solution

```go
func StoreUser(client *redis.Client, userKey string, data map[string]string) error {
    ctx := context.Background()
    return client.HSet(ctx, userKey, data).Err()
}

func GetUser(client *redis.Client, userKey string) (map[string]string, error) {
    ctx := context.Background()
    return client.HGetAll(ctx, userKey).Result()
}

func DeleteUser(client *redis.Client, userKey string) error {
    ctx := context.Background()
    return client.Del(ctx, userKey).Err()
}
```
