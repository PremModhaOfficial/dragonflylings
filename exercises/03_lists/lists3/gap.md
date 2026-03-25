# Gap Notebook — lists3

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **Can BLPOP watch multiple queues at once?**
   Format: `- [ ] [your question] -- [why it matters]`
   `BLPOP key1 key2 key3 timeout` watches all three. Which queue takes priority?

2. **What happens if you close the connection while BLPOP is blocking?**
   The worker is mid-block. The TCP connection drops. What happens to the job that arrives?

3. **How many workers can share the same BLPOP queue?**
   If 5 workers all BLPOP the same queue, and 1 job arrives, which worker gets it?

## Push Further

- What's `BRPOPLPUSH`? It's deprecated in Redis 6.2 — what replaced it?
- Look up `LMOVE` and `BLMOVE`. How do they enable "reliable queue" patterns (at-least-once processing)?
- Compare BLPOP to Redis Streams consumer groups. When would you prefer each?
