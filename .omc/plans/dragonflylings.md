# Dragonflylings: Complete Build Plan

**Date:** 2025-03-25
**Status:** AWAITING CONFIRMATION
**Estimated learner completion:** 20-40 hours
**Estimated build time:** ~1 hour with agent team

---

## 1. Project Directory Structure

```
dragonflylings/
├── cmd/
│   └── dragonflylings/
│       └── main.go                 # CLI runner entry point
├── internal/
│   └── runner/
│       ├── runner.go               # Exercise runner (compile, test, verify)
│       ├── watcher.go              # Filesystem watcher for watch mode
│       ├── progress.go             # State tracking (JSON file)
│       ├── hint.go                 # Hint display logic
│       └── printer.go             # Terminal output formatting (colors, boxes)
├── exercises/
│   ├── exercises.toml              # Manifest: single source of truth
│   ├── 00_connect/
│   │   ├── connect1/
│   │   │   ├── main.go            # Broken code (learner edits this)
│   │   │   ├── main_test.go       # Tests that must pass
│   │   │   ├── hint.md            # Progressive hints (reveal one at a time)
│   │   │   ├── explain.md         # "Explain it" challenge prompt
│   │   │   └── gap.md             # Gap notebook prompt ("what don't you know yet?")
│   │   ├── connect2/
│   │   │   └── ...
│   │   └── README.md              # Module intro: mental model, analogy, predict prompt
│   ├── 01_strings/
│   │   └── ...
│   ├── ... (all modules)
│   └── 12_capstone/
│       └── ...
├── lib/
│   └── testutil/
│       └── testutil.go            # Shared test helpers (connect, cleanup, assert)
├── feynman/
│   ├── gap_notebook.md            # Learner's running gap notebook
│   └── explanations/              # Learner writes their "explain it" answers here
│       └── .gitkeep
├── docker-compose.yml             # Dragonfly on port 6380
├── exercises.toml                 # Symlink or copy to exercises/exercises.toml
├── go.mod
├── go.sum
├── Makefile                       # Convenience targets
└── README.md                      # Getting started guide
```

---

## 2. Exercise Manifest Format (exercises.toml)

```toml
# Each exercise is the single source of truth for ordering and metadata

[[exercises]]
name = "connect1"
dir = "00_connect/connect1"
mode = "test"                    # "compile" or "test"
module = "00_connect"
hint_count = 3                   # Number of progressive hints available
feynman = "explain"              # "explain" | "predict" | "limit-test" | "none"
description = "Establish a connection to Dragonfly and verify with PING"

[[exercises]]
name = "connect2"
dir = "00_connect/connect2"
mode = "test"
module = "00_connect"
hint_count = 2
feynman = "limit-test"
description = "Handle connection failures gracefully with timeouts"
```

---

## 3. Exercise File Format Specification

### main.go (Broken Code)
```go
// exercises/00_connect/connect1/main.go
package main

// EXERCISE: connect1 - Your First PING
//
// PREDICT: Before writing any code, answer in your head:
//   What network protocol does Redis/Dragonfly use?
//   What do you think PING returns? Why would that command exist?
//
// TODO: Fix the code below so the test passes.
// The connection to Dragonfly should use localhost:6380.

import (
    "context"
    "github.com/redis/go-redis/v9"
)

func Connect() *redis.Client {
    // FIX ME: The options are wrong
    return redis.NewClient(&redis.Options{
        Addr: "localhost:6379", // <-- wrong port
    })
}

func Ping(client *redis.Client) (string, error) {
    ctx := context.Background()
    // FIX ME: call the right method
    return client.Echo(ctx, "hello").Result() // <-- should be Ping
}
```

### main_test.go (Tests)
```go
package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
    client := Connect()
    require.NotNil(t, client)
    defer client.Close()

    result, err := Ping(client)
    require.NoError(t, err)
    assert.Equal(t, "PONG", result)
}
```

### hint.md (Progressive Hints)
```markdown
## Hint 1
Dragonfly runs on port 6380 in our docker-compose setup, not the Redis default.

## Hint 2
The `Ping` method on `*redis.Client` sends a PING command and returns "PONG".
`Echo` sends ECHO, which is a different command entirely.

## Hint 3
```go
func Ping(client *redis.Client) (string, error) {
    ctx := context.Background()
    return client.Ping(ctx).Result()
}
```
```

### explain.md (Feynman Challenge)
```markdown
## Explain It Challenge

Imagine your friend just asked: "Why does Redis have a PING command?
Isn't that just wasting a network round-trip?"

Write your answer in `feynman/explanations/connect1.md`.

Your explanation should cover:
- What PING actually verifies (not just "is the server up")
- Why a health check matters in distributed systems
- What could go wrong if you skip connection verification

QUALITY CHECK: Could a 12-year-old follow your explanation?
If you used jargon, replace it with an analogy.
```

### gap.md (Gap Notebook Prompt)
```markdown
## Gap Notebook

After completing this exercise, add to `feynman/gap_notebook.md`:

1. What's one thing about Redis connections you're still unsure about?
2. What would happen if Dragonfly was under heavy load during your PING?
3. Do you know the difference between TCP connection and Redis protocol handshake?

Format: `- [ ] [your question] -- [why it matters]`
```

---

## 4. Complete Exercise Curriculum

### Module 00: Connect (3 exercises) -- "The Handshake"
**Mental Model:** A Redis connection is like a phone call. You dial (TCP), someone answers (RESP handshake), you say "hello?" (PING), they say "I'm here" (PONG). Only then do you talk.

| # | Exercise | Mode | Feynman | Description |
|---|----------|------|---------|-------------|
| 1 | connect1 | test | explain | Fix connection options and PING to verify connectivity |
| 2 | connect2 | test | predict | Handle connection timeouts -- what happens when Dragonfly is unreachable? |
| 3 | connect3 | test | limit-test | Connection pooling: configure PoolSize, MinIdleConns, observe behavior |

---

### Module 01: Strings (5 exercises) -- "The Atoms"
**Mental Model:** Strings are Redis's atoms -- the indivisible unit everything else builds on. But "string" is misleading: they're really "binary-safe byte sequences up to 512MB." A string can hold a serialized protobuf as easily as the word "hello."

| # | Exercise | Mode | Feynman | Description |
|---|----------|------|---------|-------------|
| 1 | strings1 | test | predict | SET and GET: fix the key-value basics |
| 2 | strings2 | test | explain | Expiration with TTL: SET with EX, TTL command, what happens after expiry |
| 3 | strings3 | test | limit-test | SETNX (set-if-not-exists): the foundation of distributed locking |
| 4 | strings4 | test | explain | INCR/DECR atomicity: why this isn't just GET + math + SET |
| 5 | strings5 | test | predict | MGET/MSET: batch operations, why fewer round trips matter |

---

### Module 02: Hashes (4 exercises) -- "The Row"
**Mental Model:** If a String is a single cell in a spreadsheet, a Hash is an entire row. One key, many fields. You wouldn't create 10 separate keys for a user's name, email, age -- you'd use one hash. It's the difference between 10 phone calls and one phone call with 10 questions.

| # | Exercise | Mode | Feynman | Description |
|---|----------|------|---------|-------------|
| 1 | hashes1 | test | predict | HSET/HGET: model a user profile as a hash |
| 2 | hashes2 | test | explain | HMSET/HMGET: batch field operations, partial reads |
| 3 | hashes3 | test | limit-test | HINCRBY for counters inside hashes (page view tracker) |
| 4 | hashes4 | test | explain | Hash vs. multiple Strings: when to use which (memory analysis) |

---

### Module 03: Lists (4 exercises) -- "The Queue"
**Mental Model:** A Redis List is a double-ended queue (deque). Think of a conveyor belt where you can place items on either end and take from either end. This makes it perfect for job queues, activity feeds, and message buffers.

| # | Exercise | Mode | Feynman | Description |
|---|----------|------|---------|-------------|
| 1 | lists1 | test | predict | LPUSH/RPUSH/LPOP/RPOP: build a simple task queue |
| 2 | lists2 | test | explain | LRANGE: paginate a list, understand 0-based vs -1 indexing |
| 3 | lists3 | test | limit-test | BLPOP: blocking pop as a worker pattern (with timeout) |
| 4 | lists4 | test | explain | LLEN + LTRIM: bounded lists for "last N items" pattern |

---

### Module 04: Sets and Sorted Sets (5 exercises) -- "The Collection"
**Mental Model:** A Set is a bag of unique items (like a jar of unique marbles). A Sorted Set is that same jar, but every marble has a score written on it, and they auto-arrange by score. Sets answer "is X a member?" in O(1). Sorted Sets answer "what are the top 10?" efficiently.

| # | Exercise | Mode | Feynman | Description |
|---|----------|------|---------|-------------|
| 1 | sets1 | test | predict | SADD/SMEMBERS/SISMEMBER: tag system with unique tags |
| 2 | sets2 | test | explain | SINTER/SUNION/SDIFF: find common friends between users |
| 3 | zsets1 | test | predict | ZADD/ZRANGE/ZSCORE: build a leaderboard |
| 4 | zsets2 | test | limit-test | ZRANGEBYSCORE + ZREMRANGEBYSCORE: sliding window rate limiter |
| 5 | zsets3 | test | explain | ZRANK for real-time ranking, ZINCRBY for score updates |

---

### Module 05: Expiration and Memory (3 exercises) -- "The Janitor"
**Mental Model:** Redis is an in-memory store. Without expiration, it's a bathtub with the drain plugged -- eventually it overflows. TTL is the drain. But Redis doesn't check every key every second (too expensive). It uses a mix of lazy deletion (check on access) and periodic sampling (random spot checks). Understanding this matters for production.

| # | Exercise | Mode | Feynman | Description |
|---|----------|------|---------|-------------|
| 1 | expire1 | test | explain | EXPIRE/PERSIST/TTL: set, check, and remove expiration |
| 2 | expire2 | test | predict | Maxmemory policies: what happens when memory is full? (allkeys-lru, volatile-lru, noeviction) |
| 3 | expire3 | test | limit-test | Key space notifications: subscribe to expiration events |

---

### Module 06: Pipelining and Transactions (5 exercises) -- "The Assembly Line"
**Mental Model:** Without pipelining, each command is a separate letter mailed to the server -- you wait for a reply before sending the next. Pipelining stuffs many letters in one envelope. Transactions (MULTI/EXEC) are a signed contract: all letters execute together, or none do. But WATCH adds a condition: "only execute if nobody changed X while I was preparing."

| # | Exercise | Mode | Feynman | Description |
|---|----------|------|---------|-------------|
| 1 | pipe1 | test | predict | Pipeline: batch 100 SETs, compare latency vs. individual commands |
| 2 | pipe2 | test | explain | Pipeline error handling: partial failures (some cmds fail, others succeed) |
| 3 | tx1 | test | explain | MULTI/EXEC: transfer balance between two keys atomically |
| 4 | tx2 | test | limit-test | WATCH/MULTI: optimistic locking, detect and retry on conflict |
| 5 | tx3 | test | predict | Pipeline vs. Transaction vs. Lua script: when to use each (comparison exercise) |

---

### Module 07: Pub/Sub (3 exercises) -- "The Megaphone"
**Mental Model:** Pub/Sub is a megaphone in a stadium. The publisher shouts, and everyone currently listening hears it. If you arrive late, you missed it -- there's no replay. This makes Pub/Sub great for real-time broadcasts but terrible for reliable messaging (that's what Streams are for).

| # | Exercise | Mode | Feynman | Description |
|---|----------|------|---------|-------------|
| 1 | pubsub1 | test | predict | SUBSCRIBE/PUBLISH: basic chat room with two goroutines |
| 2 | pubsub2 | test | explain | Pattern subscriptions (PSUBSCRIBE): wildcard channels |
| 3 | pubsub3 | test | limit-test | Pub/Sub failure modes: what happens when subscriber disconnects and reconnects? (fire-and-forget proof) |

---

### Module 08: Streams (5 exercises) -- "The Log"
**Mental Model:** If Pub/Sub is a megaphone, a Stream is a recorded lecture. Every message is stored with an ID (timestamp-sequence). Consumer Groups are study groups: each group processes every message, but within a group, each student handles different messages. If a student drops out mid-message (crashes), the Pending Entries List (PEL) tracks it for redelivery.

| # | Exercise | Mode | Feynman | Description |
|---|----------|------|---------|-------------|
| 1 | streams1 | test | predict | XADD/XREAD: produce and consume events from a stream |
| 2 | streams2 | test | explain | XRANGE/XLEN: query stream history by time range |
| 3 | streams3 | test | explain | Consumer groups: XGROUP CREATE, XREADGROUP, XACK |
| 4 | streams4 | test | limit-test | PEL and redelivery: XPENDING + XCLAIM for crash recovery |
| 5 | streams5 | test | predict | XTRIM + maxlen: bounded streams for production use |

---

### Module 09: Lua Scripting (4 exercises) -- "The Stored Procedure"
**Mental Model:** Lua scripts in Redis are like stored procedures in SQL: they run atomically on the server side. Instead of 5 round trips (read, compute, write, read, write), you send one script that does all 5 steps without anyone else interfering. The tradeoff: you block the server while your script runs.

| # | Exercise | Mode | Feynman | Description |
|---|----------|------|---------|-------------|
| 1 | lua1 | test | predict | EVAL basics: a Lua script that does conditional SET (compare-and-swap) |
| 2 | lua2 | test | explain | KEYS vs ARGV: why Redis requires you to declare keys upfront |
| 3 | lua3 | test | limit-test | EVALSHA + script caching: avoid re-sending the script every time |
| 4 | lua4 | test | explain | Dragonfly gotcha: multi-shard Lua requires {hashtag} key prefixes |

---

### Module 10: Dragonfly-Specific Features (4 exercises) -- "The Dragon"
**Mental Model:** Dragonfly looks like Redis from the outside (same protocol, same commands) but inside it's a different beast. Redis is single-threaded (one cook in the kitchen). Dragonfly is multi-threaded with shared-nothing shards (many cooks, each owning their own counter). This means some Redis assumptions break.

| # | Exercise | Mode | Feynman | Description |
|---|----------|------|---------|-------------|
| 1 | dragon1 | test | explain | Multi-threaded architecture: observe parallel throughput with concurrent goroutines |
| 2 | dragon2 | test | predict | Snapshot without fork: trigger BGSAVE, observe zero memory spike (vs Redis 2x memory) |
| 3 | dragon3 | test | limit-test | Native JSON: JSON.SET / JSON.GET without modules (compare to RedisJSON) |
| 4 | dragon4 | test | explain | Dragonfly gotchas: WAIT behavior, OBJECT ENCODING differences, no Sentinel |

---

### Module 11: Production Patterns (6 exercises) -- "The Battlefield"
**Mental Model:** Knowing Redis commands is like knowing chess piece movements. Production is playing actual games. This module teaches the patterns that separate "I've used Redis" from "I've run Redis in production."

| # | Exercise | Mode | Feynman | Description |
|---|----------|------|---------|-------------|
| 1 | prod1 | test | explain | Cache-aside pattern: read-through with thundering herd protection (singleflight) |
| 2 | prod2 | test | limit-test | Distributed lock: SETNX with expiry, handle lock extension, compare to Redlock |
| 3 | prod3 | test | explain | Rate limiter: sliding window using sorted sets (production-grade) |
| 4 | prod4 | test | predict | Session storage: hash-based sessions with expiry, graceful degradation |
| 5 | prod5 | test | limit-test | Hot key mitigation: detect and handle with local caching layer |
| 6 | prod6 | test | explain | Circuit breaker: wrap Redis client with fallback when Dragonfly is degraded |

---

### Module 12: Capstone (2 exercises) -- "The Summit"
**Mental Model:** No new concepts. You combine everything. These exercises have no hints. You're the expert now.

| # | Exercise | Mode | Feynman | Description |
|---|----------|------|---------|-------------|
| 1 | capstone1 | test | explain | Build a real-time event processing pipeline: Streams + Consumer Groups + Lua + Sorted Set leaderboard |
| 2 | capstone2 | test | explain | Build a rate-limited API cache: cache-aside + sliding window rate limiter + circuit breaker + pub/sub invalidation |

---

**Total: 13 modules, 53 exercises**

---

## 5. Runner CLI Design

### Commands

```
dragonflylings run              # Run next incomplete exercise
dragonflylings run <name>       # Run specific exercise
dragonflylings watch            # Watch mode: auto-run on file save
dragonflylings hint             # Show next hint for current exercise
dragonflylings hint <name>      # Show next hint for specific exercise
dragonflylings list             # Show all exercises with status (done/pending)
dragonflylings progress         # Show progress summary per module
dragonflylings reset            # Reset all progress
dragonflylings reset <name>     # Reset specific exercise
dragonflylings explain          # Open the explain challenge for current exercise
dragonflylings gap              # Open gap notebook
dragonflylings verify           # Re-verify all completed exercises still pass
```

### Runner Logic

1. Parse `exercises.toml` to get ordered exercise list
2. Load progress state from `.dragonflylings-state.json`
3. For each exercise attempt:
   a. **Compile check:** `go build ./exercises/{dir}/...`
   b. **Test check (if mode=test):** `go test ./exercises/{dir}/... -v -count=1 -timeout 30s`
   c. On success: mark done in state, print congratulations + Feynman prompt
   d. On failure: print error output + "use `dragonflylings hint` for help"
4. Watch mode: use `fsnotify` to watch `exercises/` directory, re-run on `.go` file changes

### State File Format (.dragonflylings-state.json)

```json
{
  "version": 1,
  "exercises": {
    "connect1": { "status": "done", "completed_at": "2025-03-25T10:00:00Z", "hints_used": 1 },
    "connect2": { "status": "pending", "hints_used": 0 },
    "strings1": { "status": "pending", "hints_used": 0 }
  },
  "current": "connect2",
  "stats": {
    "total": 53,
    "completed": 1,
    "hints_used": 1,
    "started_at": "2025-03-25T09:00:00Z"
  }
}
```

### Progress Display

```
DRAGONFLYLINGS PROGRESS
=======================

  00 Connect        [###.] 1/3
  01 Strings        [.....] 0/5
  02 Hashes         [.....] 0/4
  03 Lists          [.....] 0/4
  04 Sets & ZSets   [.....] 0/5
  05 Expiration     [.....] 0/3
  06 Pipelining/Tx  [.....] 0/5
  07 Pub/Sub        [.....] 0/3
  08 Streams        [.....] 0/5
  09 Lua Scripting  [.....] 0/4
  10 Dragonfly      [.....] 0/4
  11 Production     [.....] 0/6
  12 Capstone       [.....] 0/2

  Total: 1/53 (1.9%) | Hints used: 1
  Current: connect2
```

---

## 6. Feynman Integration Points Per Module

Each module follows the PREDICT > BREAK > REBUILD > EXPLAIN > LIMIT-TEST cycle:

| Module | Predict Moment | Break Moment | Rebuild Insight | Explain Challenge | Limit-Test |
|--------|---------------|-------------|----------------|-------------------|------------|
| 00 Connect | "What does PING return?" | Wrong port, wrong method | Connection = TCP + RESP + auth | "Why do health checks exist?" | Pool exhaustion under load |
| 01 Strings | "Is INCR just GET+1+SET?" | Race condition proof | Atomic operations at protocol level | "Explain TTL to a 12yo" | What happens at 512MB? |
| 02 Hashes | "When is a hash better than separate keys?" | Memory comparison | Hash ziplist encoding optimization | "Hash vs JSON blob tradeoffs" | Field count threshold for performance |
| 03 Lists | "What happens if BLPOP times out?" | Blocking without timeout = hang | Blocking ops need context deadlines | "Queue vs List vs Stream" | What if producer is faster than consumer? |
| 04 Sets/ZSets | "Can a sorted set have duplicate scores?" | Yes, and ordering is lexicographic | Score = priority, member = identity | "Build a leaderboard from scratch" | ZRANGEBYSCORE at scale |
| 05 Expiration | "Does Redis check every key's TTL constantly?" | Lazy + sampling proof | Probabilistic expiration model | "Why lazy deletion?" | Expired keys still using memory |
| 06 Pipeline/Tx | "Do transactions roll back on error?" | No rollback -- partial execution | WATCH = optimistic locking | "Pipeline vs Tx vs Lua decision tree" | WATCH race window |
| 07 Pub/Sub | "What happens to messages when no one listens?" | Lost forever -- prove it | Fire-and-forget by design | "When NOT to use Pub/Sub" | Slow subscriber backpressure |
| 08 Streams | "How are Streams different from Pub/Sub?" | Replay proof | Append-only log with consumer groups | "Explain PEL to a junior dev" | Consumer group rebalancing |
| 09 Lua | "Why not just pipeline instead of scripting?" | Atomicity proof (interleaving) | Server-side execution = atomic block | "When does Lua become dangerous?" | Long-running script blocking |
| 10 Dragonfly | "Is Dragonfly just faster Redis?" | Behavioral differences proof | Shared-nothing vs single-thread | "Explain fiber-per-shard to a 12yo" | Multi-shard Lua gotchas |
| 11 Production | "Is a distributed lock just SETNX?" | Lock expiry during long operations | Lock extension + fencing tokens | "Design a rate limiter on a whiteboard" | Thundering herd at scale |
| 12 Capstone | N/A (synthesis) | Self-directed | Integration | "Explain your entire system" | Self-directed stress testing |

---

## 7. Team Composition and Task Assignment

### Agent Team (7 agents, ~1 hour)

| Agent | Role | Model | Tasks | Time |
|-------|------|-------|-------|------|
| **Executor 1: Scaffolder** | Project structure + CLI | opus | Create directory structure, go.mod, docker-compose.yml, Makefile, README. Build complete runner CLI (main.go, runner.go, watcher.go, progress.go, hint.go, printer.go). Write testutil helpers. | 20 min |
| **Executor 2: Exercises A** | Modules 00-04 (21 exercises) | sonnet | Write all main.go, main_test.go, hint.md, explain.md, gap.md for Connect, Strings, Hashes, Lists, Sets/ZSets | 25 min |
| **Executor 3: Exercises B** | Modules 05-08 (16 exercises) | sonnet | Write all exercise files for Expiration, Pipelining/Tx, Pub/Sub, Streams | 25 min |
| **Executor 4: Exercises C** | Modules 09-12 (16 exercises) | sonnet | Write all exercise files for Lua, Dragonfly-specific, Production Patterns, Capstone | 25 min |
| **Executor 5: Manifest** | exercises.toml + module READMEs | sonnet | Write complete exercises.toml manifest, all 13 module README.md files with mental models and analogies | 10 min |
| **Verifier** | Test all exercises | sonnet | Start Dragonfly container, run each exercise's solution against tests, verify all pass. Fix any broken tests. | 15 min |
| **Code Reviewer** | Quality pass | opus | Review exercise pedagogy, code quality, hint quality, Feynman integration. Flag issues for fix. | 10 min |

### Execution Timeline

```
T+00:00  Scaffolder starts: project structure + CLI runner
T+00:00  Exercises A/B/C start in parallel (wait for scaffolder's testutil at ~T+05:00)
T+00:05  Manifest agent starts (needs to know final exercise list)
T+00:20  Scaffolder done. CLI runner complete.
T+00:25  Exercise agents finishing up
T+00:25  Verifier starts: spins up Dragonfly, runs all tests
T+00:35  Code reviewer starts: pedagogy and code quality review
T+00:40  Verifier reports issues, exercise agents fix
T+00:50  Final verification pass
T+01:00  Done. All 53 exercises passing, CLI working, Feynman content reviewed.
```

### Dependency Graph

```
Scaffolder ──────────────────> Verifier ──> Code Reviewer
    │                              ^
    ├── Exercises A (00-04) ───────┤
    ├── Exercises B (05-08) ───────┤
    ├── Exercises C (09-12) ───────┤
    └── Manifest + READMEs ────────┘
```

### Critical Path: Scaffolder (20 min) + Exercises (25 min parallel) + Verification (15 min) = 60 min

---

## 8. Guardrails

### Must Have
- Every exercise must compile and have passing tests when the correct solution is applied
- Every exercise must FAIL with the initial broken code (the break must be real)
- Real Dragonfly container required (docker-compose.yml included)
- Progressive hints: never give away the answer in hint 1 or 2
- Each module README must establish the mental model BEFORE exercises begin
- Gap notebook prompts must ask questions the exercises DON'T answer (push beyond)
- Runner CLI must handle Dragonfly being down gracefully (clear error message)

### Must NOT Have
- No mock/fake Redis -- all exercises hit real Dragonfly
- No exercises that test Go knowledge rather than Redis/Dragonfly knowledge
- No hints that just say "read the docs" -- hints must teach
- No exercises longer than ~30 minutes each
- No external dependencies beyond go-redis/v9, testify, and fsnotify
- No exercises that require multiple Dragonfly instances (keep single-node)

---

## 9. Success Criteria

1. `docker compose up -d` starts Dragonfly
2. `go run ./cmd/dragonflylings list` shows all 53 exercises
3. `go run ./cmd/dragonflylings run` starts at connect1
4. Each exercise has a clear, intentional bug that maps to a concept
5. `go run ./cmd/dragonflylings watch` re-runs on file save
6. `go run ./cmd/dragonflylings hint` shows progressive hints
7. A learner completing all 53 exercises can:
   - Explain every Redis data structure and its time complexity
   - Design cache-aside, rate limiting, and distributed locking from scratch
   - Debug production Redis/Dragonfly issues
   - Explain Dragonfly's architecture differences from Redis
   - Write Lua scripts for atomic operations
   - Handle Streams with consumer groups and failure recovery

---

## 10. Open Decisions for Execution

These are resolved in the plan but noted for executor awareness:

- **Test isolation:** Each test should FLUSHDB in setup or use unique key prefixes to avoid cross-test pollution. Recommend unique key prefixes (safer than FLUSHDB in case of parallel runs).
- **Go module name:** Keep existing `dragonflyLearnings` or rename to `dragonflylings`? Recommend rename to `dragonflylings` for consistency.
- **Existing cmd/main.go:** Archive to `cmd/sandbox/main.go` as a free-form playground, distinct from the exercise runner.
