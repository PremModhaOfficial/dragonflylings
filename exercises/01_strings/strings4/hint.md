# Hints for strings4

## Hint 1 — The Concept

`INCR` is an atomic increment operation. It reads the current value, adds 1, and writes it back — all as a single indivisible operation that no other command can interleave with.

Without atomicity, a "read-modify-write" sequence (GET → parse → SET) has a race condition: two goroutines can both read the same value, both add 1, and both write back — resulting in only one increment recorded, not two.

This is called a "lost update." With INCR, it's impossible — Redis guarantees that every INCR is counted exactly once.

## Hint 2 — The Specific Issue

The broken `IncrementCounter` does three separate operations:
1. `client.Get()` — reads current value
2. Parse and add 1 — in Go
3. `client.Set()` — writes new value

Between steps 1 and 3, another goroutine can run the same sequence. The test proves this: 50 concurrent goroutines only increment the counter some of the time (some increments are "lost").

Replace the entire block with `client.Incr(ctx, counterKey).Result()` — one atomic operation.

## Hint 3 — Near Solution

```go
func IncrementCounter(client *redis.Client, counterKey string) (int64, error) {
    ctx := context.Background()
    return client.Incr(ctx, counterKey).Result()
}
```

INCR initializes to 0 if the key doesn't exist, then returns 1. No need to handle the missing key case separately.
