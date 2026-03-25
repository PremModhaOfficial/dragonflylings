# Explain It — prod4

Your company's session storage is Redis-backed. The Redis instance goes down for 2 minutes. Three engineers propose different responses:

- Engineer A: *"Return 401 Unauthorized for all requests during the outage."*
- Engineer B: *"Return empty session data — the app degrades gracefully."*
- Engineer C: *"Cache sessions in-process so we can serve during the outage."*

Write your answer in `feynman/explanations/prod4.md`.

Evaluate each approach:
1. What does "graceful degradation" mean for session storage? (What features break? What stays working?)
2. When is returning an empty session (B) acceptable vs. dangerous?
3. What are the risks of in-process session caching (C) when Redis comes back up?
4. What would you actually do? (There's no single right answer — explain your trade-offs.)

**The sliding TTL question:** Most sessions expire after 30 minutes of INACTIVITY, not 30 minutes from creation. How does `GetSession` refreshing the TTL implement this? What's the edge case if a user makes a request every 29 minutes forever?

QUALITY CHECK: Could you explain "graceful degradation" to a non-technical product manager in one paragraph?
