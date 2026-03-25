# Gap Notebook — pubsub2

After completing this exercise, add to `feynman/gap_notebook.md`:

1. Can you PSUBSCRIBE to "?" (single character wildcard)? What channels would that match?
2. If you both Subscribe to "news" and PSubscribe to "news.*", do you get duplicate messages for "news"?
3. What is the overhead of having 10,000 PSUBSCRIBE patterns vs 10,000 exact SUBSCRIBE channels?
4. The PMessage has both Channel and Pattern fields -- when would you need both?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- Does Redis/Dragonfly support character classes in PSUBSCRIBE patterns, e.g. `news.[abc]`? Try it.
- If you PSUBSCRIBE to `*`, do you receive every published message globally? What are the CPU implications at 100k msg/sec?
- Can you both SUBSCRIBE to `"news"` and PSUBSCRIBE to `"news.*"` on the same connection? Would a publish to `"news"` be delivered once or twice?
