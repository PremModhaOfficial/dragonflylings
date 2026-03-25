# Gap Notebook — lua3

After completing this exercise, add to `feynman/gap_notebook.md`:

1. What happens to EVALSHA'd scripts when Dragonfly restarts — are they persisted to the RDB snapshot?
2. `SCRIPT FLUSH` clears all cached scripts. When would you need to call this in production?
3. Is the SHA1 of a script deterministic across Redis versions and Dragonfly? Can you hardcode the SHA in your config?
4. Can two different scripts have the same SHA? (What would that imply about SHA1?)
5. What is `SCRIPT EXISTS`? When would you use it over just catching NOSCRIPT errors?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. EVALSHA fails with NOSCRIPT after a server restart. Design a Go wrapper that transparently falls back to EVAL on NOSCRIPT, re-caches the SHA, and retries — without the caller ever seeing the error.
2. If you have 50 different Lua scripts deployed across a microservice fleet, how do you manage script lifecycle? Is there an analogy to database schema migrations? What's your rollback strategy?
3. Lua scripts are compiled on first EVAL. Does Dragonfly use LuaJIT or standard Lua? What are the performance implications of many small scripts vs a few large composite scripts?
