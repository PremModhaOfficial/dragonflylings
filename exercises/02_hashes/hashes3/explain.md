# Explain It — hashes3

## The Challenge

You're building a dashboard that shows page views per page per day. You could use:
- Option A: One string key per page per day (`pageviews:home:2024-03-25`)
- Option B: One hash per day with page fields (`pageviews:2024-03-25` → `home: 142`)

Write your design decision in `feynman/explanations/hashes3.md`.

## Your Explanation Should Cover

1. **Memory comparison** — a key in Redis has overhead (metadata, pointer, etc.) regardless of its value. Option A creates N keys per day (N = number of pages). Option B creates 1. How does this scale over 365 days × 100 pages?

2. **Query patterns** — "Show me all pages for today." With Option A, you need `SCAN` or `KEYS pageviews:*:2024-03-25`. With Option B, it's `HGETALL pageviews:2024-03-25`. Which is faster? Which is safer in production?

3. **Atomicity of HINCRBY** — two servers increment `home` simultaneously. With HINCRBY, can you lose an increment? Why or why not?

## Quality Check

Your explanation should work as a design document that justifies the hash approach to a skeptical colleague.
