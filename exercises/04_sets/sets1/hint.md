# Hints for sets1

## Hint 1 — The Concept

A Redis **Set** is an unordered collection of **unique** strings. SADD ignores duplicates — adding "go" three times results in one "go" in the set. This is the fundamental property that makes sets useful for tag systems, unique visitor tracking, and membership testing.

`SISMEMBER key value` checks membership in O(1) — regardless of set size. This is one of Redis's most powerful operations: "Does X belong to this group?" answered instantly.

With a List (LPUSH), duplicates are stored and SISMEMBER doesn't work (wrong data type).

## Hint 2 — The Specific Issue

Two bugs to fix:

1. `client.LPush(ctx, itemKey+":tags", tag)` → `client.SAdd(ctx, itemKey+":tags", tag)`
2. `client.LRange(ctx, itemKey+":tags", 0, -1)` → `client.SMembers(ctx, itemKey+":tags")`

Note: `SMembers` returns `([]string, error)` — same as `LRange`. The return type is compatible.

`HasTag` already uses `SIsMember` — it was correct all along, but it fails with LPUSH because the key type is "list" not "set."

## Hint 3 — Near Solution

```go
func AddTag(client *redis.Client, itemKey, tag string) error {
    ctx := context.Background()
    return client.SAdd(ctx, itemKey+":tags", tag).Err()
}

func GetTags(client *redis.Client, itemKey string) ([]string, error) {
    ctx := context.Background()
    return client.SMembers(ctx, itemKey+":tags").Result()
}
```
