# Module 07: Pub/Sub - "The Megaphone"

## Mental Model

Pub/Sub is a megaphone in a stadium. The publisher grabs the megaphone and shouts a message. Everyone currently in the stadium who is tuned to that channel hears it. If you walked in after the shout, you missed it — there's no replay, no history, no acknowledgment. The message is gone.

This fire-and-forget nature makes Pub/Sub excellent for real-time broadcasts: "user X came online," "new post published in channel Y," "cache invalidation signal." It makes Pub/Sub terrible for reliable messaging where you need "every consumer must process every message at least once." For that, use Streams (next module).

Pattern subscriptions (`PSUBSCRIBE`) let you subscribe to a wildcard: `notifications.*` matches `notifications.email`, `notifications.push`, `notifications.sms`. This is how you build a single consumer that handles a family of related channels without listing each one.

One critical operational note: a connection in subscribe mode can only send subscribe/unsubscribe commands — it cannot send regular Redis commands. In go-redis, you need a separate client for publishing vs. subscribing, or you'll get "connection in subscribe mode" errors.

## Before You Begin

**PREDICT:** Before writing any code, answer this:

> Subscriber A subscribes to channel "news". Publisher B sends 5 messages. Subscriber C then subscribes. Publisher B sends 5 more messages. How many messages does each subscriber receive? Now: what happens if the network drops subscriber A's connection for 30 seconds and reconnects? How many messages did A miss?

Write your prediction before starting pubsub1.

## Before You Start

```bash
# Verify Dragonfly is running
redis-cli -p 6380 PING

# Try Pub/Sub in two terminals:
# Terminal 1 — subscribe:
redis-cli -p 6380 SUBSCRIBE news

# Terminal 2 — publish (open a new terminal):
redis-cli -p 6380 PUBLISH news "hello world"
# Terminal 1 should show the message immediately

# Try pattern subscribe:
redis-cli -p 6380 PSUBSCRIBE "notifications.*"
# Then publish: redis-cli -p 6380 PUBLISH notifications.email "test"
```

## Key Concepts

- `client.Subscribe(ctx, "channel")` — subscribe; returns a `*PubSub` object
- `pubsub.ReceiveMessage(ctx)` — block until next message arrives
- `client.Publish(ctx, "channel", payload)` — publish to channel
- `client.PSubscribe(ctx, "pattern*")` — pattern subscribe with glob matching
- Subscribe mode: the connection is dedicated; no regular commands allowed
- Fire-and-forget: no message persistence, no delivery guarantee
- Messages missed while disconnected are gone forever
- go-redis auto-reconnects on subscribe connections, but missed messages aren't replayed

## What You'll Build

| Exercise | What you fix | Feynman challenge |
|----------|-------------|-------------------|
| pubsub1 | SUBSCRIBE/PUBLISH with two goroutines | Predict what happens with no active subscribers |
| pubsub2 | PSUBSCRIBE with wildcard channel matching | Explain when you'd use pattern vs. direct subscribe |
| pubsub3 | Prove fire-and-forget: messages lost on disconnect | Push it: can you build reliable pub/sub on top of this? |

## Resources

- [Redis Pub/Sub docs](https://redis.io/docs/manual/pubsub/)
- [SUBSCRIBE command](https://redis.io/commands/subscribe/)
- [PSUBSCRIBE command](https://redis.io/commands/psubscribe/)
- [Dragonfly Pub/Sub](https://www.dragonflydb.io/docs/command-reference/pubsub)
