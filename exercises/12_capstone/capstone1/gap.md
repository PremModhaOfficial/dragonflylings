# Gap Notebook — capstone1

After completing this capstone, add to `feynman/gap_notebook.md`:

**Questions this capstone raised for you:**

1. Consumer group rebalancing: if you add a second consumer, how does Redis distribute messages? What happens to in-flight messages from the first consumer?

2. Dead letter queue: some messages fail processing repeatedly. How would you implement a DLQ pattern with Redis Streams?

3. Exactly-once semantics: XACK guarantees at-least-once delivery. How would you achieve exactly-once in this pipeline?

4. Leaderboard sharding: with 10 million players, a single sorted set becomes a hot key. How do you shard a leaderboard?

5. Stream retention: if you use XADD MAXLEN, what happens to old messages that haven't been acknowledged yet?

**The meta-question:** This capstone combined 4 Redis features. Which integration was most complex? Which part were you most uncertain about before completing it?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. This capstone uses a single consumer. Design a parallel consumer group with multiple goroutines, each calling XREADGROUP and processing different messages concurrently — what coordination is needed to prevent two goroutines from claiming the same message, and how does XACK interact with concurrent processing?
2. XPENDING shows messages claimed but not yet acknowledged. Implement a reaper goroutine that picks up pending messages older than 30 seconds and reprocesses them — handle the case where reprocessing is not idempotent (e.g., adding a score twice to the leaderboard).
3. The leaderboard sorted set requires O(log N) per score update. At 1 million players with burst score events, ZADD becomes a bottleneck. Design a two-level leaderboard: a hot tier (top 1000 in sorted set) and a cold tier (bulk storage), with a promotion/demotion mechanism triggered by score thresholds.
