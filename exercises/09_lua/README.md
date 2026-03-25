# Module 09: Lua Scripting — "The Stored Procedure"

## Mental Model

Lua scripts in Redis are like stored procedures in SQL: they run **atomically on the server side**. Instead of 5 round trips (read, compute, write, read, write), you send one script that does all 5 steps without anyone else interfering.

```
Without Lua:                    With Lua:
Client → GET key       ←→      Client → EVAL script ←→
         (check value)                   (GET + compare
Client → SET key (maybe)                  + SET atomically)
```

**The key guarantee:** Between the first and last command in your Lua script, no other client can sneak in. It's a critical section baked into the server.

**The tradeoff:** Your script blocks the entire server while it runs. A 100ms Lua script is 100ms where no other client gets service. Keep scripts short and fast.

## Dragonfly vs Redis

In standard Redis (single-threaded), Lua scripts work fine with multi-key access. In Dragonfly (multi-threaded, sharded), keys that a Lua script touches must all live on the **same shard**. The way to control shard placement: **hashtags**.

- `user:42:balance` — shard determined by full key name (unpredictable)
- `{user:42}:balance` — shard determined by `user:42` only
- `{user:42}:lock` — same `{user:42}` tag → **same shard** ✓

Without hashtags, a multi-key Lua script in Dragonfly can fail with:
`ERR CROSSSLOT Keys in request don't hash to the same slot`

## The EVAL Protocol

```
EVAL script numkeys [key [key ...]] [arg [arg ...]]
         │      │         │                │
         │      │    Goes into KEYS[]   Goes into ARGV[]
         │      │    (1-indexed)        (1-indexed)
         │      └── Must declare ALL keys upfront
         └── The Lua script body
```

## Exercises

1. **lua1** — Compare-and-Swap: atomically "set if current value matches"
2. **lua2** — KEYS vs ARGV: why keys must be declared (not hardcoded in scripts)
3. **lua3** — EVALSHA: cache scripts by SHA to avoid resending code on every call
4. **lua4** — Dragonfly hashtag gotcha: multi-shard scripts require `{tag}` prefixes

## Before You Start

**Try it first** — run these before writing any Go:

```bash
# Evaluate a Lua script inline:
redis-cli -p 6380 EVAL "return redis.call('SET', KEYS[1], ARGV[1])" 1 test:key hello
redis-cli -p 6380 GET test:key

# Load a script and get its SHA (for EVALSHA):
redis-cli -p 6380 SCRIPT LOAD "return redis.call('GET', KEYS[1])"
# Paste the SHA printed above:
# redis-cli -p 6380 EVALSHA <sha-here> 1 test:key

# See what error looks like when redis.call() fails:
redis-cli -p 6380 SET test:notahash hello
redis-cli -p 6380 EVAL "return redis.call('HINCRBY', KEYS[1], 'field', 1)" 1 test:notahash
```

**PREDICT:** Without running any code, answer:
- What's the difference between `EVAL` and `EVALSHA`?
- Can a Lua script return multiple values?
- What happens if your Lua script calls `redis.call()` and it errors — does the script stop or continue?
