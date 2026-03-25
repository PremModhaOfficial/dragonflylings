# Gap Notebook — pipe1

After completing this exercise, add to `feynman/gap_notebook.md`:

1. What happens if the connection drops mid-pipeline? Do some commands execute and others not?
2. Is there a maximum number of commands you should put in a pipeline?
3. Can you mix reads and writes in a pipeline? What's the limitation?
4. Does pipelining help when Redis is on localhost (sub-millisecond RTT)?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- If network RTT is 1ms and you send 1000 commands, what's the theoretical wall-clock difference between individual vs. pipelined? Show your math.
- Does pipelining increase peak memory usage on the server? On the client? Which side buffers first?
- Can you pipeline commands that target different databases (interleaving SELECT)? What does Redis do?
