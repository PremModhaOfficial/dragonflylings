# Gap Notebook — sets2

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **What's the time complexity of SINTER for large sets?**
   Format: `- [ ] [your question] -- [why it matters]`
   If alice follows 10,000 people and bob follows 10,000 people, how long does SINTER take?

2. **Can SINTER work on more than 2 sets?**
   `SINTER key1 key2 key3` — what does it return? When would you need 3-way intersection?

3. **What's SINTERCARD?**
   Added in Redis 7.0. Returns just the COUNT of the intersection without returning all members. When is this more efficient?

## Push Further

- What's `SMOVE`? How does it atomically move a member from one set to another?
- Design "mutual friends" as a feature: given two users, show their common friends with profile data. What Redis calls do you need?
- How would you find the most common tags across 1000 posts using set operations?
