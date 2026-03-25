# Explain It — connect3

## The Challenge

Your team is debugging a performance issue: "Our Redis calls are fast individually, but under load everything slows to a crawl."

You look at the config: `PoolSize: 1`.

Write your diagnosis in `feynman/explanations/connect3.md`.

## Your Explanation Should Cover

1. **What a connection pool is and why it exists** — creating a TCP connection takes time (the three-way handshake). If every Redis command created a new connection, the overhead would dominate. Pools reuse connections.

2. **What happens with PoolSize=1 under concurrent load** — draw it out. 10 goroutines, 1 connection. What's the queue? What's the wait time? How does this relate to Amdahl's Law?

3. **The tradeoff: bigger pool vs resource cost** — each connection uses memory on both the client and Dragonfly sides. There's a sweet spot. What factors would you consider when choosing PoolSize for production?

## Quality Check

Your explanation should make someone think "I've seen this exact bug before!" even if they've never heard of connection pools.

Bonus: explain why `MinIdleConns` matters for "bursty" traffic patterns (quiet → sudden spike → quiet again).
