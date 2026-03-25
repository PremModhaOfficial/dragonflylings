# Explain It — lists3

## The Challenge

Your teammate implemented a worker like this:

```go
for {
    job, err := client.LPop(ctx, "jobs").Result()
    if err == redis.Nil {
        time.Sleep(100 * time.Millisecond) // poll every 100ms
        continue
    }
    process(job)
}
```

They say it works. Write your code review in `feynman/explanations/lists3.md`.

## Your Explanation Should Cover

1. **The polling problem** — 100ms sleep means up to 100ms of job latency. In a busy system, you're also making 10 Redis calls per second per worker just to check if there's work. With 100 workers, that's 1,000 pointless Redis calls per second.

2. **How BLPOP solves it** — the OS puts the connection to sleep. When a new job arrives, Dragonfly wakes exactly the right worker. Zero polling, near-zero latency.

3. **The timeout parameter** — why does BLPOP have a timeout? (Hint: connections can silently die. A worker blocked forever might be a zombie. Reconnecting periodically is good practice.)

## Quality Check

Your review should propose the BLPOP version with working code, and explain why you'd pick 30 seconds as a reasonable timeout for a production worker.
