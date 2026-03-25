## Hint 1

`CreateGroup` uses `"$"` as starting ID. In XGROUP CREATE, `"$"` means "start from the latest entry — only process messages added AFTER this group is created." To process ALL existing messages, use `"0"` (before any entry).

## Hint 2

`ReadGroup` uses `"0"` in the Streams slice. In XREADGROUP:
- `">"` = give me new, undelivered messages (the normal read path)
- `"0"` = give me MY pending messages that I received but haven't ACKed yet (recovery path)

## Hint 3

`AckMessages` appends `":ack"` to the stream name. XAck takes the exact stream key. Fix: `client.XAck(ctx, stream, group, ids...)` — no suffix.

```go
// Three fixes:
client.XGroupCreateMkStream(ctx, stream, group, "0")   // was "$"
Streams: []string{stream, ">"}                          // was "0"
client.XAck(ctx, stream, group, ids...)                 // no ":ack" suffix
```
