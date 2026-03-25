# Gap Notebook — hashes2

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **Is there an HMSET command? Is it still recommended?**
   Format: `- [ ] [your question] -- [why it matters]`
   HMSET was deprecated in Redis 4.0. What replaced it? What does HSET do now that it couldn't before?

2. **What's the maximum number of fields you can HMGET at once?**
   Is there a limit? What happens to performance as you request more and more fields?

3. **Can you HMGET across different hash keys in one call?**
   MGET lets you read string keys across different keys. Is there a hash equivalent for multi-key batch reads?

## Push Further

- What's `HRANDFIELD`? When would randomly sampling hash fields be useful?
- Look up `HGETALL` — when would you use it over `HMGET`? What's the risk with large hashes?
- Is there a way to atomically read some fields AND update others in one operation?
