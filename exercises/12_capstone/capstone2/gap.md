# Gap Notebook — capstone2

After completing this capstone, add to `feynman/gap_notebook.md`:

**Questions this capstone raised for you:**

1. The rate limiter and circuit breaker both protect against overload, but differently. Can they conflict? (What if the circuit is open but rate limit hasn't been hit — should the request be rate-limited or circuit-broken?)

2. Pub/sub invalidation is fire-and-forget. If a subscriber is down during invalidation, it never learns about the update. What's the worst case? How do you recover?

3. This system has no singleflight. Add it: where exactly does singleflight fit, and what's the interaction with the circuit breaker?

4. The rate limiter uses Redis — but if Redis is degraded (circuit open), the rate limiter also can't function. What should you do?

5. Looking at all four patterns combined: which one was hardest to reason about in isolation? Which interaction between patterns surprised you?

**The synthesis question:** You've now implemented cache-aside, rate limiting, distributed locking, session storage, hot key mitigation, circuit breaking, pub/sub, and streams. Which of these would you use in your next project? Which was most surprising?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. This system has no singleflight: 50 concurrent requests for the same uncached key can trigger 50 origin fetches (or 50 rate-limit checks). Add singleflight to the GetOrFetch path — explain how it interacts with the circuit breaker state and whether a singleflight call should count as a circuit failure if the origin (not Redis) fails.
2. The circuit breaker and rate limiter both depend on Redis. If Redis is degraded (circuit open), the rate limiter also cannot function. Design a local fallback rate limiter (in-memory token bucket) that activates automatically when Redis is unavailable — and explain the consistency trade-off across instances.
3. Pub/sub invalidation is fire-and-forget: an instance that is down during invalidation serves stale data indefinitely. Design a version-tag scheme where each cached value embeds a version number validated on read against a Redis counter — so stale entries self-correct on the next access without requiring perfect pub/sub delivery.
