# Gap Notebook — strings5

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **Does MSET execute atomically?**
   Format: `- [ ] [your question] -- [why it matters]`
   If MSET sets 5 keys and your app crashes mid-execution, are some keys set and others not?

2. **What's the maximum number of keys you can pass to MGET?**
   Is there a hard limit? What happens to latency if you pass 10,000 keys?

3. **Why doesn't MSET support TTL?**
   MSET doesn't have an expiry option. If you need to SET many keys each with different TTLs, what's your option?

## Push Further

- Look up `MSETNX` — it's MSET but only sets if ALL keys don't exist. How is that useful?
- What's pipelining and how does it differ from MSET/MGET?
- Try benchmarking: 100 individual SETs vs 1 MSET with 100 keys. What ratio do you see?
