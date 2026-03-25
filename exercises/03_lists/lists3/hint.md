# Hints for lists3

## Hint 1 — The Concept

`LPOP` is non-blocking: if the list is empty, it returns immediately with "nil" (no job, go away).

`BLPOP` is blocking: if the list is empty, it **waits** until either a job arrives or the timeout expires. This is perfect for workers — the worker process sits idle (no CPU spinning) until there's actual work to do.

The key difference in a worker loop:
- With LPOP: worker must sleep + retry in a loop (polling), which wastes CPU or adds latency
- With BLPOP: worker sleeps at the OS level until Dragonfly wakes it up (efficient, zero latency)

## Hint 2 — The Specific Issue

Replace `client.LPop(ctx, queueKey).Result()` with `client.BLPop(ctx, timeout, queueKey).Result()`.

Important: `BLPop` returns `([]string, error)`:
- `result[0]` = the key that had the item (useful when watching multiple queues)
- `result[1]` = the actual value

Extract the job from `result[1]`.

## Hint 3 — Near Solution

```go
func WaitForJob(client *redis.Client, queueKey string, timeout time.Duration) (string, error) {
    ctx := context.Background()
    result, err := client.BLPop(ctx, timeout, queueKey).Result()
    if err != nil {
        return "", err
    }
    return result[1], nil
}
```
