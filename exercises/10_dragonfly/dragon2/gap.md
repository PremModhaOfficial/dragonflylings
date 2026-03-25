# Gap Notebook — dragon2

After completing this exercise, add to `feynman/gap_notebook.md`:

1. What happens if Dragonfly runs out of disk space during BGSAVE? How is this reported?
2. `LASTSAVE` returns a Unix timestamp of the last successful save. How would you use this to detect a stuck snapshot?
3. Does Dragonfly's forkless snapshot guarantee a point-in-time consistent snapshot? What's the consistency model?
4. What's the difference between RDB and AOF persistence? Does Dragonfly support both?
5. If `rdb_bgsave_in_progress` stays 1 indefinitely, what are the possible causes?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

1. Dragonfly's forkless snapshot avoids the "fork tax" Redis pays at high memory. What exactly is the fork tax (hint: copy-on-write page faults), and at what memory size does it become significant?
2. Design a backup strategy for Dragonfly: snapshot frequency, offsite storage, and restore validation — how do you verify a backup is restorable without taking down production?
3. If a snapshot is in progress and you must restart Dragonfly urgently (e.g., security patch), what's the safest procedure? What data can you lose, and what's the recovery path?
