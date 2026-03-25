## Hint 1

`GetUsedMemory` should use `client.Info(ctx, "memory")` not `client.DBSize`. The INFO command returns a multi-line string. Look for the line `used_memory:<bytes>`.

The helper `parseInfoField(info, "used_memory")` is already written for you. Also `parseMemoryBytes(info)` combines both steps.

## Hint 2

For `WaitForSnapshot`: after `BGSAVE` starts, poll `INFO persistence` and check `rdb_bgsave_in_progress`. When it's `"0"`, the snapshot is done.

The helper `pollUntilDone(ctx, client, interval)` is already implemented. `WaitForSnapshot` just needs to call it.

## Hint 3

Complete fix:

```go
func GetUsedMemory(ctx context.Context, client *redis.Client) (int64, error) {
    info, err := client.Info(ctx, "memory").Result()
    if err != nil {
        return 0, err
    }
    return parseMemoryBytes(info)
}

func WaitForSnapshot(ctx context.Context, client *redis.Client) error {
    return pollUntilDone(ctx, client, 100*time.Millisecond)
}
```
