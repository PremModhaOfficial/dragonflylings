# Gap Notebook — expire3

After completing this exercise, add to `feynman/gap_notebook.md`:

1. Keyspace notifications are "fire-and-forget" — same as Pub/Sub. If your subscriber is down when a key expires, what happens?
2. What's the latency between a key expiring and the notification arriving?
3. notify-keyspace-events is disabled by default. Why? (Hint: think about production overhead)
4. Can you subscribe to ALL events for ALL databases at once? What would the channel pattern look like?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- If Dragonfly restarts, do active subscribers need to re-subscribe? Is the `notify_keyspace_events` setting persistent across restarts?
- Dragonfly currently only supports `Ex` for keyspace events. What does Redis support that Dragonfly doesn't? (Check Dragonfly's compatibility docs)
- Why subscribe to `__keyevent@0__:expired` (keyevent) instead of `__keyspace@0__:mykey` (keyspace) for expiry monitoring? What's the semantic difference?
