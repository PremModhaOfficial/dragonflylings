# Gap Notebook — prod6

After completing this exercise, add to `feynman/gap_notebook.md`:

1. This circuit breaker is per-process. If you have 20 service instances, each has its own circuit state. Is this a problem?
2. What's the difference between a circuit breaker and a retry-with-backoff? When would you use each?
3. Half-open allows ONE probe call. What if you want to ramp up traffic gradually after recovery instead of flipping immediately to closed?
4. The circuit breaker here counts ALL errors as failures. Should `redis.Nil` count? What about context.Canceled?
5. Libraries like `sony/gobreaker` and `afex/hystrix-go` implement circuit breakers. When would you use a library vs this implementation?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. With 20 service instances each holding independent circuit state, some instances may open (protecting themselves) while others keep calling a degraded Redis. Design a shared circuit state stored in Redis itself — and handle the meta-problem: the shared state call must not itself be circuit-broken in the same way.
2. The half-open state allows exactly one probe. Implement graduated recovery instead: after a successful probe, allow 10% of traffic through, then 50%, then 100%, backing off to open if any failure occurs during the ramp. What data structure tracks the current ramp percentage?
3. A circuit breaker detects faults; a bulkhead limits concurrent load. Explain the difference, then design a combined pattern: circuit breaker for fault detection + a semaphore (`chan struct{}`) that limits concurrent Redis calls even when the circuit is closed.
