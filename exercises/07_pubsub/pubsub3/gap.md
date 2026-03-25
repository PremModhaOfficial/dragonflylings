# Gap Notebook — pubsub3

After completing this exercise, add to `feynman/gap_notebook.md`:

1. go-redis automatically reconnects subscriptions after disconnects -- but missed messages during downtime are lost. How would you know what you missed?
2. Can you build a "reliable Pub/Sub" on top of Pub/Sub + Streams? What would that look like?
3. PUBLISH returns the number of clients that received the message. Is this useful for ensuring delivery?
4. What happens to a subscriber that can't keep up with publish rate? Where does backpressure apply?

Format: `- [ ] [your question] -- [why it matters]`

## Push Further

- go-redis has a `PubSub.Ping()` method for long-lived idle subscriptions. Why does idle detection matter for a TCP connection in subscribe mode?
- If 10,000 subscribers are listening and you publish a 1MB message, estimate the memory Dragonfly needs just for fanout delivery.
- Sketch a "at-least-once" Pub/Sub pattern using Redis Pub/Sub + a persistent log (could be a Stream or database). Where does each guarantee break down?
