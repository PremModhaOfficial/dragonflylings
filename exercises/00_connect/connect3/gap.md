# Gap Notebook — connect3

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **What happens when all pool connections are in use and a new request arrives?**
   Format: `- [ ] [your question] -- [why it matters]`
   Hint: look at `PoolTimeout` in the options. What error do you get? How do you handle it?

2. **Is there a maximum useful PoolSize?**
   At some point, more connections don't help. What limits you — the client, the network, or Dragonfly itself?

3. **What's the cost of a connection that's never used?**
   `MinIdleConns` keeps connections open even when idle. Each costs memory and a file descriptor. Is there a `MaxIdleConns` or similar limit?

## Push Further

- What does `PoolStats.Misses` count? When would you see misses in production?
- Dragonfly is multi-threaded (unlike Redis). Does that change how you'd configure pool size?
- Look up `ConnMaxLifetime` and `ConnMaxIdleTime`. Why would you want connections to expire?
