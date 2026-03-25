# Gap Notebook — expire2

After completing this exercise, add to `feynman/gap_notebook.md`:

1. With allkeys-lru, which key gets evicted? Is it always the LEAST recently used, or approximately?
2. What happens to a volatile-lru policy when there are no keys with TTLs but memory is full?
3. Can eviction be observed in real-time? Is there a metric?
4. If you set maxmemory but not maxmemory-policy, what is the default behavior?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- Dragonfly has `rss_oom_deny_ratio` — what does it do that `maxmemory` alone doesn't? Why would you want denials before the limit is reached?
- Run `INFO stats` on Dragonfly and find the eviction counter. What field name is it? Does it increment during your test?
- If `maxmemory-policy` is `allkeys-lru` and a hot key is accessed every 100ms, can it still be evicted? Under what exact condition?
