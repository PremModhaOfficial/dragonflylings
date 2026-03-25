# Hints for lists4

## Hint 1 — The Concept

`LTRIM key start stop` modifies the list **in place**, keeping only elements between `start` and `stop` (inclusive) and discarding the rest.

To keep the last N items (most recent, assuming LPush = most recent at index 0):
```
LTRIM key 0 N-1
```

This is the "bounded list" pattern: after every push, trim. The list never grows beyond `maxItems` elements. It's used for activity feeds, recent events, sliding windows.

## Hint 2 — The Specific Issue

The missing line is `client.LTrim(ctx, listKey, 0, maxItems-1).Err()` after the `LPush`.

Without LTRIM: after 10 pushes with maxItems=5, the list has 10 items. The test checks `LLEN == 5` and fails.

With LTRIM: after each push, trim to `0..maxItems-1`. The list always has at most `maxItems` items.

## Hint 3 — Near Solution

```go
func AddEvent(client *redis.Client, listKey, event string, maxItems int64) error {
    ctx := context.Background()
    if err := client.LPush(ctx, listKey, event).Err(); err != nil {
        return err
    }
    return client.LTrim(ctx, listKey, 0, maxItems-1).Err()
}
```
