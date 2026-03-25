# Hints for hashes1

## Hint 1 — The Concept

A Redis Hash is like a row in a database table: one key, many fields. Instead of storing `user:1:name`, `user:1:email`, `user:1:age` as three separate string keys, you store them as three fields of the hash key `user:1`.

The advantage: all user data lives under one key. Deleting the user deletes one key. Listing all fields uses one command. Updating one field doesn't require reading the others.

Think of it this way: String keys are like filing one sheet per fact. A Hash is like a filing folder — one folder per entity, all facts inside.

## Hint 2 — The Specific Issue

Two commands need changing:

1. `client.Set(ctx, userKey+":"+field, value, 0)` → `client.HSet(ctx, userKey, field, value)`
   - `HSet` takes: key, field, value (not a compound key)

2. `client.Get(ctx, userKey+":"+field)` → `client.HGet(ctx, userKey, field)`
   - `HGet` takes: key, field

With the fix, all fields live under `userKey` as a hash type, not as separate string keys.

## Hint 3 — Near Solution

```go
func SetUserField(client *redis.Client, userKey, field, value string) error {
    ctx := context.Background()
    return client.HSet(ctx, userKey, field, value).Err()
}

func GetUserField(client *redis.Client, userKey, field string) (string, error) {
    ctx := context.Background()
    return client.HGet(ctx, userKey, field).Result()
}
```
