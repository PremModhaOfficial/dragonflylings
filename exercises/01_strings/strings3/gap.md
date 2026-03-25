# Gap Notebook — strings3

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **What happens if the lock holder crashes before calling ReleaseLock?**
   Format: `- [ ] [your question] -- [why it matters]`
   The TTL handles this — but what if the operation takes longer than the TTL?

2. **Can ReleaseLock accidentally release another process's lock?**
   Imagine: worker-1 acquires lock, takes too long, lock expires, worker-2 acquires lock, worker-1 finishes and calls ReleaseLock. Who does it release?

3. **What is Redlock and why was it controversial?**
   Single-node SETNX is fine for most cases. What problem does Redlock solve?

## Push Further

- How do you implement "lock with owner check" (only release if you own it)?
- Look up `SET key value NX EX seconds` — this is the modern alternative to SETNX. What's the advantage?
- What's a "fencing token" and why do distributed lock experts recommend it?
