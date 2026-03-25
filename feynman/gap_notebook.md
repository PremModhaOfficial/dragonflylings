# Gap Notebook

Your running log of things you don't fully understand yet.
Add an entry after each exercise. Return to it weekly.

Format: `- [ ] [your question] -- [why it matters]`

---

## Example Entries

- [ ] What happens to a Redis connection if the TCP socket drops mid-command? -- Understanding failure modes matters for retry logic
- [ ] Is INCR atomic across multiple Dragonfly threads? -- Dragonfly is multi-threaded, Redis is single-threaded; atomicity guarantees differ
- [ ] When does the ziplist encoding switch to hashtable for hashes? -- Encoding affects memory usage and access patterns in production

---

## Your Entries

### Prediction
- newClient is just the runtme entity not a connection 
- If you configure a pool of 10 connections, when are those connections established? -- whne the first  connection is established?

<!-- Add your gaps below as you work through exercises -->
