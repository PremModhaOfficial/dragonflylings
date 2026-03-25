# Module 12: Capstone - "The Summit"

## Mental Model

No new commands. No new concepts. No hints.

You've learned every primitive: Strings, Hashes, Lists, Sets, Sorted Sets, Streams, Pub/Sub, Lua, pipelining, transactions, expiration, Dragonfly specifics, and production patterns. The capstone exercises combine them.

The summit analogy is deliberate. You've been climbing all the techniques as separate faces of a mountain. Now you stand at the top and can see how they connect. A real-time leaderboard needs Streams to ingest events, Lua for atomic score updates, and Sorted Sets for rankings. A production API cache needs cache-aside for reads, rate limiting for protection, a circuit breaker for resilience, and Pub/Sub to invalidate stale entries across service instances.

These are not toy exercises. They're simplified versions of systems that run in production. If you can build them here, you can build them at work.

**There are no hints in this module.** If you're stuck, go back to the module where the relevant concept lives, re-read the README, re-do the exercise. The hints are in the work you've already done.

## Before You Begin

**REFLECT:** Before starting, answer this:

> Look back at your `feynman/gap_notebook.md`. What questions did you write down in modules 00-11 that you still haven't answered? Pick the two most important ones and answer them now — either from memory, from the docs, or by writing a small test. Then start the capstone.

This isn't busywork. Unresolved gaps become production bugs.

**Warm up** — run the core commands before writing any code:

```bash
# Streams (capstone1):
redis-cli -p 6380 XADD cap:events '*' player alice points 100
redis-cli -p 6380 XADD cap:events '*' player bob points 200
redis-cli -p 6380 XREAD COUNT 10 STREAMS cap:events 0

# Consumer groups:
redis-cli -p 6380 XGROUP CREATE cap:events processors 0 MKSTREAM
redis-cli -p 6380 XREADGROUP GROUP processors worker-1 COUNT 10 STREAMS cap:events '>'
redis-cli -p 6380 XPENDING cap:events processors - + 10

# Leaderboard sorted set:
redis-cli -p 6380 ZINCRBY cap:leaderboard 100 alice
redis-cli -p 6380 ZINCRBY cap:leaderboard 200 bob
redis-cli -p 6380 ZREVRANGE cap:leaderboard 0 4 WITHSCORES

# Pub/Sub invalidation (capstone2) — open two terminals:
# Terminal 1: redis-cli -p 6380 SUBSCRIBE cache:invalidate:products
# Terminal 2: redis-cli -p 6380 PUBLISH cache:invalidate:products "product:42"
```

## Key Concepts

All of them. Specifically, each capstone exercise will require:

**capstone1** — Real-time event pipeline:
- `XADD` / `XREADGROUP` / `XACK` for event ingestion and processing
- Consumer groups for parallel worker processing
- Lua scripts for atomic leaderboard updates
- Sorted Sets for ranked leaderboard queries

**capstone2** — Rate-limited API cache:
- Cache-aside pattern with singleflight
- Sliding window rate limiter (Sorted Sets)
- Circuit breaker wrapping the client
- Pub/Sub for cross-instance cache invalidation

## What You'll Build

| Exercise | System | No hints |
|----------|--------|----------|
| capstone1 | Real-time event processing pipeline with leaderboard | ✓ |
| capstone2 | Rate-limited API cache with circuit breaker and invalidation | ✓ |

## Resources

You know where to look. Here's the map:

- Module 06: pipelining and transactions
- Module 08: Streams and consumer groups
- Module 09: Lua scripting
- Module 04: Sorted Sets for leaderboards
- Module 11: cache-aside, rate limiter, circuit breaker, Pub/Sub invalidation
- [Redis patterns library](https://redis.io/docs/manual/patterns/)
- [Dragonfly docs](https://www.dragonflydb.io/docs)
