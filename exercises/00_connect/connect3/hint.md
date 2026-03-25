# Hints for connect3

## Hint 1 — The Concept

go-redis maintains a **connection pool** — a set of reusable TCP connections to Dragonfly. When a goroutine wants to run a command, it borrows a connection from the pool, uses it, and returns it.

If `PoolSize=1`, only one goroutine can use Dragonfly at a time. The rest queue up waiting. This turns parallel goroutines into sequential ones — defeating the purpose of concurrency.

`MinIdleConns` pre-creates connections before they're needed. Without it, connections are created on demand — the first burst of requests after a quiet period pays the TCP setup cost.

Think of it like a coffee shop: `PoolSize` = number of coffee machines, `MinIdleConns` = how many machines are pre-heated even when idle.

## Hint 2 — The Specific Issue

Two fields need fixing in `NewPool()`:

1. `PoolSize: 1` — this is a severe bottleneck. Change it to `10` to allow 10 concurrent Dragonfly connections.

2. `MinIdleConns: 0` — this means zero pre-warmed connections. Change it to `5` to keep 5 connections ready at all times.

The test sends 20 concurrent PINGs and checks that `TotalConns >= 2`, which only passes if the pool can grow beyond 1.

## Hint 3 — Near Solution

```go
func NewPool(addr string) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:         addr,
        PoolSize:     10,
        MinIdleConns: 5,
        PoolTimeout:  2 * time.Second,
    })
}
```
