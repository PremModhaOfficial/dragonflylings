# Dragonflylings

Learn Redis and Dragonfly by fixing intentionally broken code. Every exercise has a real bug that teaches a real concept. No hand-holding—you fix, the tests tell you if you're right.

## Prerequisites

- Go 1.21+
- Docker

## Getting Started

```bash
# 1. Start Dragonfly
docker compose up -d

# 2. Build the CLI
go build -o dragonflylings ./cmd/dragonflylings

# 3. See all exercises
./dragonflylings list

# 4. Run the first exercise
./dragonflylings run
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `run` | Run the next incomplete exercise |
| `run <name>` | Run a specific exercise (e.g. `run connect1`) |
| `watch` | Auto-run on file save (great for live feedback) |
| `hint` | Show the next progressive hint for current exercise |
| `hint <name>` | Show hint for a specific exercise |
| `list` | Show all exercises with status |
| `progress` | Show progress summary per module |
| `reset` | Reset all progress |
| `reset <name>` | Reset a specific exercise |
| `explain` | Open the Feynman explain challenge for current exercise |
| `gap` | Open your gap notebook |
| `verify` | Re-verify all completed exercises still pass |

## How It Works

Each exercise has:
- `main.go` — broken code you need to fix
- `main_test.go` — tests that must pass
- `hint.md` — progressive hints (use sparingly)
- `explain.md` — Feynman challenge after you pass
- `gap.md` — questions to add to your gap notebook

Fix the code in `main.go` until the tests pass. The runner tells you what's wrong.

## Feynman Technique Integration

After each exercise, you'll get a Feynman prompt. Write your explanation in `feynman/explanations/<name>.md`. Track questions you can't answer yet in `feynman/gap_notebook.md`.

The goal: by the time you finish all 53 exercises, you can explain every concept to a 12-year-old.

## Modules

| # | Module | Exercises | Core Concept |
|---|--------|-----------|--------------|
| 00 | Connect | 3 | TCP + RESP handshake, connection pooling |
| 01 | Strings | 5 | SET/GET, TTL, atomicity, batching |
| 02 | Hashes | 4 | HSET/HGET, field operations, encoding |
| 03 | Lists | 4 | LPUSH/RPOP, blocking ops, bounded lists |
| 04 | Sets & Sorted Sets | 5 | SADD, ZADD, leaderboards, rate limiting |
| 05 | Expiration | 3 | TTL, eviction policies, keyspace notifications |
| 06 | Pipelining & Tx | 5 | Batching, MULTI/EXEC, WATCH |
| 07 | Pub/Sub | 3 | SUBSCRIBE/PUBLISH, pattern subs, failure modes |
| 08 | Streams | 5 | XADD/XREAD, consumer groups, PEL |
| 09 | Lua Scripting | 4 | EVAL, atomicity, EVALSHA, Dragonfly gotchas |
| 10 | Dragonfly-Specific | 4 | Multi-threading, snapshots, native JSON |
| 11 | Production Patterns | 6 | Cache-aside, locks, rate limiting, circuit breaker |
| 12 | Capstone | 2 | Synthesis — no hints available |

**Total: 53 exercises**

## Dragonfly vs Redis

Dragonfly speaks the Redis protocol—all the same commands work. The differences:
- Multi-threaded (not single-threaded like Redis)
- Runs on port **6380** in this setup (not 6379)
- Snapshot without fork (no memory spike)
- Native JSON support without modules
