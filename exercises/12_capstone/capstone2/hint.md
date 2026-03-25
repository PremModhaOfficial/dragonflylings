# No Hints

This is the final capstone. You've built every component of this system in previous modules:

- Rate limiting: see prod3 (sliding window)
- Circuit breaker: see prod6 (state machine)
- Cache-aside: see prod1
- Pub/Sub: see module 07

The bugs are described in the `main.go` comments. The tests tell you exactly what must pass.

**Estimated time: 45-60 minutes**

If you finish early: add singleflight to `Fetch()` for thundering herd protection (see prod1).
