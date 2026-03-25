## Hint 1

A pipeline in go-redis is like a to-do list: writing items to the list doesn't do anything until you submit it. Look at what's missing after the `for` loop in `SetManyPipelined`.

## Hint 2

The `pipe.Pipeline()` call creates a buffer. Commands like `pipe.Set(...)` add to the buffer. To actually send all buffered commands to Dragonfly in one round trip, you need to call a method on `pipe` that flushes the buffer and returns results.

## Hint 3

```go
func SetManyPipelined(client *redis.Client, ctx context.Context, prefix string, n int) error {
    pipe := client.Pipeline()
    for i := 0; i < n; i++ {
        key := fmt.Sprintf("%s:%d", prefix, i)
        pipe.Set(ctx, key, i, 0)
    }
    _, err := pipe.Exec(ctx)  // <-- this is what was missing
    return err
}
```
