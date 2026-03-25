# Gap Notebook — prod2

After completing this exercise, add to `feynman/gap_notebook.md`:

1. What happens if the Redis/Dragonfly server crashes while a lock is held? Is the lock recoverable from RDB snapshot?
2. Lock extension (refreshing TTL before it expires) requires checking you still hold the lock. What's the Lua script for this?
3. What's the difference between a distributed lock and a distributed semaphore? When would you need a semaphore?
4. Redlock uses N Redis nodes for stronger guarantees. How many nodes are needed for safety? What's the quorum rule?
5. If a Go process is GC-paused for 6 seconds and the lock TTL is 5 seconds, what's the worst case? How do fencing tokens prevent the damage?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. Your lock TTL is 5 seconds but the protected operation can take up to 30 seconds. Design a watchdog goroutine that refreshes the TTL before expiry — handle the case where the lock is lost mid-operation (another holder has taken it) and the watchdog must signal the original goroutine to abort.
2. Redlock acquires locks on N independent Redis nodes and requires N/2+1 successes. Implement the quorum-check logic in Go: given N responses (success, failure, or timeout), compute whether quorum was reached and calculate the effective lock validity time accounting for clock drift.
3. Distributed locks are advisory: a buggy client can ignore them. Design a fencing token scheme where the protected resource (e.g., a database row with a `last_lock_version` column) can enforce lock validity without trusting the lock holder.
