# Module 05: Expiration and Memory - "The Janitor"

## Mental Model

Redis is an in-memory store. Without expiration, it's a bathtub with the drain plugged — eventually it overflows. TTL is the drain. You set an expiry and the key disappears when time runs out. Simple concept, subtle behavior.

The subtlety is *how* Redis actually deletes expired keys. It doesn't check every key every millisecond — that would be too expensive. Instead it uses two strategies together: **lazy deletion** (check if a key is expired when someone tries to access it — if yes, delete and return nil), and **periodic sampling** (every 100ms, sample 20 random volatile keys, delete the expired ones, repeat if more than 25% were expired). This means a key might exist slightly after its TTL fires if nobody accessed it. For most applications this is fine. For some (keyspace notifications, strict rate limiters) it matters.

When memory is completely full, Redis/Dragonfly can't just refuse writes — it would break clients. Instead, it evicts keys according to a policy. `noeviction` returns errors (safe but breaks the app). `allkeys-lru` evicts the least-recently-used key regardless of TTL (treats Redis as a pure cache). `volatile-lru` only evicts keys that have TTLs set (protects persistent data). Understanding these policies is the difference between a Redis that silently drops your data and one that fails loudly.

## Before You Begin

**PREDICT:** Before writing any code, answer this:

> You set a key with `EXPIRE mykey 10`. Ten seconds later you call `GET mykey`. You get nil. Was the key definitely deleted exactly at the 10-second mark? Could it have been deleted earlier? Could it still exist somewhere? What guarantees does TTL actually give you?

Write your prediction before starting expire1.

## Before You Start

```bash
# Verify Dragonfly is running with keyspace events enabled
redis-cli -p 6380 PING
# Should return: PONG

# Check that notify_keyspace_events is set (needed for expire3)
redis-cli -p 6380 CONFIG GET notify_keyspace_events
# Should return: Ex (if not, run: docker compose up -d --force-recreate)

# Try TTL yourself
redis-cli -p 6380 SET mykey hello EX 10
redis-cli -p 6380 TTL mykey      # returns seconds remaining
redis-cli -p 6380 PTTL mykey     # returns milliseconds remaining
redis-cli -p 6380 PERSIST mykey  # removes TTL
redis-cli -p 6380 TTL mykey      # now returns -1 (no expiry)
```

## Key Concepts

- `EXPIRE key seconds` — set TTL on existing key
- `PEXPIRE key milliseconds` — millisecond precision TTL
- `TTL key` — remaining TTL in seconds (-1 = no expiry, -2 = key gone)
- `PERSIST key` — remove TTL, make key permanent
- `SET key value EX seconds` — set value and TTL in one command (preferred)
- Lazy deletion: expired keys deleted on access
- Active expiry: periodic background sampling of volatile keys
- Eviction policies: `noeviction`, `allkeys-lru`, `volatile-lru`, `allkeys-random`, `volatile-ttl`
- Keyspace notifications: subscribe to `__keyevent@0__:expired` for expiry events

## What You'll Build

| Exercise | What you fix | Feynman challenge |
|----------|-------------|-------------------|
| expire1 | EXPIRE/PERSIST/TTL: set, inspect, and remove expiry | Explain why PERSIST might be needed |
| expire2 | Maxmemory policies: configure and observe eviction | Predict which key gets evicted with LRU |
| expire3 | Keyspace notifications for expiry events | Push it: what's the latency between expiry and notification? |

## Resources

- [Redis key expiration docs](https://redis.io/docs/manual/keyspace-notifications/)
- [EXPIRE command](https://redis.io/commands/expire/)
- [Redis eviction policies](https://redis.io/docs/manual/eviction/)
- [Dragonfly expiration](https://www.dragonflydb.io/docs/managing-dragonfly/eviction-policy)
