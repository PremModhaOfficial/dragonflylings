## Hint 1

`FindPendingMessages` uses `client.XPending(ctx, stream, group)`. This returns aggregate stats (total count, min/max IDs, per-consumer counts) but NOT individual message IDs. To get IDs, use `client.XPendingExt(ctx, &redis.XPendingExtArgs{...})`.

## Hint 2

`XPendingExtArgs` requires `Stream`, `Group`, `Start` (`"-"`), `Stop` (`"+"`), and `Count`. The result is `[]redis.XPendingExt` where each item has an `ID` field.

## Hint 3

```go
// Fix FindPendingMessages:
entries, err := client.XPendingExt(ctx, &redis.XPendingExtArgs{
    Stream: stream,
    Group:  group,
    Start:  "-",
    Stop:   "+",
    Count:  100,
}).Result()
// extract IDs:
for _, e := range entries { ids = append(ids, e.ID) }

// Fix ClaimMessages:
MinIdle: 0,  // claim immediately, regardless of idle time
```
