# Module 08: Streams - "The Log"

## Mental Model

If Pub/Sub is a megaphone in a stadium, a Stream is a recorded lecture. Every message is appended to an immutable log with a unique ID (millisecond-timestamp + sequence number). You can replay from any point in history. Late consumers can catch up. The log persists until you explicitly trim it.

Consumer Groups make Streams a distributed work queue. Think of a study group watching the recorded lecture: every group watches every message, but within a group, different students handle different parts — student A handles messages 1, 3, 5 while student B handles 2, 4, 6. No message is delivered to two students in the same group. This is how you fan work out across multiple workers while ensuring each item is processed exactly once.

The Pending Entries List (PEL) is what separates Streams from Pub/Sub for reliability. When a consumer reads a message, it goes into the PEL — a "you acknowledged this right?" list. If the consumer crashes before calling `XACK`, the message stays in the PEL. Another consumer can call `XPENDING` to find abandoned messages and `XCLAIM` to take ownership. This is Redis's answer to "at-least-once delivery."

## Before You Begin

**PREDICT:** Before writing any code, answer this:

> A Stream has 1000 messages. Consumer group "workers" has 3 consumers. Worker-1 reads messages 1-10 but crashes before ACKing them. Worker-2 and Worker-3 read messages 11-1000 and ACK all of them. Now Worker-1 restarts. What does it see when it reads from the stream? How does it know about the unACKed messages?

Write your prediction before starting streams3.

## Before You Start

```bash
# Verify Dragonfly is running
redis-cli -p 6380 PING

# Try streams yourself:
redis-cli -p 6380 XADD mystream '*' event "login" user "alice"
redis-cli -p 6380 XADD mystream '*' event "purchase" user "bob"
redis-cli -p 6380 XLEN mystream          # number of entries
redis-cli -p 6380 XRANGE mystream - +    # read all entries
redis-cli -p 6380 XREAD COUNT 10 STREAMS mystream 0  # read from start

# Create a consumer group:
redis-cli -p 6380 XGROUP CREATE mystream workers 0 MKSTREAM
redis-cli -p 6380 XREADGROUP GROUP workers consumer1 COUNT 10 STREAMS mystream '>'
```

## Key Concepts

- `XADD key * field value [field value ...]` — append message (auto-ID)
- `XREAD COUNT n STREAMS key id` — read n messages after given ID (`0` = from start, `$` = new only)
- `XRANGE key start end` — read messages in ID range (`-` = min, `+` = max)
- `XLEN key` — number of messages in stream
- `XGROUP CREATE key group id MKSTREAM` — create consumer group
- `XREADGROUP GROUP group consumer COUNT n STREAMS key >` — read undelivered messages
- `XACK key group id` — acknowledge processed message (removes from PEL)
- `XPENDING key group - + count` — list unacknowledged messages
- `XCLAIM key group consumer min-idle-time id` — steal ownership of idle message
- `XTRIM key MAXLEN count` — cap stream length

## What You'll Build

| Exercise | What you fix | Feynman challenge |
|----------|-------------|-------------------|
| streams1 | XADD/XREAD: produce and consume a stream | Predict what ID format looks like |
| streams2 | XRANGE by time range, XLEN | Explain stream IDs as logical timestamps |
| streams3 | Consumer groups: XGROUP + XREADGROUP + XACK | Explain what the `>` means in XREADGROUP |
| streams4 | PEL recovery: XPENDING + XCLAIM | Push it: what's the right idle timeout before reclaiming? |
| streams5 | XTRIM MAXLEN for bounded streams | Explain MAXLEN ~ vs. exact: the tradeoff |

## Resources

- [Redis Streams docs](https://redis.io/docs/data-types/streams/)
- [XADD command](https://redis.io/commands/xadd/)
- [XREADGROUP command](https://redis.io/commands/xreadgroup/)
- [XPENDING command](https://redis.io/commands/xpending/)
- [Dragonfly Streams](https://www.dragonflydb.io/docs/data-types/streams)
