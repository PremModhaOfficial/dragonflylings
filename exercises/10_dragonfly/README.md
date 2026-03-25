# Module 10: Dragonfly-Specific Features - "The Dragon"

## Mental Model

Dragonfly looks exactly like Redis from the outside. Same protocol (RESP). Same commands. Your go-redis client doesn't know the difference. But inside, they're built completely differently.

Redis is single-threaded. One cook in one kitchen, handling every order in sequence. This makes Redis simple to reason about (no lock contention, no parallelism bugs) but means it can't use more than one CPU core for command processing.

Dragonfly is multi-threaded with a shared-nothing architecture. Many cooks, each owning their own section of the kitchen (a shard of the keyspace). Commands that touch keys in the same shard run in parallel with commands in other shards. This is how Dragonfly can saturate all CPU cores and achieve 25x Redis throughput on the same hardware.

The same-menu-different-kitchen metaphor breaks down in one critical way: some Redis assumptions about single-threaded ordering no longer hold. Cross-shard operations (like transactions touching keys in different shards) require coordination. Lua scripts touching multiple shards need hashtag hints. `WAIT` semantics differ. And some Redis operational patterns (Sentinel, certain `OBJECT ENCODING` values) simply don't exist in Dragonfly's kitchen. This module shows you where the seams are.

## Before You Begin

**PREDICT:** Before writing any code, answer this:

> Redis is single-threaded, so `SET a 1; SET b 2; SET c 3` executes in strict order with no concurrency. If you run 1000 of these `SET` commands from 10 goroutines simultaneously against Dragonfly, what guarantees does Dragonfly make about ordering? About atomicity of each individual SET? What could be different from Redis?

Write your prediction before starting dragon1.

**Try it first** — run these to feel the environment:

```bash
# Check Dragonfly version and thread count:
redis-cli -p 6380 INFO server | grep -E "dragonfly_version|thread_count|redis_version"

# See internal encoding for different data types:
redis-cli -p 6380 SET enc:string "hello"
redis-cli -p 6380 OBJECT ENCODING enc:string
redis-cli -p 6380 ZADD enc:zset 1.0 member
redis-cli -p 6380 OBJECT ENCODING enc:zset

# Trigger a snapshot and watch its progress:
redis-cli -p 6380 BGSAVE
redis-cli -p 6380 INFO persistence | grep -E "rdb_bgsave_in_progress|rdb_last_bgsave_status"

# Try native JSON (no module needed in Dragonfly):
redis-cli -p 6380 JSON.SET user:1 '$' '{"name":"alice","score":100}'
redis-cli -p 6380 JSON.GET user:1 '$.name'
```

## Key Concepts

- Dragonfly uses shard-per-thread: keyspace is split across CPU-pinned threads
- Cross-shard commands (MSET with keys on different shards) use internal coordination
- `BGSAVE` in Dragonfly: snapshot without `fork()` — no memory spike, no copy-on-write
- `JSON.SET key path value` / `JSON.GET key path` — native JSON support, no module needed
- `OBJECT ENCODING key` — may return different values than Redis (Dragonfly uses different internal types)
- `WAIT numreplicas timeout` — behavior differs; Dragonfly currently runs standalone
- No Redis Sentinel support; Dragonfly uses its own HA approach
- Hashtag `{tag}` in keys forces keys to the same shard

## What You'll Build

| Exercise | What you fix | Feynman challenge |
|----------|-------------|-------------------|
| dragon1 | Parallel goroutines: observe multi-core throughput | Explain why Redis can't do this |
| dragon2 | BGSAVE: trigger and observe memory behavior | Predict memory usage during Redis vs. Dragonfly snapshot |
| dragon3 | JSON.SET/JSON.GET natively | Push it: what JSON path operations are supported? |
| dragon4 | Identify and work around Dragonfly gotchas | Explain each gotcha: why it exists, how to handle it |

## Resources

- [Dragonfly architecture overview](https://www.dragonflydb.io/docs/architecture)
- [Dragonfly vs Redis comparison](https://www.dragonflydb.io/docs/managing-dragonfly/compatibility)
- [Dragonfly JSON support](https://www.dragonflydb.io/docs/data-types/json)
- [Dragonfly BGSAVE](https://www.dragonflydb.io/docs/managing-dragonfly/backups)
- [Redis OBJECT ENCODING](https://redis.io/commands/object-encoding/)
