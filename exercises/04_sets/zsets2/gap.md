# Gap Notebook — zsets2

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **Is IsAllowed atomic? What's the race condition?**
   Format: `- [ ] [your question] -- [why it matters]`
   ZRemRangeByScore + ZCount + ZAdd are 3 separate commands. Two concurrent requests could both pass the count check. How do you fix this?

2. **What happens to the ZSet when requests stop?**
   Old entries accumulate until the next request triggers cleanup. Is there a background cleanup option?

3. **What's the time complexity of ZREMRANGEBYSCORE?**
   If a user makes 10,000 requests, then we clean up, how long does the cleanup take?

## Push Further

- Implement the rate limiter using Lua script to make it atomic
- Compare the ZSet rate limiter to the token bucket algorithm — which is more forgiving for burst traffic?
- How would you implement "per user per endpoint" rate limiting? What's the key naming scheme?
