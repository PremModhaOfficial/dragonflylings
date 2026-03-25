# Gap Notebook — streams5

After completing this exercise, add to `feynman/gap_notebook.md`:

1. MAXLEN trims by COUNT. Is there a way to trim by AGE (delete entries older than X minutes)?
2. If you add 1000 entries with MAXLEN~100, could the stream temporarily have 150 entries?
3. Does XTRIM interact with consumer groups? What happens to unACKed PEL entries that get trimmed?
4. What is MINID and how does it differ from MAXLEN for trimming?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- XTRIM with `~` (approximate) is much faster than exact trimming. Why? What data structure optimization does this exploit?
- If you trim away a stream entry that's still in a consumer group's PEL (unACKed), what happens when that consumer tries to XACK the now-gone ID?
- MINID trims by ID (timestamp) rather than count. When is this more useful than MAXLEN? Give a time-based retention policy example.
