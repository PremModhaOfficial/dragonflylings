# Gap Notebook — sets1

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **What does SADD return when a member already exists?**
   Format: `- [ ] [your question] -- [why it matters]`
   SADD returns the number of elements actually added (not already present). When is this count useful?

2. **Can you do SADD with multiple members at once?**
   `SADD key a b c` — does it work? What's the return value?

3. **Is SMEMBERS safe to call on a set with 1 million members?**
   SMEMBERS returns everything. What command should you use instead for large sets?

## Push Further

- What's `SRANDMEMBER`? How would you use it to randomly recommend a tag?
- What's `SPOP`? How is it different from SRANDMEMBER?
- Design "users who are online" using a set. What operations do you need for join, leave, count, and "is X online?"
