# Hints for connect2

## Hint 1 — The Concept

`redis.NewClient()` does **not** make a network connection. It just creates a Go struct with configuration. The actual TCP connection happens lazily — on the first command you run.

This means: right after `NewClient()`, your code has no idea if Dragonfly is running. To verify connectivity, you must send a command. The conventional choice is `PING`.

The key insight: **verification requires action**. Passive creation tells you nothing about reachability.

## Hint 2 — The Specific Issue

The function creates the client but never calls `Ping()`. So it always returns `(client, nil)` — even when port 19999 has nothing listening.

You need to:
1. Call `client.Ping(ctx)` after creating the client
2. If it errors, close the client and return `(nil, err)`
3. If it succeeds, return `(client, nil)`

Use a context with timeout for the Ping — you don't want to wait forever.

## Hint 3 — Near Solution

```go
func Connect(addr string) (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr:        addr,
        DialTimeout: 500 * time.Millisecond,
    })
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()
    if err := client.Ping(ctx).Err(); err != nil {
        client.Close()
        return nil, err
    }
    return client, nil
}
```
