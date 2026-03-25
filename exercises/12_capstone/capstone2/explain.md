# Explain It — capstone2

You've built a production-grade API cache. Now design it for the real world.

Write your answer in `feynman/explanations/capstone2.md`.

**Part 1: Component Interaction**

Draw all four components (cache-aside, rate limiter, circuit breaker, pub/sub) and show:
- The order of checks on each request
- What each component returns on failure
- Which failures should trip the circuit breaker and which shouldn't

**Part 2: The Failure Scenarios**

For each failure scenario, trace through your system:
1. Dragonfly is completely unreachable (network partition)
2. Dragonfly is slow (100ms+ responses instead of <1ms)
3. Cache is poisoned (wrong value was cached)
4. A burst of 10,000 requests from one user ID

**Part 3: Operational Concerns**

- How would you monitor the rate limiter? (What's the alert: "user X hit rate limit" or "X% of requests are rate limited"?)
- Circuit breaker per-instance vs. shared state: your 10 instances each have separate circuit state. What's the impact?
- Pub/sub invalidation: one message reaches all 10 instances. What if an instance is restarting when the message is published?

**Part 4: Improvements**

Looking at your implementation, what would you add for production?
1. One thing for observability
2. One thing for correctness
3. One thing for performance

QUALITY CHECK: Could you present this system design in 20 minutes to a senior engineer who has never seen this codebase?
