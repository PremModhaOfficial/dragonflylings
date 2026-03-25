# Gap Notebook — connect1

After completing this exercise, add these questions to `feynman/gap_notebook.md`:

1. **What's one thing about Redis connections you're still unsure about?**
   Format: `- [ ] [your question] -- [why it matters]`
   Example: `- [ ] Does go-redis reconnect automatically? -- I don't know what happens if Dragonfly restarts mid-session`

2. **What would happen if Dragonfly was under heavy load during your PING?**
   Would PING still return PONG? Would it be slow? Would it error? Think about what "heavy load" means for a single-threaded server vs a multi-threaded one.

3. **Do you know the difference between a TCP connection and the Redis protocol handshake?**
   TCP is layer 4 (transport). Redis RESP is the application protocol on top. They're separate. A TCP connection can succeed even if the Redis server is in a bad state.

## Push Further

- What port does Redis use by default, and why 6379? (Look it up — the answer is delightfully human)
- Can you have multiple clients talking to Dragonfly simultaneously? What limits them?
- What does `client.Close()` actually do? Does it immediately disconnect?
