# Explain It — zsets2

## The Challenge

Your teammate says: "I implemented a rate limiter with a simple counter and TTL. Why do you need a sorted set?"

```go
count, _ := client.Incr(ctx, "rate:user:42").Result()
client.Expire(ctx, "rate:user:42", time.Minute)
if count > 100 { return rateLimited }
```

Write your comparison in `feynman/explanations/zsets2.md`.

## Your Explanation Should Cover

1. **The fixed window problem** — with INCR + TTL, the window resets every minute at a fixed boundary. A user can make 100 requests at 11:59 and 100 more at 12:00 — 200 requests in 2 minutes, but never "rate limited." The sliding window doesn't have this problem.

2. **How ZSet timestamps enable sliding windows** — each request has a precise timestamp. At any point, counting members in `[now-60s, now]` gives the exact request count in the last 60 seconds, regardless of minute boundaries.

3. **The memory tradeoff** — ZSet stores one entry per request. The INCR approach stores one integer total. For 100 req/min × 1M users, how much memory does each approach use?

## Quality Check

Draw a timeline showing where the fixed window rate limiter fails but the sliding window succeeds. Make it concrete with specific timestamps.
