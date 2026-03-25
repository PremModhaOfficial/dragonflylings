# Gap Notebook — streams1

After completing this exercise, add to `feynman/gap_notebook.md`:

1. Stream IDs encode time -- can you use this to query "all events from the last hour" without a timestamp field?
2. What happens if the system clock goes backward? Are stream IDs still monotonic?
3. XREAD with Block=0 is non-blocking. What if you want to wait for NEW messages indefinitely?
4. What is the maximum number of fields a single stream entry can have?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- If your consumer crashes and restarts, how do you know which stream ID to resume reading from? Where would you persist the last-seen ID?
- Two events XADD'd within the same millisecond get IDs like `1711234567890-0` and `1711234567890-1`. What is the maximum sequence number within one millisecond?
- XREVRANGE reads a stream in reverse order. When would you prefer it over XRANGE — can you give a concrete use case?
