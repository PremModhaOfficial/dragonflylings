# Gap Notebook — hashes1

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **How many fields can a Redis hash hold?**
   Format: `- [ ] [your question] -- [why it matters]`
   Is there a limit? What happens to memory encoding when a hash gets large?

2. **Can hash fields have TTLs?**
   The whole key can have a TTL (EXPIRE). But can individual fields expire? What's the workaround?

3. **What does HGETALL return? What's the order?**
   Is it insertion order? Alphabetical? Random? Does it matter for your application?

## Push Further

- What's `HSCAN`? When would you use it instead of `HGETALL`?
- What's `HKEYS` and `HVALS`? When would you need just keys or just values?
- Look up "hash ziplist encoding" — Redis uses a compact encoding for small hashes. At what field count does it switch to a hash table?
