# Gap Notebook — strings4

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **Is INCR truly atomic in Dragonfly, which is multi-threaded?**
   Format: `- [ ] [your question] -- [why it matters]`
   Redis is single-threaded so atomicity is easy. Dragonfly uses multiple threads. How does it guarantee INCR is still atomic?

2. **What happens if you INCR a key that holds "hello" (not a number)?**
   Try it. What error do you get? What does this tell you about Redis's type system?

3. **Can INCR overflow? What's the maximum value?**
   INCR operates on 64-bit signed integers. What's the max? What happens when you exceed it?

## Push Further

- What's `DECRBY`? Can you decrement by a negative number?
- Look up `INCRBYFLOAT` — can you do `INCRBYFLOAT mykey 0.1` ten times and expect exactly `1.0`? (Floating point!)
- What's the difference between using INCR as a counter vs using a sorted set score for ranking?
