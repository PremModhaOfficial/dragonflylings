# Gap Notebook — tx2

After completing this exercise, add to `feynman/gap_notebook.md`:

1. At what level of concurrency does optimistic locking become worse than a mutex?
2. WATCH watches specific keys -- what if you need to watch the result of a computation, not a key?
3. Can you WATCH and then do non-transactional reads before TxPipelined? What are the risks?
4. What is a "fencing token" and how does it solve a problem that WATCH cannot?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- WATCH expires after EXEC is called (whether the transaction succeeds or not). Does a successful EXEC also unwatch all keys automatically?
- Under extreme contention (100 goroutines all WATCH the same key), your retry loop could spin indefinitely. How would you add an exponential backoff with a max retry limit?
- WATCH is per-connection. If two goroutines share one `*redis.Client`, do their WATCH calls interfere? Why or why not?
