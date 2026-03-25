# Explain It — prod6

Your service is experiencing a Dragonfly outage. Without a circuit breaker, your service is making 10,000 Redis calls per second, each timing out after 5 seconds.

Write your answer in `feynman/explanations/prod6.md`.

Your explanation must cover:
1. What happens to goroutines during a Redis timeout (where are they, what are they doing)?
2. How does this cascade: why do Redis timeouts eventually crash the entire service?
3. How the circuit breaker stops the cascade: what does the service serve during the open state?
4. The half-open probe: why not just reopen immediately when cooldown expires?

**The fallback question:** When the circuit is open and we return `ErrCircuitOpen`, what should callers do? The exercise doesn't implement a fallback — design one for a key-value cache scenario.

**Monitoring question:** What metrics would you expose from a circuit breaker? What would a Grafana dashboard for circuit breaker health look like?

QUALITY CHECK: Draw the three-state machine as a diagram. Add it to your explanation.
