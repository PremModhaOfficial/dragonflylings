# Explain It — prod5

Your SRE alerts you: *"One Dragonfly shard is at 100% CPU. All traffic is hitting the 'home:featured-products' key. We're getting 50,000 GET requests per second from 20 service instances."*

Write your answer in `feynman/explanations/prod5.md`.

Your explanation must cover:
1. Why adding more Dragonfly nodes doesn't help (which shard gets all the traffic?)
2. How local process-level caching solves the problem (what happens to those 50,000 req/sec?)
3. The consistency trade-off: with 20 instances each caching for 1 second, how stale can the data be?
4. When would you NOT use local caching? (What data must always be fresh?)

**The write problem:** When the featured products list updates, 20 instances each have a stale local cache. How do you invalidate them? (Hint: think Pub/Sub.)

**Memory concern:** sync.Map can grow unbounded if you cache millions of different keys. How do you bound it?

QUALITY CHECK: Can you explain this to a load test engineer who's about to scale-test your service?
