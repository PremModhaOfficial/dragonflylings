# Module 04: Sets & Sorted Sets — "The Collection"

## Mental Model

**Set**: A bag of unique items (like a jar of unique marbles). Sets answer "is X a member?" in O(1). They also support powerful set operations: intersection, union, difference — perfect for "common friends," "users online," or "items with all these tags."

**Sorted Set (ZSet)**: That same bag of unique marbles, but every marble has a score written on it, and they auto-arrange by score. Sorted Sets answer "what are the top 10?" efficiently. Use them for leaderboards, rate limiters, and priority queues.

```
Set:                     Sorted Set (ZSet):
─────────────────        ──────────────────────────────
{"go", "redis",          member     score    rank(desc)
 "database"}             ──────     ─────    ──────────
                         "diana"    2800     0  ← rank 0 = best
No ordering.             "bob"      2300     1
No scores.               "charlie"  1800     2
Just membership.         "alice"    1500     3
                         "eve"      950      4
```

## Predict Before Starting

Before writing any code, answer in your head:
1. What does SADD return when you add a duplicate?
2. What is SINTER? SUNION? SDIFF?
3. Can a Sorted Set have two members with the same score?
4. What's the difference between ZRANK and ZSCORE?

Write your predictions in `feynman/gap_notebook.md`.

## Key Concepts

### Sets
| Command | Description |
|---------|-------------|
| SADD key member [member ...] | Add members (duplicates ignored) |
| SMEMBERS key | Get all members |
| SISMEMBER key member | Check membership O(1) |
| SINTER key [key ...] | Intersection of sets |
| SUNION key [key ...] | Union of sets |
| SDIFF key [key ...] | Members in first set not in others |

### Sorted Sets
| Command | Description |
|---------|-------------|
| ZADD key score member | Add/update member with score |
| ZRANGE key start stop | Members in ascending score order |
| ZREVRANGE key start stop | Members in descending score order |
| ZSCORE key member | Get a member's score |
| ZRANK / ZREVRANK key member | Get rank (ascending/descending) |
| ZINCRBY key increment member | Atomically add to score |
| ZRANGEBYSCORE / ZREMRANGEBYSCORE | Range/remove by score value |

## Exercises

- **sets1**: SADD/SMEMBERS/SISMEMBER — tag system with deduplication
- **sets2**: SINTER/SUNION/SDIFF — common friends between users
- **zsets1**: ZADD/ZRANGE/ZSCORE — game leaderboard (highest score first)
- **zsets2**: ZRANGEBYSCORE + ZREMRANGEBYSCORE — sliding window rate limiter
- **zsets3**: ZRANK + ZINCRBY — real-time ranking with score updates

## Before You Start

```bash
# Experiment with sets
redis-cli -p 6380 SADD tags:post:1 "go" "redis" "go"  # "go" deduped
redis-cli -p 6380 SMEMBERS tags:post:1                  # {go, redis}

# Experiment with sorted sets
redis-cli -p 6380 ZADD lb 1500 alice 2300 bob 950 eve
redis-cli -p 6380 ZREVRANGE lb 0 -1 WITHSCORES  # bob,2300,alice,1500,eve,950
```
