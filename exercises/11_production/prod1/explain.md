# Explain It — prod1

Your on-call engineer pages you at 2am: *"Cache key 'featured-products' just expired and the database is getting hammered with 800 queries per second from a single service. The DB is falling over."*

Write your answer in `feynman/explanations/prod1.md`.

Your explanation must cover:
1. Exactly what happened: why did the DB get 800 QPS instead of 1?
2. How `singleflight` prevents this at the process level
3. The remaining gap: what if you have 10 instances of this service? Does singleflight help?
4. What additional approaches work across multiple service instances?

**The deeper trade-off:** singleflight means all 50 waiting goroutines get the SAME value even if the first one was slightly stale. Is this a problem? When?

**Production scenario:** Your `fetch()` function sometimes returns errors (DB flaky). With singleflight, all 50 waiters get that error. What should you do about error propagation in singleflight?

QUALITY CHECK: Could you explain the thundering herd to a product manager who doesn't code?
