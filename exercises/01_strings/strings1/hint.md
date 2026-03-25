# Hints for strings1

## Hint 1 — The Concept

Redis keys are just strings — they can be anything. But a good key naming convention uses colons as separators to create a namespace hierarchy: `user:42:name`, `user:42:email`, `session:abc123`.

The key you write with SET must exactly match the key you read with GET. Redis has no schema, no foreign keys, no type enforcement — the only contract is that you use the same key string on both sides.

## Hint 2 — The Specific Issue

Look at `SetUsername`: it stores the value under `"user:" + userID + ":name"`.

Now look at `GetUsername`: it reads from `"user:name"` — a completely different key. The `userID` part is missing, and there's no suffix.

The fix: make `GetUsername` construct the same key pattern that `SetUsername` uses.

## Hint 3 — Near Solution

```go
func GetUsername(client *redis.Client, userID string) (string, error) {
    ctx := context.Background()
    return client.Get(ctx, "user:"+userID+":name").Result()
}
```
