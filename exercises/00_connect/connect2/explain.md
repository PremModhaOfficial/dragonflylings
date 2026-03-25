# Explain It — connect2

## The Challenge

Your colleague wrote this code:

```go
client := redis.NewClient(&redis.Options{Addr: "localhost:6380"})
// ... 50 lines later ...
result, err := client.Get(ctx, "user:42").Result()
```

They say: "I checked — the client is not nil. Why is my GET failing?"

Write your explanation in `feynman/explanations/connect2.md`.

## Your Explanation Should Cover

1. **Why "client is not nil" doesn't mean "connected"** — explain the lazy connection model in plain terms. What does `NewClient()` actually create?

2. **When the error actually surfaces** — at what exact line does the code first touch the network? What would the error message look like?

3. **What the fix looks like and why** — not just the code change, but *why* this pattern (verify on startup) is a best practice in distributed systems.

## Quality Check

Imagine explaining this to someone who has only used SQLite before. In SQLite, opening the database file either works or errors immediately. Redis is different — can you explain *why* without using jargon?

Your explanation passes if someone reading it says: "Oh, that's why my Redis app sometimes fails in production but not locally!"
