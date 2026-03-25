## Hint 1

`sync.WaitGroup` tracks in-flight goroutines. You call `wg.Add(1)` before launching each goroutine, `wg.Done()` when it finishes, and `wg.Wait()` to block until all are done.

The current code calls `wg.Add(1)` but the goroutine never calls `wg.Done()`. So `wg.Wait()` would block forever (if it were called). And it isn't called at all — so the function returns immediately.

Two fixes needed: add `wg.Done()` inside the goroutine, and add `wg.Wait()` before the return.

## Hint 2

`wg.Done()` must go inside the goroutine, called via `defer` so it always runs even if the Set fails:

```go
go func(k string) {
    defer wg.Done()  // ← add this
    if err := client.Set(ctx, k, value, 0).Err(); err == nil {
        atomic.AddInt64(&count, 1)
    }
}(key)
```

## Hint 3 (Pipelining bug)

Redis pipelining works by buffering all commands in the client and sending them in one batch. The key is that you create the pipeline **once**, queue all commands, then call `Exec` once. Creating a new pipeline inside the loop cancels that benefit — each pipeline has exactly one command and is flushed immediately.

Move `client.Pipeline()` and `pipe.Exec(ctx)` outside the loop:

```go
pipe := client.Pipeline()
for _, key := range keys {
    pipe.Set(ctx, key, value, 0)
}
cmds, err := pipe.Exec(ctx)
// count successes from cmds
```

## Hint 4 (WaitGroup complete fix — two lines added):

```go
func SetConcurrent(ctx context.Context, client *redis.Client, keys []string, value string) int64 {
    var wg sync.WaitGroup
    var count int64

    for _, key := range keys {
        wg.Add(1)
        go func(k string) {
            defer wg.Done()                                    // ← fix 1
            if err := client.Set(ctx, k, value, 0).Err(); err == nil {
                atomic.AddInt64(&count, 1)
            }
        }(key)
    }
    wg.Wait()   // ← fix 2
    return count
}
```
