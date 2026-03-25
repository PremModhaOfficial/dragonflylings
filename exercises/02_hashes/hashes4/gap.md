# Gap Notebook — hashes4

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **What's the hash-max-ziplist-entries config?**
   Format: `- [ ] [your question] -- [why it matters]`
   Small hashes use a compact encoding. After 128 fields (default), Redis switches to a hash table. This affects memory significantly. What's the Dragonfly default?

2. **Can you store nested objects in a hash field?**
   Redis hash fields are strings. If a user has an "address" object with 5 sub-fields, how do you store it in a hash?

3. **What's HSETNX? How is it different from HSET?**
   When would you want to set a hash field only if it doesn't already exist?

## Push Further

- Benchmark: Create a hash with 100 fields vs 100 string keys. Compare `MEMORY USAGE` output.
- What happens to `HGetAll` when a hash has 100,000 fields? Is there a safer alternative?
- Look up "Redis hash encoding" — understand listpack vs hashtable encoding and the threshold between them.
