# Gap Notebook — lua4

After completing this exercise, add to `feynman/gap_notebook.md`:

1. What happens if you use hashtags but your tag value is always the same (e.g., `{global}`)? How many shards actually serve your traffic?
2. Does Dragonfly run in single-shard mode by default locally? How do you check the current shard count?
3. Can you use hashtags with regular (non-Lua) multi-key commands like `MGET`? What's the rule?
4. Is the hashtag spec the same between Redis Cluster and Dragonfly? Are there any differences?
5. What would a monitoring alert for cross-slot errors look like? What metric would you watch?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. You're designing a multi-tenant system where each tenant's data must be co-located for atomic Lua operations. Design a key naming scheme using hashtags that supports tenant isolation AND cross-tenant admin queries.
2. Hashtag routing fixes data locality at key creation time. What happens when Dragonfly adds a thread (shard rebalancing)? Can keys be migrated to new shards, and what's the downtime model?
3. If all hot keys share `{global}` as a hashtag, you've created a single-thread bottleneck despite Dragonfly's multi-threading. How would you detect this in production metrics, and what's the refactoring path without downtime?
