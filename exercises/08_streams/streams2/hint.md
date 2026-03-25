## Hint 1

`CountEvents` calls `client.XLen(ctx, stream+":count")`. XLen takes the stream name directly -- no suffix needed. Remove the `+":count"` part.

## Hint 2

`QueryRange` calls `client.XRange(ctx, stream, stop, start)` -- the third and fourth arguments are reversed. XRANGE signature is `XRange(ctx, key, start, stop)`.

## Hint 3

```go
func CountEvents(client *redis.Client, ctx context.Context, stream string) (int64, error) {
    return client.XLen(ctx, stream).Result()  // just stream, no suffix
}

func QueryRange(...) ([]redis.XMessage, error) {
    return client.XRange(ctx, stream, start, stop).Result()  // start before stop
}
```
