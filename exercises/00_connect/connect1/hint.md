# Hints for connect1

## Hint 1 — The Concept

Dragonfly in this project is configured to run on port **6380**, not the Redis default of 6379. This is intentional — it lets you run both Redis and Dragonfly side by side without a port conflict.

In production, you'd typically use environment variables or config files to set the address. For learning, we hardcode it.

Think about: what does the port number *mean*? It's just a number the OS uses to route traffic to the right process. Redis chose 6379 (there's a fun story about it). We chose 6380.

## Hint 2 — The Specific Issue

There are **two bugs**:

1. The `Addr` field in `Connect()` uses port `6379`. Change it to `6380`.

2. In `Ping()`, `client.Echo(ctx, "hello")` sends the Redis `ECHO` command — it echoes the string "hello" back. That's not what we want. We want the `PING` command, which returns `"PONG"`.

Look at the go-redis documentation or autocomplete for the right method name.

## Hint 3 — Near Solution

```go
func Connect() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr: "localhost:6380",
    })
}

func Ping(client *redis.Client) (string, error) {
    ctx := context.Background()
    return client.Ping(ctx).Result()
}
```
