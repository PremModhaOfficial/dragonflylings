# Module 11: Production Patterns - "The Battlefield"

## Mental Model

Knowing Redis commands is like knowing how chess pieces move. You can describe every rule perfectly. Production is playing actual games under time pressure, against an opponent who makes unexpected moves, while the clock ticks.

This module is the jump from chess rules to chess games. The commands are the same — `SET`, `GET`, `ZADD`, `SETNX` — but now they're assembled into patterns that handle real failure modes: cache stampedes when the database is slow, lock expiry when a process hangs, cascading failures when Dragonfly is degraded, hot keys that overwhelm a single shard.

Every pattern here is battle-tested and appears in production codebases at scale. Cache-aside with singleflight prevents the "thundering herd" problem where 1000 requests simultaneously discover a cold cache and all hit the database. The sliding window rate limiter is used by every major API. The circuit breaker ensures that when Dragonfly is unhealthy, your application degrades gracefully instead of cascading into total failure. Learn these patterns now; you'll use them the rest of your career.

## Before You Begin

**PREDICT:** Before writing any code, answer this:

> A popular blog post's cache entry expires at exactly 3:00 PM. At 3:00:00.001 PM, 500 simultaneous requests arrive and all find the cache cold. What happens next? How many database queries will execute? How would you prevent this without simply making the cache entry never expire?

Write your prediction. This is the thundering herd problem, and you'll solve it in prod1.

**Try it first** — the patterns by hand before writing Go:

```bash
# Cache-aside by hand:
redis-cli -p 6380 SET cache:article:1 '{"title":"hello"}' EX 60
redis-cli -p 6380 TTL cache:article:1

# Distributed lock by hand:
redis-cli -p 6380 SET lock:resource token-abc NX EX 30
redis-cli -p 6380 GET lock:resource
# Only delete if you still own it (check token first):
redis-cli -p 6380 DEL lock:resource

# Sliding window rate limit by hand:
redis-cli -p 6380 ZADD ratelimit:user1 $(date +%s%3N) "req-1"
redis-cli -p 6380 ZCARD ratelimit:user1
redis-cli -p 6380 ZREMRANGEBYSCORE ratelimit:user1 -inf $(($(date +%s%3N) - 60000))

# Circuit breaker state — track manually:
redis-cli -p 6380 SET cb:failures 0
redis-cli -p 6380 INCR cb:failures
redis-cli -p 6380 GET cb:failures
```

## Key Concepts

- Cache-aside: check cache → on miss, fetch from DB → populate cache → return result
- Singleflight (`golang.org/x/sync/singleflight`): deduplicate concurrent cache misses
- Distributed lock: `SET lock:key uuid NX EX 30` — acquire with expiry; `DEL` only if value matches uuid
- Lock extension: refresh TTL before expiry if work takes longer than expected
- Sliding window rate limiter: `ZADD` timestamp → `ZREMRANGEBYSCORE` old entries → `ZCARD` count
- Session storage: Hash per session with `EXPIRE`; degrade gracefully if Redis unavailable
- Hot key mitigation: local in-process cache (sync.Map) as L1 cache in front of Redis
- Circuit breaker: track failure rate; open circuit → return fallback; half-open → probe

## What You'll Build

| Exercise | What you fix | Feynman challenge |
|----------|-------------|-------------------|
| prod1 | Cache-aside + singleflight thundering herd protection | Explain why singleflight alone isn't enough |
| prod2 | Distributed lock with expiry + safe release | Push it: what breaks if two processes have clock skew? |
| prod3 | Sliding window rate limiter with sorted sets | Explain why sliding window is better than fixed window |
| prod4 | Hash-based sessions with graceful degradation | Predict behavior when Dragonfly is unreachable |
| prod5 | Hot key mitigation with local L1 cache | Explain the consistency tradeoff in L1 caching |
| prod6 | Circuit breaker wrapping the Redis client | Explain the three circuit states and transitions |

## Resources

- [Redis cache patterns](https://redis.io/docs/manual/patterns/caching/)
- [Distributed locks with Redis](https://redis.io/docs/manual/patterns/distributed-locks/)
- [Rate limiting patterns](https://redis.io/glossary/rate-limiting/)
- [golang.org/x/sync/singleflight](https://pkg.go.dev/golang.org/x/sync/singleflight)
- [Dragonfly production guide](https://www.dragonflydb.io/docs/managing-dragonfly/tuning)
