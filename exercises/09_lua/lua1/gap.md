# Gap Notebook — lua1

After completing this exercise, add to `feynman/gap_notebook.md`:

1. What happens if a Lua script calls `redis.call()` and Redis returns an error — does the entire script abort, or does it continue?
2. Can you use `pcall()` in Lua to catch Redis errors mid-script? What are the risks of doing so?
3. Lua scripts block the server while running. How long is "too long" for a production Lua script?
4. Is Dragonfly's CAS implementation truly atomic, or does its multi-threaded nature create a window?
5. CAS assumes you know the current value to compare against. What if the value was set by a different process and you never read it?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. What is optimistic locking? How does CAS relate to it, and how does it compare to pessimistic locking (SETNX + lock TTL)? When would you choose each?
2. `redis.call()` raises an error and aborts the script; `redis.pcall()` returns an error table you can inspect. Design a multi-step CAS that uses `pcall` to recover from a partial failure and still leave the key in a consistent state.
3. CAS is the foundation of many higher-level patterns. How would you implement a counter that saturates at MAX (never goes below zero) using only a single Lua CAS script — without a separate lock?
