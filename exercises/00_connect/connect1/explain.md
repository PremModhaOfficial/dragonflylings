# Explain It — connect1

## The Challenge

Your friend just asked:

> "Why does Redis have a PING command? Isn't that just wasting a network round-trip? My database doesn't have a PING."

Write your answer in `feynman/explanations/connect1.md`.

## Your Explanation Should Cover

1. **What PING actually verifies** — it's not just "is the server up?" Think about what has to succeed for PING to return PONG: TCP connection, RESP protocol parsing, server processing, response serialization, and transmission back. Each layer can fail independently.

2. **Why a health check matters in distributed systems** — your app starts, creates a client, then immediately tries to run a critical query. If Dragonfly is down, when do you find out? Compare two scenarios: with and without an explicit PING on startup.

3. **What could go wrong if you skip connection verification** — think about connection pools. If you start 10 goroutines all assuming Dragonfly is up, and it's actually down, what's the user experience?

## Quality Check

Could a 12-year-old follow your explanation?

If you used "RESP protocol", "distributed systems", or "connection pool" — replace each with an analogy before you're done.

Example of good analogizing: "The TCP connection is like picking up the phone. RESP is like the dialing tones. PING/PONG is saying 'hello, can you hear me?' before you start the actual conversation."
