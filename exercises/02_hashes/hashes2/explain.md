# Explain It — hashes2

## The Challenge

Your API loads a user's profile page by reading 6 hash fields: name, avatar, bio, role, created_at, last_seen. A new teammate asks: "Why use HMGET? Isn't looping with HGET the same thing?"

Write your explanation in `feynman/explanations/hashes2.md`.

## Your Explanation Should Cover

1. **Round trips vs bandwidth** — HMGET sends one request and waits for one response. HGet loop sends 6 requests and waits 6 times. The cost isn't bandwidth (the data is the same size either way) — it's latency multiplied by round trips.

2. **The partial read consistency issue** — with 6 HGETs, another process can update the user between reads. With HMGET, all 6 fields are returned together. Are they guaranteed to be consistent? (Yes — HMGET is a single command.)

3. **Error handling difference** — HMGET returns nil for missing fields (not an error). HGET returns `redis.Nil` error. Why does this difference matter for robustness?

## Quality Check

Include a concrete latency calculation: if each round trip is 1ms, how much time does your API spend waiting for Redis for a 6-field profile load? With HMGET vs HGET loop?
