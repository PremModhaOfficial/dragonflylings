# Gap Notebook — zsets3

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **What does ZREVRANK return for a member that doesn't exist?**
   Format: `- [ ] [your question] -- [why it matters]`
   Is it redis.Nil? -1? An error? How do you handle "player not on leaderboard"?

2. **Can ZINCRBY create a new member?**
   If "alice" isn't in the ZSet and you ZINCRBY 100, what happens?

3. **What's ZRANGEBYSCORE with WITHSCORES?**
   How do you get both the member name AND their score in one call?

## Push Further

- What's `ZMSCORE`? When would you need scores for multiple members at once?
- How would you implement a "weekly reset" for the leaderboard without deleting all data? (Hint: multiple leaderboard keys)
- Look up `ZRANGEBYLEX` — you can use ZSets as a sorted string index. What's a use case?
