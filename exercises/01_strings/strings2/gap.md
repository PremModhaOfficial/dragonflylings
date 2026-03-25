# Gap Notebook — strings2

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **Does Redis check every key's TTL every second?**
   Format: `- [ ] [your question] -- [why it matters]`
   Hint: Redis uses lazy expiry + periodic sampling. An expired key might still exist in memory briefly after its TTL.

2. **What happens if you SET a key that already has a TTL?**
   Does the TTL reset? Does the old TTL persist? Try it.

3. **Can you extend a key's TTL without knowing its current value?**
   What command would you use to reset the TTL without touching the value?

## Push Further

- What's `PTTL` vs `TTL`? (Hint: precision)
- What's `EXPIREAT`? When would you use it over `EXPIRE`?
- In Redis 7+, what's `EXAT`? Does Dragonfly support it?
