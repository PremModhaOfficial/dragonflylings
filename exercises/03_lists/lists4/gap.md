# Gap Notebook — lists4

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **Is LTRIM + LPUSH atomic in Dragonfly?**
   Format: `- [ ] [your question] -- [why it matters]`
   If not, what's the risk? Can you make it atomic with pipelining?

2. **What does LTRIM do when the range covers the whole list?**
   `LTRIM mylist 0 -1` — does it do anything? What about `LTRIM mylist 0 0` on a 5-element list?

3. **Can you trim from the right instead of the left?**
   `LTRIM mylist -N -1` keeps the last N items. When would you use this instead of `0 N-1`?

## Push Further

- What's the time complexity of LTRIM? Is it O(1) for trimming from the front?
- Design a "sliding window" rate limiter using a list + LTRIM. Compare it to the sorted set approach in module 04.
- What's the memory difference between keeping 5 items in a list vs 5 items in a sorted set?
