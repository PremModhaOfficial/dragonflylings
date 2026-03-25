# Gap Notebook — zsets1

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **Can two members have the same score?**
   Format: `- [ ] [your question] -- [why it matters]`
   What happens when scores are equal? What's the tiebreaker? Is it insertion order?

2. **What does ZRANGEBYSCORE return?**
   How is it different from ZRANGE? When would you use score-based range vs index-based range?

3. **What's ZRANGEBYLEX?**
   When scores are all 0, members are sorted lexicographically. How does ZRANGEBYLEX let you do prefix searches?

## Push Further

- What's `ZREVRANGEBYSCORE` and when would you use it?
- How would you implement a "paginated leaderboard" with ZREVRANGE? What's the page formula?
- What's `ZPOPMAX` and `ZPOPMIN`? How would you use them to implement a priority queue?
