# Module 06: Pipelining and Transactions - "The Assembly Line"

## Mental Model

Without pipelining, every Redis command is a separate letter: you write it, seal the envelope, walk to the mailbox, post it, walk home, and wait for a reply before writing the next letter. Each round trip to the server adds latency — typically 0.1ms to 1ms on a local network, but it adds up when you're sending 1000 commands.

Pipelining stuffs many letters in one envelope. You queue up 100 commands, send them in a single network write, and read all 100 responses in a single read. The server executes them in order. No waiting between commands. Latency drops from 100 × (round trip) to 1 × (round trip) + execution time.

Transactions (`MULTI/EXEC`) are a signed contract: all queued commands execute together as a unit, or none do. But this is not isolation in the SQL sense — other clients can still read keys between your `MULTI` and `EXEC`. For true optimistic concurrency, you add `WATCH`: "only execute this transaction if nobody modified key X since I watched it." If something changed, `EXEC` returns nil — your code detects this and retries. This is the foundation of compare-and-swap patterns in Redis.

The key lesson of this module: pipelining is about latency, transactions are about atomicity, and Lua scripts (next module) are about both at once. Knowing when to use each is a production skill.

## Before You Begin

**PREDICT:** Before writing any code, answer this:

> You pipeline 100 SET commands. Command #50 fails (the key is wrong type). What happens to commands 1-49? What happens to commands 51-100? Now contrast this with MULTI/EXEC — if command #50 fails in a transaction, what happens to the others?

Write your prediction. The answer surprises most people.

## Before You Start

```bash
# Verify Dragonfly is running
redis-cli -p 6380 PING

# Try a pipeline manually (--pipe mode sends raw RESP)
redis-cli -p 6380 SET counter 0
redis-cli -p 6380 MULTI
redis-cli -p 6380 INCR counter
redis-cli -p 6380 INCR counter
redis-cli -p 6380 EXEC

# Try WATCH/MULTI/EXEC for optimistic locking
redis-cli -p 6380 SET balance 100
redis-cli -p 6380 WATCH balance
redis-cli -p 6380 MULTI
redis-cli -p 6380 DECRBY balance 10
redis-cli -p 6380 EXEC
```

## Key Concepts

- `client.Pipeline()` — create a pipeline; commands are queued, not sent
- `pipe.Exec(ctx)` — flush all queued commands, return slice of results
- `MULTI` / `EXEC` — begin and commit a transaction
- `DISCARD` — cancel a queued transaction
- `WATCH key [key ...]` — optimistic lock; `EXEC` returns nil if any watched key changed
- Pipeline errors: each command has its own error; pipeline itself doesn't fail if one command fails
- Transaction errors: compile-time (queuing) vs. runtime (execution) error behavior differs
- `TxPipelined` in go-redis: wraps pipeline in MULTI/EXEC automatically

## What You'll Build

| Exercise | What you fix | Feynman challenge |
|----------|-------------|-------------------|
| pipe1 | Pipeline 100 SETs, measure vs. sequential | Predict the speedup factor |
| pipe2 | Handle partial pipeline failures correctly | Explain pipeline error semantics |
| tx1 | MULTI/EXEC atomic balance transfer | Explain why MULTI doesn't mean isolated |
| tx2 | WATCH + retry loop for optimistic locking | Push it: at what contention level does WATCH break down? |
| tx3 | Choose between pipeline, transaction, and Lua | Explain the decision tree |

## Resources

- [Redis pipelining docs](https://redis.io/docs/manual/pipelining/)
- [Redis transactions docs](https://redis.io/docs/manual/transactions/)
- [MULTI command](https://redis.io/commands/multi/)
- [WATCH command](https://redis.io/commands/watch/)
- [Dragonfly transactions](https://www.dragonflydb.io/docs/command-reference/transactions)
