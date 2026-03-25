# Gap Notebook — tx1

After completing this exercise, add to `feynman/gap_notebook.md`:

1. MULTI/EXEC doesn't roll back. If DecrBy succeeds but IncrBy fails (can it?), what happens to the money?
2. Redis documentation says transactions are "atomic" -- but what does that mean exactly?
3. Can you WATCH keys inside a TxPipelined? Or do you need a separate Watch call?
4. What happens if the client crashes between MULTI and EXEC?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- Redis says MULTI/EXEC is "atomic" but not "isolated." Concretely, can another client read a key mid-transaction? Write a thought experiment.
- If you call MULTI, queue 5 commands, then call DISCARD — what happens to those 5 commands? What does the server return?
- Can you nest MULTI inside MULTI? Try it with `redis-cli` and see the exact error message Redis returns.
