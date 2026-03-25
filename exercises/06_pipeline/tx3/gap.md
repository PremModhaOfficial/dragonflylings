# Gap Notebook — tx3

After completing this exercise, add to `feynman/gap_notebook.md`:

1. Lua scripts block the entire Redis server while running. When does this become a problem?
2. Is there a scenario where you'd use Pipeline INSIDE a Watch/TxPipelined?
3. MULTI/EXEC doesn't support conditional logic (no if/else in a transaction). Lua does. When specifically does this matter?
4. What is a "Lua script cache" in Redis and why does it exist?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- Lua scripts are cached server-side by SHA1. What happens to the cache when Dragonfly restarts? How does SCRIPT LOAD interact with this?
- Can a Lua script call external services or make network calls? What is it explicitly prohibited from doing?
- In a Redis Cluster, a Lua script can only access keys in the same hash slot. How does this constraint change your script design?
