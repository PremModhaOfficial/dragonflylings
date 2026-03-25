# Gap Notebook — prod3

After completing this exercise, add to `feynman/gap_notebook.md`:

1. The sorted set sliding window stores one entry per request. For very high-rate buckets (10,000 req/min), is there a more memory-efficient approach?
2. What happens if two requests arrive at the exact same nanosecond? They'd have the same member value — ZADD would deduplicate them. How do you prevent this?
3. Is the pipeline approach truly atomic? What happens if Dragonfly crashes between ZREMRANGEBYSCORE and ZADD?
4. Can you make the rate limiter work across multiple keys (e.g., per-user AND per-IP) atomically? What's the challenge?
5. What's the difference between rate limiting and throttling? When would you use each?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. The sorted-set window stores one member per request. At 10,000 req/sec with a 60-second window, one user consumes 600,000 members. Design a fixed-bucket approximation (e.g., one member per second bucket) that bounds memory to O(window_size) regardless of request rate — and quantify the accuracy trade-off.
2. Implement a token bucket rate limiter in Redis using a single Lua script: the bucket refills at R tokens/sec up to capacity C, and each request atomically checks and consumes one token. Compare the memory and accuracy characteristics to the sorted-set approach.
3. Your rate limiter is per-request-key. A multi-tenant SaaS needs per-tenant limits configurable at runtime (tenant A: 100 req/min, tenant B: 1000 req/min). How do you store and look up per-tenant config with sub-millisecond latency without a separate DB call on every request?
