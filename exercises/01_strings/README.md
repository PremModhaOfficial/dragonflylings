# Module 01: Strings — "The Atoms"

## Mental Model

Strings are Redis's atoms — the indivisible unit everything else builds on. But "string" is misleading: they're really "binary-safe byte sequences up to 512MB." A string can hold a serialized protobuf as easily as the word "hello."

Everything in Redis that isn't a collection type is stored as a String. Numbers are stored as strings too, but Redis knows how to operate on them atomically (INCR, INCRBYFLOAT).

```
Key                     Value (always a string internally)
─────────────────────   ────────────────────────────────
"user:42:name"          "alice"
"session:abc123"        "eyJhbGciOiJSUzI1..."  (JWT token)
"counter:pageviews"     "10483"               (number as string)
"feature:dark_mode"     "true"                (bool as string)
"user:42:avatar"        "\x89PNG\r\n..."      (binary blob)
```

## Predict Before Starting

Before writing any code, answer in your head:
1. Is `INCR counter` the same as GET + parse + add 1 + SET?
2. What does GET return for a key that doesn't exist?
3. SET stores a string. Can it store a number? How does INCR work then?
4. How many round trips does storing 5 keys individually take vs MSET?

Write your predictions in `feynman/gap_notebook.md`.

## Key Concepts

| Command | Description |
|---------|-------------|
| SET key value \[EX seconds\] | Store a string, optionally with expiry |
| GET key | Retrieve a string (redis.Nil if missing) |
| SETNX key value | Set only if key doesn't exist (foundation of locks) |
| INCR / DECR | Atomic increment/decrement |
| INCRBY / DECRBY | Atomic increment/decrement by N |
| MSET / MGET | Batch set/get in one round trip |
| TTL key | Get remaining time-to-live (-1 = no expiry, -2 = gone) |

## Exercises

- **strings1**: Fix key construction — SET stores under a different key than GET reads
- **strings2**: Add expiration — SET with TTL so sessions automatically expire
- **strings3**: SETNX for distributed locking — set only if not exists
- **strings4**: INCR atomicity — why GET+math+SET loses updates under concurrency
- **strings5**: MGET/MSET batch operations — reduce N round trips to 1

## Before You Start

```bash
# Verify Dragonfly is running
redis-cli -p 6380 PING

# Experiment with strings manually
redis-cli -p 6380 SET mykey "hello"
redis-cli -p 6380 GET mykey
redis-cli -p 6380 TTL mykey   # -1 means no expiry
```
