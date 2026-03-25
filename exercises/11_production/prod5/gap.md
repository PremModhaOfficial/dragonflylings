# Gap Notebook — prod5

After completing this exercise, add to `feynman/gap_notebook.md`:

1. How do you detect a hot key before it causes problems? What Dragonfly/Redis metric would alert you?
2. Key splitting: instead of one hot key, spread reads across N copies (e.g., `key:0`, `key:1`, ...). How does this compare to local caching?
3. sync.Map has no max-size or LRU eviction. What happens to memory if your service caches 10 million different keys?
4. Pub/Sub for cache invalidation across instances: what happens if an instance is down when the invalidation message is published?
5. How does CDN caching relate to this pattern? When would you add a CDN in front of this?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. The `sync.Map` local cache has no TTL or LRU eviction. Model the worst case: 10,000 unique product IDs are accessed once and never again. How much memory does the map consume, and at what point does it become a memory leak? What's the minimal addition that bounds it?
2. Key splitting distributes a hot key across N shards (`product:42:shard:0` … `product:42:shard:N-1`). Implement a read function that picks a random shard and falls back to shard 0 on miss — then explain why reads across shards are eventually consistent and when that matters.
3. Local cache and Redis can diverge permanently if an invalidation message is missed (subscriber down, pub/sub message dropped). Design a version-tag reconciliation strategy: how do you detect stale local entries and self-correct them without a full cache flush?
