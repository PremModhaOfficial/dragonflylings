# Explain It — strings5

## The Challenge

Your API loads a user's profile by making 8 individual GET calls to Redis (name, email, avatar, bio, role, created_at, last_login, preferences). At 200 requests/second, this is 1,600 Redis commands per second.

A code reviewer comments: "You should use MGET here."

Write your explanation of why in `feynman/explanations/strings5.md`.

## Your Explanation Should Cover

1. **What a "round trip" costs** — even on localhost, a round trip takes microseconds. In production (different hosts, network hops), it's 0.5-5ms. Multiply that by 8 calls × 200 req/s.

2. **What MGET changes** — instead of 8 round trips, it's 1. The response is slightly larger (8 values in one message), but the total latency is dominated by the round-trip count, not message size.

3. **When MGET is NOT the right answer** — MGET only works if all keys exist on the same Redis node. In a Redis Cluster, keys on different slots require separate commands (or `{}` key tagging). Dragonfly has its own sharding model.

## Quality Check

Your explanation should include a rough calculation: if MGET saves 7 round trips × 1ms each × 200 req/s = ? ms saved per second. Make the math concrete.
