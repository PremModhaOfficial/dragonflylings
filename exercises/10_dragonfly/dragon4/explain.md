# Explain It — dragon4

Your team is migrating from Redis to Dragonfly. Three engineers report different bugs:

1. *"Our health check fails — it calls WAIT and gets 0 back, which we treat as an error."*
2. *"Our memory optimizer crashes — it calls OBJECT ENCODING and panics on Dragonfly."*
3. *"Our HA setup is broken — we use Redis Sentinel and Dragonfly doesn't support it."*

Write your answer in `feynman/explanations/dragon4.md`.

For each bug, explain:
- Why this works on Redis but breaks on Dragonfly
- What the correct behavior is (what Dragonfly actually does)
- How to fix the code to work with both Redis and Dragonfly

**The meta-question:** All three bugs share a pattern. What is it? Write a general principle for writing Redis-compatible code that also runs correctly on Dragonfly.

**Sentinel alternative:** If Dragonfly doesn't support Sentinel, how do you achieve high availability with Dragonfly? What does the Dragonfly team recommend?

QUALITY CHECK: Can you explain these three bugs in a 5-minute team standup without slides or demos?
